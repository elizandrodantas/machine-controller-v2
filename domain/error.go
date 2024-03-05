package domain

import "errors"

type HttpError struct {
	Message string `json:"message"`
}

type errHttp error

var (
	ErrInternalError errHttp = errors.New("internal error, try again")

	// PARAMETER INVALID
	ErrUsernameOrPasswordInvalid errHttp = errors.New("username or password invalid")
	ErrUserAlreayExist           errHttp = errors.New("username already exists")
	ErrScopeInvalid              errHttp = errors.New("scope is not valid")
	ErrParamIdNotUuid            errHttp = errors.New("id is not a valid UUID")
	ErrPageQueryParameterInvalid errHttp = errors.New("page query parameter is not a valid integer")
	ErrInvalidParameters         errHttp = errors.New("invalid parameters")
	ErrInvalidQuery              errHttp = errors.New("invalid query")

	// ALREADY EXIST
	ErrServiceAlreayExist errHttp = errors.New("service already exists")
	ErrScopeAlreadyExist  errHttp = errors.New("scope already registered for this user")

	// AUTHORIZATION
	ErrOnlyCreatorCanChange errHttp = errors.New("only creator can change")
	ErrOnlyCreatorCanDelete errHttp = errors.New("only creator can delete")
	ErrUnauthorized         errHttp = errors.New("unauthorized")
	ErrServiceUnavailable   errHttp = errors.New("service unavailable")
	ErrUserIsNotActive      errHttp = errors.New("user is not active")
	ErrForbidden            errHttp = errors.New("forbidden")

	// NOT FOUND
	ErrNoteNotFound    errHttp = errors.New("note not found")
	ErrServiceNotFound errHttp = errors.New("service not found")
	ErrUserNotExist    errHttp = errors.New("user not exist")
	ErrMachineNotFound errHttp = errors.New("machine not found")
)

func ErrorHttpMessage(m string) *HttpError {
	return &HttpError{m}
}

func ErrorHttpMessageFromError(err error) *HttpError {
	return &HttpError{err.Error()}
}
