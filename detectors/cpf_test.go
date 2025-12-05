package detectors

import (
	"sync"
	"testing"
)

func TestCPFDetector(t *testing.T) {
	d := NewCPFDetector()

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		// Valid Cases
		{"Valid Formatted", "Client CPF: 111.444.777-35 confirmed.", 1},
		{"Valid Plain", "Document: 11144477735", 1},
		{"Valid Embedded", "ID: abc111.444.777-35xyz", 1},
		{"Second Valid", "Data: 529.982.247-25", 1},

		// Invalid Cases (Business Logic / Checksum)
		{"Invalid Checksum", "CPF 123.456.789-00", 0}, // Checksum math fail
		{"All Equals", "111.111.111-11", 0},           // Invalid by definition (blacklist)
		{"All Zeros", "000.000.000-00", 0},            // Invalid by definition

		// Formatting & Noise Edge Cases
		{"Short Sequence", "123.456", 0},
		{"Almost CPF", "111.444.777-3", 0},        // Missing last digit
		{"Wrong Separators", "111/444/777.35", 0}, // Wrong separator chars (slash instead of dot)
		{"Mixed Separators", "111.444.777.35", 0}, // Dot at the end instead of dash (strict check?)
		{"Numeric Noise", "Order total 12345-67890", 0},
		{"Unicode Noise", "CPF 111.444.777-35 🚀", 1}, // Ensure unicode safety
		{"Boundary Start", "111.444.777-35 is the ID", 1},
		{"Boundary End", "The ID is 111.444.777-35", 1},
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
func TestCPFDetector_Concurrency(t *testing.T) {
	d := NewCPFDetector()
	payload := "Thread safe test for CPF 529.982.247-25 running."
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
// Run with: go test -fuzz=FuzzCPF -fuzztime=10s
func FuzzCPFDetector(f *testing.F) {
	d := NewCPFDetector()

	// Seed corpus
	f.Add("111.444.777-35")
	f.Add("12345678901")
	f.Add("Random text with numbers 12345")

	f.Fuzz(func(t *testing.T, orig string) {
		// MUST NOT PANIC
		_ = d.Scan(orig)
	})
}

func TestCPFDetector_LongParagraph(t *testing.T) {
	d := NewCPFDetector()
	text := `
Registry:
 - Client 1: 111.444.777-35
 - Client 2: CPF 52998224725
 - Client 3: masked 000.000.000-00 (invalid, should be ignored)
`
	if got := len(d.Scan(text)); got != 2 {
		t.Fatalf("expected 2 valid CPFs, got %d", got)
	}
}

func BenchmarkCPFDetector_LongText(b *testing.B) {
	d := NewCPFDetector()
	payload := `
Report Data:
Client A CPF 111.444.777-35
Client B CPF 529.982.247-25
Client C CPF 862.883.667-57
`
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.Scan(payload)
	}
}
