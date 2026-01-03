package haozpay

import "fmt"

type SDKError struct {
	Code       int
	Message    string
	RequestID  string
	StatusCode int
}

func (e *SDKError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("[%d] %s (RequestID: %s, StatusCode: %d)", 
			e.Code, e.Message, e.RequestID, e.StatusCode)
	}
	return fmt.Sprintf("[%d] %s (StatusCode: %d)", e.Code, e.Message, e.StatusCode)
}

func NewSDKError(code int, message string, statusCode int) *SDKError {
	return &SDKError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

func NewSDKErrorWithRequestID(code int, message string, statusCode int, requestID string) *SDKError {
	return &SDKError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		RequestID:  requestID,
	}
}

type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("config error: %s - %s", e.Field, e.Message)
}

func ErrInvalidConfig(message string) *ConfigError {
	return &ConfigError{
		Field:   "config",
		Message: message,
	}
}

var (
	ErrTimeout         = NewSDKError(1001, "request timeout", 0)
	ErrNetworkError    = NewSDKError(1002, "network error", 0)
	ErrInvalidResponse = NewSDKError(1003, "invalid response", 0)
	ErrInvalidRequest  = NewSDKError(1004, "invalid request", 0)
	ErrUnauthorized    = NewSDKError(1005, "unauthorized", 401)
	ErrForbidden       = NewSDKError(1006, "forbidden", 403)
	ErrNotFound        = NewSDKError(1007, "not found", 404)
	ErrServerError     = NewSDKError(1008, "server error", 500)
)