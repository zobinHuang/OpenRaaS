package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

type Type string

const (
	Authorization        Type = "AUTHORIZATION"          // Authentication Failures
	BadRequest           Type = "BAD_REQUEST"            // Validation errors
	Conflict             Type = "CONFLICT"               // Already exists - 409
	Internal             Type = "INTERNAL"               // Server (500) and fallback errors
	NotFound             Type = "NOTFOUND"               // For not finding resource
	PayloadTooLarge      Type = "PAYLOADTOOLARGE"        // for uploading tons of data over the limit - 413
	ServiceUnavailable   Type = "SERVICE_UNAVALIABLE"    // for long running handlers
	UnsupportedMediaType Type = "UNSUPPORTED_MEDIA_TYPE" // for http 415
)

type Error struct {
	Type    Type                   `json:"type"`
	Message string                 `json:"message"`
	Content map[string]interface{} `json:"content"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Status() int {
	switch e.Type {
	case Authorization:
		return http.StatusUnauthorized
	case BadRequest:
		return http.StatusBadRequest
	case Conflict:
		return http.StatusConflict
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	case PayloadTooLarge:
		return http.StatusRequestEntityTooLarge
	case ServiceUnavailable:
		return http.StatusServiceUnavailable
	case UnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	default:
		return http.StatusInternalServerError
	}
}

func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.Status()
	}
	return http.StatusInternalServerError
}

func NewAuthorization(reason string) *Error {
	return &Error{
		Type:    Authorization,
		Message: reason,
	}
}

func NewBadRequest(reason string) *Error {
	return &Error{
		Type:    BadRequest,
		Message: reason,
	}
}

func NewConflict(name string, value string) *Error {
	content := make(map[string]interface{})

	content["conflict_type"] = name
	content["conflict_value"] = value

	return &Error{
		Type:    Conflict,
		Message: fmt.Sprintf("resource: %v with value: %v already exists", name, value),
		Content: content,
	}
}

func NewInternal() *Error {
	return &Error{
		Type:    Internal,
		Message: fmt.Sprintf("internal server error"),
	}
}

func NewNotFound(name, value string) *Error {
	content := make(map[string]interface{})

	content["notfound_type"] = name
	content["notfound_value"] = value

	return &Error{
		Type:    NotFound,
		Message: fmt.Sprintf("resource: %v with value %v not found", name, value),
		Content: content,
	}
}

func NewPayloadTooLarge(maxBodySize, contentLength int64) *Error {
	return &Error{
		Type:    PayloadTooLarge,
		Message: fmt.Sprintf("May payload size of %v exceeded, Actual payload size: %v", maxBodySize, contentLength),
	}
}

func NewServiceUnavailable() *Error {
	return &Error{
		Type:    ServiceUnavailable,
		Message: fmt.Sprintf("Service unavailable or timed out"),
	}
}

func NewUnsupportedMediaType(reason string) *Error {
	return &Error{
		Type:    UnsupportedMediaType,
		Message: reason,
	}
}
