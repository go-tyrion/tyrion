package error

type ErrorCode int

type Error struct {
	code    ErrorCode
	message string
	file    string
	line    int
}

func (e *Error) Error() string {
	return e.message
}

//
func New(message string) *Error {
	e := &Error{
		message: message,
	}

	return e
}

func NewWithCode(code ErrorCode, message string) *Error {
	e := &Error{
		message: message,
		code:    code,
	}

	return e
}

func Wrap(err error) *Error {
	return &Error{
		message: err.Error(),
	}
}

func WrapWithCode(code ErrorCode, err error) *Error {
	return &Error{
		code:    code,
		message: err.Error(),
	}
}
