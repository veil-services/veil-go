package detectors

import (
	"sync"
	"testing"
)

func TestEmailDetector(t *testing.T) {
	d := NewEmailDetector()

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		// Valid Cases
		{"Simple Email", "Contact me at john@example.com", 1},
		{"Two Emails", "Ping john@example.com and jane@test.org", 2},
		{"Subdomain Email", "Reach us at support@eu.mail.example.co.uk", 1},
		{"Uppercase", "Send to JOHN@EXAMPLE.IO", 1},
		{"Plus Alias", "Reach out via ops+alerts@service.io", 1},
		{"Numeric User", "123456@numbers.com", 1},
		{"Dot in User", "firstname.lastname@company.com", 1},
		{"Dash in Domain", "info@my-company.net", 1},

		// Invalid Cases
		{"Missing TLD", "john@example", 0}, // Debateable, but usually invalid in public web
		{"Missing User", "@example.com", 0},
		{"Missing At", "john.example.com", 0},
		{"Double Dot", "john..doe@example.com", 0},
		{"Spaces", "john @ example . com", 0},

		// Noise & Boundary
		{"Noise Without Email", "Hello world this is a test", 0},
		{"Embedded in Text", "The email is<john@example.com>.", 1}, // Should extract just the email
		{"Punctuation End", "Write to abuse@twitter.com.", 1},      // Should exclude the trailing dot
		{"Unicode Noise", "Email: john@example.com 🚀", 1},
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
func TestEmailDetector_Concurrency(t *testing.T) {
	d := NewEmailDetector()
	payload := "Thread safe test for admin@veil.services running."
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
// Run with: go test -fuzz=FuzzEmail -fuzztime=10s
func FuzzEmailDetector(f *testing.F) {
	d := NewEmailDetector()

	// Seed corpus
	f.Add("test@example.com")
	f.Add("user+tag@sub.domain.co.uk")
	f.Add("not an email")

	f.Fuzz(func(t *testing.T, orig string) {
		// MUST NOT PANIC
		_ = d.Scan(orig)
	})
}

func TestEmailDetector_LongText(t *testing.T) {
	d := NewEmailDetector()

	input := `
Hello team,

Please invite alice.wong@corp.com and the EU contact ops-eu@service.co.uk to tomorrow's sync.
For escalations, cc support+vip@help.io.

Thanks!
`

	matches := d.Scan(input)
	if len(matches) != 3 {
		t.Fatalf("expected 3 emails, got %d (%v)", len(matches), matches)
	}
}

func BenchmarkEmailDetector_LongText(b *testing.B) {
	d := NewEmailDetector()
	payload := `
Meeting notes:
- Owner: maria@product.io
- Backup: john.smith@ops.acme.com
- Escalations: pager-duty@alerts.west.cloud
- Vendors: contact@vendor-one.io, billing@vendor-two.net
No other addresses should be captured.
`
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.Scan(payload)
	}
}
