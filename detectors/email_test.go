package detectors

import "testing"

func TestEmailDetector(t *testing.T) {
	d := NewEmailDetector()

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Simple Email", "Contact me at john@example.com", 1},
		{"Two Emails", "Ping john@example.com and jane@test.org", 2},
		{"Subdomain Email", "Reach us at support@eu.mail.example.co.uk", 1},
		{"Uppercase", "Send to JOHN@EXAMPLE.IO", 1},
		{"Plus Alias", "Reach out via ops+alerts@service.io", 1},
		{"Invalid Missing TLD", "john@example", 0},
		{"Invalid Missing User", "@example.com", 0},
		{"Invalid Double Dot", "john..doe@example.com", 0},
		{"Noise Without Email", "Hello world", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(d.Scan(tt.input)); got != tt.expected {
				t.Errorf("expected %d matches, got %d", tt.expected, got)
			}
		})
	}
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
