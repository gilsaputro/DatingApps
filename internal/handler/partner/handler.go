package partner

import (
	"gilsaputro/dating-apps/internal/service/partner"
)

// PartnerHandler list dependencies for partner handler
type PartnerHandler struct {
	service      partner.PartnerServiceMethod
	timeoutInSec int
}

// Option set options for http handler config
type Option func(*PartnerHandler)

const (
	defaultTimeout = 5
)

// NewPartnerHandler is func to create http partner handler
func NewPartnerHandler(service partner.PartnerServiceMethod, options ...Option) *PartnerHandler {
	handler := &PartnerHandler{
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
		func(h *PartnerHandler) {
			if timeoutinsec <= 0 {
				timeoutinsec = defaultTimeout
			}
			h.timeoutInSec = timeoutinsec
		})
}
