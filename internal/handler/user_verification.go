package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juevigrace/diva-server/internal/middlewares"
	"github.com/juevigrace/diva-server/internal/models"
	"github.com/juevigrace/diva-server/internal/models/dtos"
	"github.com/juevigrace/diva-server/internal/models/errs"
	"github.com/juevigrace/diva-server/internal/models/responses"
	"github.com/juevigrace/diva-server/internal/service"
)

type UserVerificationHandler struct {
	sService  *service.SessionService
	uService  *service.UserService
	uaService *service.UserActionsService
	vService  *service.UserVerificationService
	upService *service.UserPermissionService
}

func NewVerificationHandler(
	sService *service.SessionService,
	uService *service.UserService,
	uaService *service.UserActionsService,
	vService *service.UserVerificationService,
	upService *service.UserPermissionService,
) *UserVerificationHandler {
	return &UserVerificationHandler{
		sService:  sService,
		uService:  uService,
		uaService: uaService,
		vService:  vService,
		upService: upService,
	}
}

// TODO: create verification actions to user data updates
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

	parsedAction := models.ActionFromString(dto.Action)
	if parsedAction == -1 {
		responses.WriteJSON(w, responses.RespondBadRequest(nil, errs.ErrActionNotFound.Error()))
		return
	}

	dbUser, err := h.uService.GetByEmail(r.Context(), dto.Email)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	res, err := h.vService.RequestVerification(r.Context(), dbUser, parsedAction)
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

	va, err := h.vService.Verify(r.Context(), actionID, dto.Token)
	if err != nil {
		responses.HandleReqError(w, err)
		return
	}

	switch va.Action.Name {
	case models.ActionPasswordUpdate:
	case models.ActionUserVerification:
		if !va.Verified {
			responses.WriteJSON(w, responses.RespondForbbiden(nil, errs.ErrActionNotVerified.Error()))
			return
		}

		if err := h.uService.UpdateVerified(r.Context(), true, va.Action.UserID); err != nil {
			responses.HandleReqError(w, err)
			return
		}

		if err := h.uaService.Delete(r.Context(), va.Action.ID); err != nil {
			responses.HandleReqError(w, err)
			return
		}
	case models.ActionEmailUpdate, models.ActionUsernameUpdate, models.ActionPhoneUpdate:
		if !va.Verified {
			responses.WriteJSON(w, responses.RespondForbbiden(nil, errs.ErrActionNotVerified.Error()))
			return
		}

		var permAction models.PermissionAction
		switch va.Action.Name {
		case models.ActionEmailUpdate:
			permAction = models.PERMISSION_USERS_EMAIL_WRITE
		case models.ActionUsernameUpdate:
			permAction = models.PERMISSION_USERS_USERNAME_WRITE
		case models.ActionPhoneUpdate:
			permAction = models.PERMISSION_USERS_PHONE_WRITE
		}

		exp := time.Now().UTC().Add(15 * time.Minute).UnixMilli()
		if err := h.upService.CreateByName(r.Context(), permAction, nil, true, &exp, va.Action.UserID); err != nil {
			responses.HandleReqError(w, err)
			return
		}

		if err := h.uaService.Delete(r.Context(), va.Action.ID); err != nil {
			responses.HandleReqError(w, err)
			return
		}
	}

	responses.WriteJSON(w, responses.RespondOk(nil, "Success"))
}
