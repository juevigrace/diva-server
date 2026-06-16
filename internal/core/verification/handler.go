package verification

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
)

type VerificationHandler struct {
	vRepo *VerificationRepo
}

func NewVerificationHandler(vRepo *VerificationRepo) *VerificationHandler {
	return &VerificationHandler{
		vRepo: vRepo,
	}
}

func (h *VerificationHandler) requestVerification(w http.ResponseWriter, r *http.Request) {
	var dto dtos.RequestActionVerificationDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	res, err := h.vRepo.RequestVerification(r.Context(), dto.Email, dto.Action)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(res.Response(), "check your email"))
}

func (h *VerificationHandler) verify(w http.ResponseWriter, r *http.Request) {
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

	if err := h.vRepo.Verify(r.Context(), actionID, dto.Token); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Success"))
}
