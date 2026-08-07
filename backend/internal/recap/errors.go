package recap

import (
	"errors"
	"fmt"
)

var ErrInsufficientActivity = errors.New("insufficient_activity")

type InputError struct {
	Code    string
	Message string
}

func (e *InputError) Error() string {
	if e.Message == "" {
		return e.Code
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

type ConfigError struct {
	Code    string
	Message string
}

func (e *ConfigError) Error() string {
	if e.Message == "" {
		return e.Code
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
