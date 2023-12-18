package authentication

import (
	"gilsaputro/dating-apps/internal/service/authentication"
)

// AuthenticationHandler list dependencies for authentication handler
type AuthenticationHandler struct {
	service      authentication.AuthenticationServiceMethod
	timeoutInSec int
}

// Option set options for http handler config
type Option func(*AuthenticationHandler)

const (
	defaultTimeout = 5
)

// NewAuthenticationHandler is func to create http auth handler
func NewAuthenticationHandler(service authentication.AuthenticationServiceMethod, options ...Option) *AuthenticationHandler {
	handler := &AuthenticationHandler{
		service:      service,
		timeoutInSec: defaultTimeout,
	}

	// Apply options
	for _, opt := range options {
		opt(handler)
	}

	return handler
}

// WithTimeoutOptions is func to set timeout config into handler
func WithTimeoutOptions(timeoutinsec int) Option {
	return Option(
		func(h *AuthenticationHandler) {
			if timeoutinsec <= 0 {
				timeoutinsec = defaultTimeout
			}
			h.timeoutInSec = timeoutinsec
		})
}
