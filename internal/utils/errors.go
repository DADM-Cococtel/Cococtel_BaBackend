package utils

type ApiError interface {
	Message() error
	Status() int
}

type apiErr struct {
	ErrorMessage error `json:"message"`
	ErrorStatus  int   `json:"status"`
}

func (e *apiErr) Message() error {
	return e.ErrorMessage
}

func (e *apiErr) Status() int {
	return e.ErrorStatus
}

func newApiError(message error, status int) *apiErr {
	return &apiErr{
		ErrorMessage: message,
		ErrorStatus:  status,
	}
}

func NewApiError(message error, status int) ApiError {
	return newApiError(message, status)
}
