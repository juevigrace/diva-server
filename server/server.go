package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/api/core/auth"
	"github.com/juevigrace/diva-server/internal/api/core/permission"
	"github.com/juevigrace/diva-server/internal/api/core/session"
	"github.com/juevigrace/diva-server/internal/api/core/user"
	"github.com/juevigrace/diva-server/internal/api/core/verification"
	"github.com/juevigrace/diva-server/internal/api/middlewares"
	"github.com/juevigrace/diva-server/internal/config"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/pkg/concurrency"
	"github.com/juevigrace/diva-server/pkg/filehelper"
	"github.com/juevigrace/diva-server/pkg/mail"
	"github.com/juevigrace/diva-server/storage"
)

type Server struct {
	srv      *http.Server
	serverCh chan error

	config *ServerConfig

	database storage.Storage
	router   *chi.Mux
	mail     *mail.Client
	files    *filehelper.FileHelper
}

func NewServer(cfg config.Config, database storage.Storage) (*Server, error) {
	server := new(Server)
	server.config = cfg.(*ServerConfig)
	server.config.LoadFromEnv()

	server.database = database

	return server, nil
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	notifyCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s.serverCh = make(chan error, 1)
	defer close(s.serverCh)

	s.start()

	select {
	case <-notifyCtx.Done():
		log.Println("shutting down gracefully, press Ctrl+C again to force")
		return s.shutdown(ctx)
	case err := <-s.serverCh:
		return err
	}
}

func (s *Server) start() {
	s.mail = mail.NewClient(s.config.ResendAPIKey, s.config.ResendFromEmail)
	s.files = filehelper.NewFileHelper(s.config.UploadsDir, make(map[string]string), 0)
	s.setupRouter()
	s.setupApi()
	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: s.router,
	}

	s.printBanner()
	go func() {
		if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
			s.serverCh <- err
		}
	}()
}

func (s *Server) setupApi() {
	apiLimiter := middlewares.NewRateLimiter(60, 1*time.Minute)

	pModule := permission.NewPermissionModule(s.database.PermissionStore())
	sModule := session.NewSessionModule(s.database.SessionStore())
	uModule := user.NewUserModule(
		s.database.UserStore(),
		s.database.UserActionStore(),
		s.database.UserPermissionStore(),
		s.database.UserPreferenceStore(),
		s.database.UserProfileStore(),
		s.database.UserStateStore(),
		pModule.Repo,
		sModule.Repo,
		sModule.Handler,
		s.files,
	)
	vModule := verification.NewVerificationModule(s.mail, s.database.UserVerificationStore(), uModule)
	aModule := auth.NewAuthModule(pModule.Repo, sModule.Repo, uModule.URepo, vModule.Repo)

	root := s.router.Route("/", func(root chi.Router) {
		root.Use(apiLimiter.Middleware)
		root.Get("/", func(w http.ResponseWriter, r *http.Request) {
			// TODO: main page
		})

		root.Route("/api", func(api chi.Router) {
			uModule.Routes(api)
			sModule.Routes(api, uModule.URepo.GetByID)
			aModule.Routes(api)
			pModule.Routes(api, sModule.Repo.GetByID, uModule.URepo.GetByID)
			vModule.Routes(api)
		})
	})

	root.With(
		middlewares.RequiresSession(sModule.Repo.GetByID, uModule.URepo.GetByID),
		middlewares.RequireRole(models.ROLE_MODERATOR, models.ROLE_ADMIN),
	).Get("/health", func(w http.ResponseWriter, r *http.Request) {
		res := responses.RespondOk(s.database.Health(context.Background()), "Success")
		responses.WriteJSON(w, res)
	})

	fileServer := http.FileServer(http.Dir(s.config.UploadsDir))
	s.router.Handle("/uploads", fileServer)

	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		res := responses.RespondNotFound(nil, "Route not found")
		responses.WriteJSON(w, res)
	})
}

func (s *Server) shutdown(ctx context.Context) error {
	return concurrency.WithTimeout(ctx, 1*time.Minute, func(ctx context.Context) error {
		if err := s.database.Close(); err != nil {
			return err
		}

		if err := s.srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("server shutdown: %v", err)
		}
		return nil
	})
}

func (s *Server) printBanner() {
	const cyan = "\033[38;5;51m"
	const reset = "\033[0m"

	banner := fmt.Sprintf(`
%s     ____  _____    _____       ____ ___
%s    / __ \/  _/ |  / /   |     / __ <  /
%s   / / / // / | | / / /| |    / / / / /
%s  / /_/ // /  | |/ / ___ |   / /_/ / /
%s /_____/___/  |___/_/  |_|   \____/_/

%s%s`,
		cyan, cyan, cyan, cyan, cyan, reset,
		fmt.Sprintf("==> http server started on port: %d", s.config.Port),
	)
	fmt.Println(banner)
}
