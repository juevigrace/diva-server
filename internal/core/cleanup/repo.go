package cleanup

import (
	"context"
	"log"
	"time"

	"github.com/juevigrace/diva-server/storage"
)

type CleanupService struct {
	sessionStore    storage.SessionStore
	permissionStore storage.UserPermissionStore
	actionStore     storage.UserActionStore
	interval        time.Duration
	stopCh          chan struct{}
}

func NewCleanupService(
	sessionStore storage.SessionStore,
	permissionStore storage.UserPermissionStore,
	actionStore storage.UserActionStore,
) *CleanupService {
	return &CleanupService{
		sessionStore:    sessionStore,
		permissionStore: permissionStore,
		actionStore:     actionStore,
		interval:        15 * time.Minute,
		stopCh:          make(chan struct{}),
	}
}

func (s *CleanupService) Start() {
	go func() {
		s.run()

		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.run()
			case <-s.stopCh:
				log.Println("cleanup service stopped")
				return
			}
		}
	}()
	log.Println("cleanup service started")
}

func (s *CleanupService) Stop() {
	close(s.stopCh)
}

func (s *CleanupService) run() {
	ctx := context.Background()

	if err := s.sessionStore.CloseExpiredSessions(ctx); err != nil {
		log.Printf("cleanup: close expired sessions: %v", err)
	}

	if err := s.permissionStore.DeleteExpiredUserPermissions(ctx); err != nil {
		log.Printf("cleanup: delete expired user permissions: %v", err)
	}

	if err := s.actionStore.DeleteExpiredActions(ctx); err != nil {
		log.Printf("cleanup: delete expired actions: %v", err)
	}
}
