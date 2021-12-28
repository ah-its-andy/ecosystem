package http

var _ error = (*I18NError)(nil)

type I18NError struct {
	errorCode string
	message   string
}

func (e *I18NError) Error() string {
	return e.errorCode
}

func (e *I18NError) Message() string {
	return e.message
}

func NewI18NError(code string, msg string) error {
	return &I18NError{
		errorCode: code,
		message:   msg,
	}
}
