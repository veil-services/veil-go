package veil

import "github.com/veil-services/veil-go/detectors"

// Option define uma função para configurar o Veil.
type Option func(*Config)

// WithEmail habilita a máscara de e-mails.
func WithEmail() Option {
	return func(c *Config) {
		c.MaskEmail = true
	}
}

// WithCPF habilita a máscara de CPFs (Brasil).
func WithCPF() Option {
	return func(c *Config) {
		c.MaskCPF = true
	}
}

// WithCNPJ habilita a máscara de CNPJs (Brasil).
func WithCNPJ() Option {
	return func(c *Config) {
		c.MaskCNPJ = true
	}
}

// WithCreditCard habilita a máscara de cartões de crédito.
func WithCreditCard() Option {
	return func(c *Config) {
		c.MaskCreditCard = true
	}
}

// WithIP habilita a máscara de endereços IPv4.
func WithIP() Option {
	return func(c *Config) {
		c.MaskIP = true
	}
}

// WithPhone habilita a máscara de telefones (Formato Global E.164 começando com +).
func WithPhone() Option {
	return func(c *Config) {
		c.MaskPhone = true
	}
}

// WithUUID habilita a máscara de UUIDs/GUIDs.
func WithUUID() Option {
	return func(c *Config) {
		c.MaskUUID = true
	}
}

// WithConsistentTokenization garante que o mesmo valor original receba o mesmo token
// durante o processo de mascaramento de uma string.
// Ex: "joao@a.com ... joao@a.com" -> "<<EMAIL_1>> ... <<EMAIL_1>>"
func WithConsistentTokenization(enabled bool) Option {
	return func(c *Config) {
		c.ConsistentTokenization = enabled
	}
}

// WithCustomDetector adiciona um detector customizado à lista.
func WithCustomDetector(d detectors.Detector) Option {
	return func(c *Config) {
		c.CustomDetectors = append(c.CustomDetectors, d)
	}
}
