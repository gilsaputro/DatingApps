package find

import (
	"gilsaputro/dating-apps/internal/service/find"
)

// FindHandler list dependencies for Find handler
type FindHandler struct {
	service      find.FindServiceMethod
	timeoutInSec int
}

// Option set options for http handler config
type Option func(*FindHandler)

const (
	defaultTimeout = 5
)

// NewFindHandler is func to create http Find handler
func NewFindHandler(service find.FindServiceMethod, options ...Option) *FindHandler {
	handler := &FindHandler{
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
		func(h *FindHandler) {
			if timeoutinsec <= 0 {
				timeoutinsec = defaultTimeout
			}
			h.timeoutInSec = timeoutinsec
		})
}
