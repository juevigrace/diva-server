package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type VerificationHandler struct {
	sService *service.SessionService
	vService *service.VerificationService
}

func NewVerificationHandler(
	sService *service.SessionService,
	vService *service.VerificationService,
) *VerificationHandler {
	return &VerificationHandler{
		sService: sService,
		vService: vService,
	}
}

func (h *VerificationHandler) Routes(r chi.Router) {
	r.Route("/verification", func(v chi.Router) {
		v.Post("/request", h.requestVerification)
		v.Post("/", h.verify)

		v.Route("/auth", func(auth chi.Router) {
			auth.Use(middlewares.SessionMiddleware(h.sService.GetByID))
			auth.Post("/", h.verifyWithAuth)
		})
	})
}

func (h *VerificationHandler) requestVerification(w http.ResponseWriter, r *http.Request) {
	var dto dtos.RequestVerificationDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.vService.RequestVerification(r.Context(), &dto); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "check your email"))
}

func (h *VerificationHandler) verify(w http.ResponseWriter, r *http.Request) {
	var dto dtos.VerificationDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	uv, err := h.vService.Verify(r.Context(), dto.Token)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}
	defer func() {
		if err := h.vService.Delete(r.Context(), dto.Token); err != nil {
			responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		}
	}()

	switch uv.UserAction.Action {
	case models.ActionPasswordReset:
		if dto.SessionData != nil {
			dto.SessionData.IpAddress = strings.Split(r.RemoteAddr, ":")[0]
		}
		session, err := h.vService.HandlePasswordReset(r.Context(), &uv.UserID, dto.SessionData)
		if err != nil {
			if errors.Is(err, models.ErrTokenInvalid) {
				responses.WriteJSON(w, responses.RespondNotFound(nil, err.Error()))
				return
			}
			responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
			return
		}

		responses.WriteJSON(w, responses.RespondOk(session, "Success"))
	case models.ActionUserVerification:
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "invalid action"))
	}
}

func (h *VerificationHandler) verifyWithAuth(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	var dto dtos.VerificationDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	uv, err := h.vService.Verify(r.Context(), dto.Token)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}
	defer func() {
		if err := h.vService.Delete(r.Context(), dto.Token); err != nil {
			responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		}
	}()

	switch uv.UserAction.Action {
	case models.ActionPasswordReset:
		responses.WriteJSON(w, responses.RespondBadRequest(nil, "invalid action"))
	case models.ActionUserVerification:
		if err := h.vService.HandleVerifyUser(r.Context(), session.User.ID); err != nil {
			responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
			return
		}
		responses.WriteJSON(w, responses.RespondOk(nil, "Success"))
	}
}
