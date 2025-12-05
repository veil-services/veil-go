package detectors

import (
	"sync"
	"testing"
)

func TestCNPJDetector(t *testing.T) {
	d := NewCNPJDetector()

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		// Valid Cases
		{"Valid BB", "Company 00.000.000/0001-91 registered.", 1},
		{"Valid Google", "Branch 06.990.590/0001-23", 1},
		{"Valid Plain", "06990590000123", 1},
		{"Valid Embedded", "Text XX00.000.000/0001-91YY", 1},

		// Invalid Cases (Business Logic)
		{"Invalid Checksum", "00.000.000/0001-00", 0},
		{"All Equals", "11.111.111/1111-11", 0}, // Would pass Mod11 if not explicitly checked
		{"All Zeros", "00.000.000/0000-00", 0},

		// Noise and Formatting Cases
		{"Short Sequence", "12.345", 0},
		{"Almost CNPJ", "12.345.678/0001-", 0},        // Missing digits
		{"Wrong Separators", "00.000.000.0001-91", 0}, // Wrong separator (dot instead of slash)
		{"Noise", "Use REF/0000-000", 0},
		{"Unicode Noise", "Company 00.000.000/0001-91 🚀", 1}, // Ensure unicode doesn't break offset
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
// Ensures the detector can be used by multiple goroutines simultaneously
// without data races (common in API Gateways).
func TestCNPJDetector_Concurrency(t *testing.T) {
	d := NewCNPJDetector()
	payload := "Text with CNPJ 06.990.590/0001-23 inside."
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
// Runs the parser against random inputs to ensure no PANICS occur.
// Run with: go test -fuzz=FuzzCNPJ -fuzztime=10s
func FuzzCNPJDetector(f *testing.F) {
	d := NewCNPJDetector()

	// Seed corpus (valid and near-valid examples to guide the fuzzer)
	f.Add("00.000.000/0001-91")
	f.Add("Random text")
	f.Add("12345678901234")

	f.Fuzz(func(t *testing.T, orig string) {
		// The only rule here is: MUST NOT PANIC.
		// The result (found or not) doesn't matter as input is random.
		_ = d.Scan(orig)
	})
}

// Benchmark
func BenchmarkCNPJDetector_LongText(b *testing.B) {
	d := NewCNPJDetector()
	payload := `
Fiscal data:
00.000.000/0001-91
06.990.590/0001-23
12.345.678/0001-00
33.649.575/0001-99
`
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.Scan(payload)
	}
}
