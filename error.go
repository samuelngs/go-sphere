package sphere

import (
	"log"
	"runtime"
)

// List of errors
var (
	ErrNotFound         = &ProtocolError{"not found"}
	ErrNotSupported     = &ProtocolError{"not supported"}
	ErrNotImplemented   = &ProtocolError{"not implemented"}
	ErrTooManyRequest   = &ProtocolError{"too many requests"}
	ErrBadScheme        = &ProtocolError{"bad scheme"}
	ErrBadStatus        = &ProtocolError{"bad status"}
	ErrBadRequestMethod = &ProtocolError{"bad method"}
	ErrUnauthorized     = &ProtocolError{"unauthorized"}
	ErrServerErrors     = &ProtocolError{"server errors"}
	ErrRequestFailed    = &ProtocolError{"request failed"}

	ErrAlreadySubscribed = &ClientError{"already subscribed"}
	ErrNotSubscribed     = &ClientError{"not subscribed"}

	ErrPacketBadScheme = &PacketError{"packet bad scheme"}
	ErrPacketBadType   = &PacketError{"packet bad type"}
)

// IError interface
type IError interface {
	Error() string
}

// Error is a trivial implementation of error.
type Error struct {
	s string
}

// Error returns error string of Error
func (e *Error) Error() string {
	return e.s
}

// ProtocolError represents WebSocket protocol errors.
type ProtocolError Error

// Error returns error string of ProtocolError
func (e *ProtocolError) Error() string {
	return e.s
}

// ClientError represents general client errors.
type ClientError Error

// Error returns error string of ChannelError
func (e *ClientError) Error() string {
	return e.s
}

// PacketError represents general client errors.
type PacketError Error

// Error returns error string of ChannelError
func (e *PacketError) Error() string {
	return e.s
}

// LogError logs the function name, line and error message
func LogError(err IError) {
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		pc, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
	}
}
