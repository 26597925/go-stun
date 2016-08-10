package stun

import (
	"io"
)

// Error codes introduced by the RFC 5389 Section 15.6
const (
	CodeTryAlternate     = 300
	CodeBadRequest       = 400
	CodeUnauthorized     = 401
	CodeUnknownAttribute = 420
	CodeStaleNonce       = 438
	CodeServerError      = 500
)

// Error codes introduced by the RFC 3489 Section 11.2.9 except listed in RFC 5389.
const (
	CodeStaleCredentials      = 430
	CodeIntegrityCheckFailure = 431
	CodeMissingUsername       = 432
	CodeUseTLS                = 433
	CodeGlobalFailure         = 600
)

var errorText = map[int]string{
	CodeTryAlternate:          "Try Alternate",
	CodeBadRequest:            "Bad Request",
	CodeUnauthorized:          "Unauthorized",
	CodeUnknownAttribute:      "Unknown Attribute",
	CodeStaleCredentials:      "Stale Credentials",
	CodeIntegrityCheckFailure: "Integrity Check Failure",
	CodeMissingUsername:       "Missing Username",
	CodeUseTLS:                "Use TLS",
	CodeStaleNonce:            "Stale Nonce",
	CodeServerError:           "Server Error",
	CodeGlobalFailure:         "Global Failure",
}

// ErrorText returns a reason phrase text for the STUN error code. It returns the empty string if the code is unknown.
func ErrorText(code int) string {
	return errorText[code]
}

// Error represents the ERROR-CODE attribute.
type Error struct {
	Code   int
	Reason string
}

// String returns the string form of the error attribute.
func (e *Error) String() string {
	return e.Reason
}

type errorCodec struct{}

func (errorCodec) Encode(m *Message, v interface{}, b []byte) (int, error) {
	var code int
	var reason string
	switch a := v.(type) {
	case Error:
		code, reason = a.Code, a.Reason
	default:
		return DefaultAttrCodec.Encode(m, v, b)
	}
	n := 4 + len(reason)
	if len(b) < n {
		return 0, io.EOF
	}
	b[0] = 0
	b[1] = 0
	b[2] = byte(code / 100)
	b[3] = byte(code % 100)
	copy(b[4:], reason)
	return n, nil
}

func (errorCodec) Decode(m *Message, b []byte) (interface{}, error) {
	if len(b) < 4 {
		return nil, io.EOF
	}
	code := int(b[2])*100 + int(b[3])
	return &Error{code, string(b[4:])}, nil
}
