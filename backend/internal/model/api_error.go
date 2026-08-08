package model

const (
	ErrCodeInvalidArgument         = "invalid_argument"
	ErrCodeProfileNotFound         = "profile_not_found"
	ErrCodeRecapNotFound           = "recap_not_found"
	ErrCodeExplanationNotAvailable = "explanation_not_available"
	ErrCodeShareNotAvailable       = "share_not_available"
	ErrCodeInsufficientActivity    = "insufficient_activity"
	ErrCodeRateLimitExceeded       = "rate_limit_exceeded"
	ErrCodeDependencyUnavailable   = "dependency_unavailable"
	ErrCodeInternalError           = "internal_error"
)

type APIError struct {
	Code      string        `json:"code"`
	Message   string        `json:"message"`
	RequestID string        `json:"request_id"`
	Details   []ErrorDetail `json:"details,omitempty"`
}

type ErrorDetail struct {
	Field  string `json:"field,omitempty"`
	Reason string `json:"reason"`
}
