package sphere

import "errors"

// List of errors
var (
	ErrNotAuthorized        = errors.New("not authorized")
	ErrUnsupportedNamespace = errors.New("unsupported namespace")
	ErrNotImplemented       = errors.New("not implemented")
	ErrBadScheme            = errors.New("bad scheme")
	ErrBadStatus            = errors.New("bad status")
	ErrBadRequestMethod     = errors.New("bad method")
	ErrNotSupported         = errors.New("not supported")
)
