package postgql

type PublicError struct {
	message string
}

type PrivateError struct {
	message string
}

//newPublicError returns new error that will be shown to client
func newPublicError(msg string) *PublicError {
	return &PublicError{message: msg}
}

//newPublicError returns new error that won`t be shown to client
func newPrivateError(msg string) *PrivateError {
	return &PrivateError{message: msg}
}

//implementing built-in interface error
func (err *PublicError) Error() string {
	return err.message
}

//implementing built-in interface error
func (err *PrivateError) Error() string {
	return err.message
}
