package ecosystem

var _ error = (*I18NError)(nil)

type I18NError struct {
	errorCode string
	message   string
}

func (e *I18NError) Error() string {
	return e.message
}

func (e *I18NError) Code() string {
	return e.errorCode
}
