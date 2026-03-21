package responses

import (
	"net/http"
	"time"
)

type APIResponse struct {
	Status  int    `json:"-"`
	Data    any    `json:"data"`
	Message string `json:"message"`
	Time    int64  `json:"time"`
}

func NewAPIResponse[T any](status int, data *T, message string) *APIResponse {
	return &APIResponse{
		Status:  status,
		Data:    data,
		Message: message,
		Time:    time.Now().UTC().UnixMilli(),
	}
}

func Respond[T any](status int, data *T, message string) *APIResponse {
	return NewAPIResponse(status, data, message)
}

func RespondOk(data any, message string) *APIResponse {
	return Respond(http.StatusOK, &data, message)
}

func RespondCreated(data any, message string) *APIResponse {
	return Respond(http.StatusCreated, &data, message)
}

func RespondAccepted(data any, message string) *APIResponse {
	return Respond(http.StatusAccepted, &data, message)
}

func RespondNoContent(data any, message string) *APIResponse {
	return Respond(http.StatusNoContent, &data, message)
}

func RespondBadRequest(data any, message string) *APIResponse {
	return Respond(http.StatusBadRequest, &data, message)
}

func RespondUnauthorized(data any, message string) *APIResponse {
	return Respond(http.StatusUnauthorized, &data, message)
}

func RespondForbidden(data any, message string) *APIResponse {
	return Respond(http.StatusForbidden, &data, message)
}

func RespondForbbiden(data any, message string) *APIResponse {
	return Respond(http.StatusForbidden, &data, message)
}

func RespondNotFound(data any, message string) *APIResponse {
	return Respond(http.StatusNotFound, &data, message)
}

func RespondNotAllowed(data any, message string) *APIResponse {
	return Respond(http.StatusMethodNotAllowed, &data, message)
}

func RespondConflict(data any, message string) *APIResponse {
	return Respond(http.StatusConflict, &data, message)
}

func RespondInternalServerError(data any, message string) *APIResponse {
	return Respond(http.StatusInternalServerError, &data, message)
}
