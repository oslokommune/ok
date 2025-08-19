package error_user_msg

type ErrorUserMessage struct {
	err          error
	errorMessage string
	details      string
}

func (e *ErrorUserMessage) Error() string {
	return e.errorMessage
}

func (e *ErrorUserMessage) Unwrap() error {
	return e.err
}

func (e *ErrorUserMessage) Details() string {
	return e.details
}

func NewError(errorMessage string, details string, err error) ErrorUserMessage {
	return ErrorUserMessage{
		errorMessage: errorMessage,
		err:          err,
		details:      details,
	}
}
