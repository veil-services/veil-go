package detectors

import (
	"sync"
	"testing"
)

func TestCreditCardDetector(t *testing.T) {
	d := NewCreditCardDetector()

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		// Valid Cases (Standard Test Numbers)
		{"Valid Visa", "Card 4111 1111 1111 1111 approved.", 1},
		{"Valid Mastercard", "Payment 5555-5555-5555-4444 processed.", 1},
		{"Valid Amex", "Use 371449635398431 for tests", 1}, // 15 digits
		{"Valid 13 Digits", "Old Visa 4222222222222", 1},   // 13 digits
		{"Multiple Cards", "Primary 4242 4242 4242 4242, backup 5555-5555-5555-4444", 2},

		// Invalid Cases (Luhn Algorithm / Logic)
		{"Invalid Luhn", "4111 1111 1111 1112", 0}, // Checksum fail
		{"Too Short", "4111", 0},
		{"Too Long", "123456789012345678901", 0}, // > 19 digits usually ignored

		// Formatting & Noise Edge Cases
		{"Letters Mixed", "4111a1111b1111c111", 0},                 // Should break sequence
		{"Numeric Noise", "Order #1234567890123456 is invalid", 0}, // Likely fails Luhn
		{"Boundary Start", "4111111111111111 is the card", 1},
		{"Boundary End", "The card is 4111111111111111", 1},
		{"Unicode Noise", "Card 4111 1111 1111 1111 🚀", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := d.Scan(tt.input)
			if got := len(matches); got != tt.expected {
				t.Errorf("input: %q\nexpected %d matches, got %d", tt.input, tt.expected, got)
			}
		})
	}
}

// 2. Concurrency Test (Thread-Safety)
func TestCreditCardDetector_Concurrency(t *testing.T) {
	d := NewCreditCardDetector()
	payload := "Thread safe test for Visa 4111 1111 1111 1111 running."
	concurrency := 100
	var wg sync.WaitGroup

	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			matches := d.Scan(payload)
			if len(matches) != 1 {
				t.Errorf("Concurrent scan failed to find match")
			}
		}()
	}
	wg.Wait()
}

// 3. Fuzz Testing (Go 1.18+ native)
// Run with: go test -fuzz=FuzzCreditCard -fuzztime=10s
func FuzzCreditCardDetector(f *testing.F) {
	d := NewCreditCardDetector()

	// Seed corpus
	f.Add("4111 1111 1111 1111")
	f.Add("5555-5555-5555-4444")
	f.Add("Random text 12345")

	f.Fuzz(func(t *testing.T, orig string) {
		// MUST NOT PANIC
		_ = d.Scan(orig)
	})
}

func TestCreditCardDetector_LongTranscript(t *testing.T) {
	d := NewCreditCardDetector()
	text := `
Agent: Please confirm the cards to whitelist.
Client: Use Visa 4111 1111 1111 1111 for orders, Master 5555-5555-5555-4444 for backup.
Agent: Any virtual cards?
Client: Yes, 371449635398431 for Amex tests only.
`
	if got := len(d.Scan(text)); got != 3 {
		t.Fatalf("expected 3 cards, got %d", got)
	}
}

func BenchmarkCreditCardDetector_LongText(b *testing.B) {
	d := NewCreditCardDetector()
	payload := `
Invoices:
1) Order #1001 card 4111 1111 1111 1111
2) Order #1002 card 5555-5555-5555-4444
3) Order #1003 card 4222 2222 2222 2222
4) Order #1004 card 3782 822463 10005
Please mask all before logging.
`
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.Scan(payload)
	}
}
