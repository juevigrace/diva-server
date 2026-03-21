package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type VerificationHandler struct {
	service *service.VerificationService
}

func NewVerificationHandler(svc *service.VerificationService) *VerificationHandler {
	return &VerificationHandler{service: svc}
}

func (h *VerificationHandler) Routes(r chi.Router) {
	r.Post("/verify", h.Verify)
}

func (h *VerificationHandler) Verify(w http.ResponseWriter, r *http.Request) {
	var dto dtos.EmailTokenDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	_, err := h.service.Verify(r.Context(), dto.Token)
	if err != nil {
		res := new(responses.APIResponse)
		if errors.Is(err, models.ErrTokenInvalid) {
			res = responses.RespondNotFound(nil, err.Error())
		} else if errors.Is(err, models.ErrTokenExpired) {
			res = responses.RespondBadRequest(nil, err.Error())
		} else {
			res = responses.RespondInternalServerError(nil, err.Error())
		}
		responses.WriteJSON(w, res)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Success"))
}
