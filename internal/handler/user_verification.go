package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type UserVerificationHandler struct {
	sService *service.SessionService
	uService *service.UserService
	vService *service.UserVerificationService
}

func NewVerificationHandler(
	sService *service.SessionService,
	uService *service.UserService,
	vService *service.UserVerificationService,
) *UserVerificationHandler {
	return &UserVerificationHandler{
		sService: sService,
		uService: uService,
		vService: vService,
	}
}

func (h *UserVerificationHandler) Routes(r chi.Router) {
	r.Route("/verification", func(v chi.Router) {
		v.Post("/request", h.requestVerification)
		v.Post("/", h.verify)
	})
}

func (h *UserVerificationHandler) requestVerification(w http.ResponseWriter, r *http.Request) {
	var dto dtos.RequestActionVerificationDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	dbUser, err := h.uService.GetByEmail(r.Context(), dto.Email)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	res, err := h.vService.RequestVerification(r.Context(), dbUser, &dto)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	responses.WriteJSON(w, responses.RespondOk(res.Response(), "check your email"))
}

func (h *UserVerificationHandler) verify(w http.ResponseWriter, r *http.Request) {
	var dto dtos.VerifyActionDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	action, err := h.vService.Verify(r.Context(), &dto)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	switch action.Name {
	case models.ActionPasswordReset:
	case models.ActionUserVerification:
		go func() {
			if err := h.uService.VerifyUser(r.Context(), action.ID); err != nil {
				return
			}
		}()
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Success"))
}
