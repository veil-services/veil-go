package veil

import "github.com/veil-services/veil-go/detectors"

// Option defines a function to configure Veil.
type Option func(*Config)

// WithEmail enables masking of emails.
func WithEmail() Option {
	return func(c *Config) {
		c.MaskEmail = true
	}
}

// WithCPF enables masking of CPFs (Brazil).
func WithCPF() Option {
	return func(c *Config) {
		c.MaskCPF = true
	}
}

// WithCNPJ enables masking of CNPJs (Brazil).
func WithCNPJ() Option {
	return func(c *Config) {
		c.MaskCNPJ = true
	}
}

// WithCreditCard enables masking of credit cards.
func WithCreditCard() Option {
	return func(c *Config) {
		c.MaskCreditCard = true
	}
}

// WithIP enables masking of IPv4 addresses.
func WithIP() Option {
	return func(c *Config) {
		c.MaskIP = true
	}
}

// WithPhone enables masking of phone numbers (Global E.164 format starting with +).
func WithPhone() Option {
	return func(c *Config) {
		c.MaskPhone = true
	}
}

// WithUUID enables masking of UUIDs/GUIDs.
func WithUUID() Option {
	return func(c *Config) {
		c.MaskUUID = true
	}
}

// WithConsistentTokenization ensures that the same original value receives the same token
// during the masking process of a single string.
// e.g. "john@a.com ... john@a.com" -> "<<EMAIL_1>> ... <<EMAIL_1>>"
func WithConsistentTokenization(enabled bool) Option {
	return func(c *Config) {
		c.ConsistentTokenization = enabled
	}
}

// WithCustomDetector adds a user-defined detector to the list.
func WithCustomDetector(d detectors.Detector) Option {
	return func(c *Config) {
		c.CustomDetectors = append(c.CustomDetectors, d)
	}
}
