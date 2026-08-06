package model

import (
	
)

const (
	ErrCodeInvalidArgument       = "invalid_argument"
	ErrCodeProfileNotFound       = "profile_not_found"
	ErrCodeRecapNotFound         = "recap_not_found"
	ErrCodeInsufficientActivity  = "insufficient_activity"
	ErrCodeRateLimitExceeded     = "rate_limit_exceeded"
	ErrCodeDependencyUnavailable = "dependency_unavailable"
	ErrCodeInternalError         = "internal_error"
)

// APIError — стандартный ответ с ошибкой (соответствует OpenAPI)
type APIError struct {
	Code      string        `json:"code"`              // код ошибки (например, "invalid_argument")
	Message   string        `json:"message"`           // человекочитаемое сообщение
	RequestID string        `json:"request_id"`        // UUID запроса для трассировки
	Details   []ErrorDetail `json:"details,omitempty"` // детали (если нужно)
}

// ErrorDetail — дополнительная информация об ошибке (поле + причина)
type ErrorDetail struct {
	Field  string `json:"field,omitempty"` // поле, в котором ошибка (если применимо)
	Reason string `json:"reason"`          // причина ошибки
}
