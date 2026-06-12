package verification

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
)

type UserVerificationHandler struct {
	vService *UserVerificationService
}

func NewVerificationHandler(
	vService *UserVerificationService,
) *UserVerificationHandler {
	return &UserVerificationHandler{
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

	res, err := h.vService.RequestVerification(r.Context(), dto.Email, dto.Action)
	if err != nil {
		responses.HandleReqError(w, err)
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

	actionID, err := uuid.Parse(dto.ActionID)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.vService.Verify(r.Context(), actionID, dto.Token); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Success"))
}
