package service

type AppError struct {
	Status  int    `json:"-"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func newAppError(status int, message string) *AppError {
	return &AppError{Status: status, Message: message}
}
