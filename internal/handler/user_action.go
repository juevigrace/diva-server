package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type UserActionsHandler struct {
	service *service.UserActionsService
}

func NewUserActionsHandler(svc *service.UserActionsService) *UserActionsHandler {
	return &UserActionsHandler{service: svc}
}

func (h *UserActionsHandler) Routes(r chi.Router) {
	r.Route("/actions", func(auth chi.Router) {
		auth.Get("/", h.getActions)
		auth.Get("/stream", h.streamActions)
	})
}

func (h *UserActionsHandler) getActions(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	actions, err := h.service.GetAll(r.Context(), &session.User.ID)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	result := make([]responses.ActionResponse, len(actions))
	for i, a := range actions {
		result[i] = responses.ActionResponse{
			ActionName: a.Action.String(),
			ID:         a.ID.String(),
		}
	}

	responses.WriteJSON(w, responses.RespondOk(result, "Success"))
}

func (h *UserActionsHandler) streamActions(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	flusher, ok := responses.SSEHeaders(w)
	if !ok {
		responses.WriteJSON(w, responses.RespondInternalServerError(nil, "SSE not supported"))
		return
	}

	actionsCh := make(chan []*responses.ActionResponse)
	errCh := make(chan error)
	go h.streamActionsWorker(r.Context(), &session.User.ID, actionsCh, errCh)
	for {
		select {
		case actions := <-actionsCh:
			response := responses.RespondOk(actions, "Success")
			data, _ := json.Marshal(response)
			w.Write([]byte("event: user-actions-stream\n"))
			bytes, err := fmt.Fprintf(w, "data: %s\n\n", data)
			if err != nil {
				errCh <- err
				continue
			}
			log.Printf("Bytes written: %d\n", bytes)
			flusher.Flush()
		case err := <-errCh:
			errorResp := responses.RespondInternalServerError(nil, err.Error())
			data, _ := json.Marshal(errorResp)
			w.Write([]byte("event: user-actions-error\n"))
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
			return
		case <-r.Context().Done():
			endResp := responses.RespondOk(nil, "Stream ended")
			data, _ := json.Marshal(endResp)
			w.Write([]byte("event: user-actions-end\n"))
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
			return
		}
	}
}

func (h *UserActionsHandler) streamActionsWorker(
	ctx context.Context,
	userID *uuid.UUID,
	out chan []*responses.ActionResponse,
	errCh chan error,
) {
	defer close(out)
	defer close(errCh)
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	// Immediate first event
	h.sendActions(ctx, userID, out, errCh)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.sendActions(ctx, userID, out, errCh)
		}
	}
}
func (h *UserActionsHandler) sendActions(
	ctx context.Context,
	userID *uuid.UUID,
	out chan []*responses.ActionResponse,
	errCh chan error,
) {
	actions, err := h.service.GetAll(ctx, userID)
	if err != nil {
		errCh <- err
		return
	}
	result := make([]*responses.ActionResponse, len(actions))
	for i, a := range actions {
		result[i] = &responses.ActionResponse{
			ActionName: a.Action.String(),
			ID:         a.ID.String(),
		}
	}
	out <- result
}
