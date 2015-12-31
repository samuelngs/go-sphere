package sphere

// Error Error struct
type Error struct {
	ErrorString string
}

// Error returns ErrorString message
func (err *Error) Error() string {
	return err.ErrorString
}

// List of errors
var (
	ErrUnsupportedNamespace = &Error{"unsupported namespace"}
	ErrNotImplemented       = &Error{"not implemented"}
	ErrBadScheme            = &Error{"bad scheme"}
	ErrBadStatus            = &Error{"bad status"}
	ErrBadRequestMethod     = &Error{"bad method"}
	ErrNotSupported         = &Error{"not supported"}
)
