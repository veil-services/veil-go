package detectors

import "testing"

func TestCreditCardDetector(t *testing.T) {
	d := NewCreditCardDetector()

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Valid Visa", "Card 4111 1111 1111 1111 aprovado.", 1},
		{"Valid Mastercard", "Cartao 5555-5555-5555-4444", 1},
		{"Valid 13 digits", "Travel card 4222222222222", 1},
		{"Valid Amex", "Use 371449635398431 for tests", 1},
		{"Multiple Cards", "Primary 4242 4242 4242 4242, backup 5555-5555-5555-4444", 2},
		{"Invalid Luhn", "4111 1111 1111 1112", 0},
		{"Too Short", "4111", 0},
		{"Letters Mixed", "4111a1111b1111c111", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(d.Scan(tt.input)); got != tt.expected {
				t.Errorf("expected %d matches, got %d", tt.expected, got)
			}
		})
	}
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

