package detectors

import (
	"sync"
	"testing"
)

func TestPhoneDetector(t *testing.T) {
	d := NewPhoneDetector()

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		// Valid E.164 & International Formats
		{"E164 US", "Call +1 555 010 9999", 1},
		{"E164 BR", "Contact +55 11 99999-9999", 1},
		{"With Dots", "Pager +33.1.23.45.67.89", 1}, // Common in Europe
		{"With Spaces", "Pager +33 1 23 45 67 89", 1},
		{"International No DDD", "Emergency +44 20 7946 0321", 1},
		{"Embedded In Text", "Meet me (+81 90 1234 5678) tomorrow", 1},
		{"Multiple Numbers", "List: +49 30 123456, +49 160 9876543", 2},

		// Invalid Cases (Strict Mode)
		{"Trailing Digit Reject", "Token +1234567890123456 extra digit", 0}, // Max 15 digits (E.164 standard)
		{"No Plus Prefix", "Dial 555-010-9999", 0},                          // Ignored in strict global detector to avoid false positives on math/dates
		{"Too Short", "+1234", 0},                                           // Usually < 7-8 digits is not a full international number
		{"Too Long", "+12345678901234567", 0},
		{"Continuous Digits Past Limit", "+1234567890123456", 0}, // 16 digits
		{"Separators Before Digit", "Call +-1234567", 0},
		{"Only Plus Sign", "+ just a plus", 0},
		{"Math Expression", "Value is +1 500", 0}, // Looks like number but too short/simple

		// Noise & Boundaries
		{"Unicode Noise", "Phone: +1 555 010 9999 🚀", 1},
		{"Boundary Start", "+55 11 99999-9999 is the number", 1},
		{"Boundary End", "Call +55 11 99999-9999", 1},
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
func TestPhoneDetector_Concurrency(t *testing.T) {
	d := NewPhoneDetector()
	payload := "Support hotline is +1 800 555 0199 open 24/7."
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
// Run with: go test -fuzz=FuzzPhone -fuzztime=10s
func FuzzPhoneDetector(f *testing.F) {
	d := NewPhoneDetector()

	// Seed corpus
	f.Add("+1 555 010 9999")
	f.Add("+55 11 99999-9999")
	f.Add("random text +123")

	f.Fuzz(func(t *testing.T, orig string) {
		// MUST NOT PANIC
		_ = d.Scan(orig)
	})
}

func TestPhoneDetector_LongConversation(t *testing.T) {
	d := NewPhoneDetector()

	input := `
		Agent: Hello! To verify, can you confirm the numbers we have on file?
		Client: Sure, my main line is +1 650 555 1234. At the office they reach me at +44 20 7946 0321.
		Agent: Perfect. Do you have any local numbers for delivery teams?
		Client: Yes, use +55 11 4002-8922 for São Paulo and +55 21 3333-2222 for Rio.
		Agent: Thanks! Anything else?
		Client: If it's urgent after hours, try my backup +81-90-1234-5678.
	`

	matches := d.Scan(input)
	if len(matches) != 5 {
		t.Fatalf("expected 5 numbers, got %d", len(matches))
	}
}

func BenchmarkPhoneDetector_LongPrompt(b *testing.B) {
	d := NewPhoneDetector()
	input := `
		Hello support, I am listing all escalation contacts:
		- North America: +1 212 555 0198
		- Backup line: +1 312 555 0110
		- EU Ops: +49 30 1234567 and +49 170 5555555
		- UK Ops: +44 20 1234 5678
		- APAC Ops: +61 2 9876 5432
		- Japan Ops: +81 90 1234 5678
		- Brazil Ops: +55 11 4002-8922
		- Emergency: +1 999 888 7777 (call anytime)
	`

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.Scan(input)
	}
}
