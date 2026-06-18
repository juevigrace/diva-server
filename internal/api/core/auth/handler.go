package auth

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/api/middlewares"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/responses"
)

type AuthHandler struct {
	aRepo *AuthRepo
}

func NewAuthHandler(aRepo *AuthRepo) *AuthHandler {
	return &AuthHandler{
		aRepo: aRepo,
	}
}

func (h *AuthHandler) signIn(w http.ResponseWriter, r *http.Request) {
	var dto dtos.SignInDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}
	dto.SessionData.IpAddress = strings.Split(r.RemoteAddr, ":")[0]

	session, err := h.aRepo.SignIn(r.Context(), &dto)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(session.Response(), "sign in successful"))
}

func (h *AuthHandler) signUp(w http.ResponseWriter, r *http.Request) {
	var dto dtos.SignUpDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}
	dto.SessionData.IpAddress = strings.Split(r.RemoteAddr, ":")[0]

	session, err := h.aRepo.SignUp(r.Context(), &dto)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondCreated(session.Response(), "sign up successful"))
}

func (h *AuthHandler) signOut(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	var dto dtos.SessionDataDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	if err := h.aRepo.SignOut(r.Context(), session.ID); err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Sign out successful"))
}

func (h *AuthHandler) ping(w http.ResponseWriter, r *http.Request) {
	responses.WriteJSON(w, responses.RespondOk(nil, "Pong"))
}

func (h *AuthHandler) refresh(w http.ResponseWriter, r *http.Request) {
	session, ok := middlewares.GetSessionFromContext(r.Context())
	if !ok {
		responses.WriteJSON(w, responses.RespondUnauthorized(nil, "session not found"))
		return
	}

	var dto dtos.SessionDataDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	res, err := h.aRepo.Refresh(r.Context(), session, &dto)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(res, "Session refreshed"))
}

func (h *AuthHandler) forgotPasswordConfirm(w http.ResponseWriter, r *http.Request) {
	var dto dtos.ForgotPasswordConfirmDto
	if err := middlewares.ValidateBody(&dto, r); err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}
	dto.SessionData.IpAddress = strings.Split(r.RemoteAddr, ":")[0]

	parsedID, err := uuid.Parse(dto.ActionID)
	if err != nil {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, err.Error()))
		return
	}

	session, err := h.aRepo.ForgotPasswordConfirm(r.Context(), parsedID, &dto.SessionData)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	responses.WriteJSON(w, responses.RespondOk(session, "Success"))
}
