package detectors

import (
	"testing"
)

func TestEmailDetector(t *testing.T) {
	d := NewEmailDetector()

	tests := []struct {
		name     string
		input    string
		expected int // number of matches
	}{
		{"Simple Email", "Contact me at john@example.com", 1},
		{"Two Emails", "john@example.com and jane@test.org", 2},
		{"No Email", "Just some text", 0},
		{"Invalid Email", "john@example (incomplete)", 0},
		{"Complex Email", "john.doe+tag@sub.example.co.uk", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := d.Scan(tt.input)
			if len(matches) != tt.expected {
				t.Errorf("expected %d matches, got %d", tt.expected, len(matches))
			}
		})
	}
}

func TestCPFDetector(t *testing.T) {
	d := NewCPFDetector()

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		// 123.456.789-09 é matematicamente válido! Trocando para final 88 que certamente falha.
		{"Invalid Checksum CPF", "CPF: 123.456.789-88", 0}, 
		// Calculado manualmente: 111.444.777-35
		{"Valid CPF Manual", "CPF: 111.444.777-35", 1},
		{"Valid CPF Clean", "11144477735", 1},
		{"Invalid Checksum", "111.222.333-44", 0},
		{"Repeated Digits", "111.111.111-11", 0}, 
		{"Short Number", "123.456", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := d.Scan(tt.input)
			if len(matches) != tt.expected {
				t.Errorf("expected %d matches, got %d for input %s", tt.expected, len(matches), tt.input)
			}
		})
	}
}

func TestCNPJDetector(t *testing.T) {
	d := NewCNPJDetector()

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		// Banco do Brasil: 00.000.000/0001-91
		{"Valid CNPJ BB", "CNPJ: 00.000.000/0001-91", 1},
		// Google Brasil: 06.990.590/0001-23
		{"Valid CNPJ Google", "06.990.590/0001-23", 1},
		{"Invalid Checksum", "00.000.000/0001-00", 0},
		{"Short Number", "12.345", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := d.Scan(tt.input)
			if len(matches) != tt.expected {
				t.Errorf("expected %d matches, got %d for input %s", tt.expected, len(matches), tt.input)
			}
		})
	}
}

func TestCreditCardDetector(t *testing.T) {
	d := NewCreditCardDetector()

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Valid Visa", "4111 1111 1111 1111", 1},
		{"Valid Mastercard", "5555 5555 5555 4444", 1},
		{"Invalid Luhn", "4111 1111 1111 1112", 0},
		{"Too Short", "4111", 0},
		{"With Dashes", "4111-1111-1111-1111", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := d.Scan(tt.input)
			if len(matches) != tt.expected {
				t.Errorf("expected %d matches, got %d", tt.expected, len(matches))
			}
		})
	}
}
