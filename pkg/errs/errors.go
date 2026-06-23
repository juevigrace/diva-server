package errs

import "errors"

var (
	ErrContextIsMissing error = errors.New("request context missing")
	ErrHeaderNotValid   error = errors.New("header format not valid")
	ErrNotAuthorized    error = errors.New("you're not authorized to access this endpoint")
	ErrForbidden        error = errors.New("forbidden resource")
	ErrIDRequired       error = errors.New("id param is required")
	ErrParamRequired    error = errors.New("this url param is required")
)

var (
	ErrPermissionDenied           = errors.New("permission denied")
	ErrPermissionExpired          = errors.New("permission expired")
	ErrPermissionExpiration       = errors.New("invalid permission expiration for this action")
	ErrPermissionNotFound   error = errors.New("permission not found")
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUsernameTaken      = errors.New("username already taken")
	ErrEmailTaken         = errors.New("email already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSamePassword       = errors.New("passwords are the same")
)

var (
	ErrTokenExpired      = errors.New("token expired")
	ErrTokenInvalid      = errors.New("token invalid")
	ErrActionNotFound    = errors.New("action not found")
	ErrActionNotVerified = errors.New("action not verified")
)

var (
	ErrTokenNotValid   = errors.New("token is not valid")
	ErrBadAudience     = errors.New("bad audience in token")
	ErrBadIssuer       = errors.New("bad issuer in token")
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionInvalid  = errors.New("invalid session")
)

var (
	ErrAdminAccessRequired = errors.New("admin access required")
)

var (
	ErrFileRequired     = errors.New("file is required")
	ErrFilePathRequired = errors.New("file path is required")
	ErrUnsupportedImage = errors.New("unsupported image type")
	ErrFileTooLarge     = errors.New("file too large")
)
