package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/juevigrace/diva-server/concurrency"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/router"
	"github.com/juevigrace/diva-server/storage"
)

type Server struct {
	srv    *http.Server
	config *ServerConfig

	database storage.Storage
	router   *router.ServerRouter
}

func NewServer(config models.Config) (*Server, error) {
	var server *Server = new(Server)
	server.config = NewServerConfig().(*ServerConfig)
	server.config.Configure(config)

	if err := server.setup(); err != nil {
		return nil, err
	}
	return server, nil
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	notifyCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	done := make(chan error, 1)
	go func() {
		s.printBanner()
		s.router.Routes()
		if err := s.srv.ListenAndServe(); err != http.ErrServerClosed {
			done <- err
		}
		close(done)
	}()

	select {
	case <-notifyCtx.Done():
		log.Println("shutting down gracefully, press Ctrl+C again to force")
		return s.shutdown(ctx)
	case err := <-done:
		return err
	}
}

func (s *Server) setup() error {
	if err := s.createStorage(); err != nil {
		return err
	}

	s.router = router.NewServerRouter(s.database)

	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: s.router.Chi,
	}

	return nil
}

func (s *Server) createStorage() error {
	databaseConf := storage.NewDatabaseConf()
	database, err := storage.New(databaseConf)
	if err != nil {
		return fmt.Errorf("failed to create storage: %w", err)
	}

	s.database = database

	return nil
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
