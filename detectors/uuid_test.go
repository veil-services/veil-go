package detectors

import (
	"sync"
	"testing"
)

func TestUUIDDetector(t *testing.T) {
	d := NewUUIDDetector()

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		// Valid Cases
		{"UUID v4 Lowercase", "Trace 123e4567-e89b-12d3-a456-426614174000", 1},
		{"UUID Uppercase", "ID 123E4567-E89B-12D3-A456-426614174000", 1},
		{"Nil UUID", "Zero 00000000-0000-0000-0000-000000000000", 1},
		{"Surrounded By Text", "before 123e4567-e89b-12d3-a456-426614174000 after", 1},
		{"Multiple UUIDs", "IDs 123e4567-e89b-12d3-a456-426614174000 and a0e9825f-561a-4835-b2a8-2921c2618991", 2},
		{"With Neighboring Hex", "foo123e4567-e89b-12d3-a456-426614174000bar", 0},

		// Invalid Cases (Strict Format)
		{"Missing Section", "123e4567-e89b-12d3-a456-42661417400", 0},        // Too short
		{"No Hyphens", "123e4567e89b12d3a456426614174000", 0},                // Strict UUID usually requires hyphens
		{"Invalid Hex Chars", "123e4567-e89b-12d3-a456-42661417400G", 0},     // 'G' is not hex
		{"Too Long", "123e4567-e89b-12d3-a456-4266141740000", 0},             // Extra digit at end
		{"Wrong Hyphen Position", "123e4567-e89b-12d3-a4564-26614174000", 0}, // Broken structure
		{"Garbage", "not-a-uuid-string", 0},

		// Noise & Boundary
		{"Unicode Noise", "UUID: 123e4567-e89b-12d3-a456-426614174000 🚀", 1},
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
func TestUUIDDetector_Concurrency(t *testing.T) {
	d := NewUUIDDetector()
	payload := "Request ID: 123e4567-e89b-12d3-a456-426614174000 processing."
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
// Run with: go test -fuzz=FuzzUUID -fuzztime=10s
func FuzzUUIDDetector(f *testing.F) {
	d := NewUUIDDetector()

	// Seed corpus
	f.Add("123e4567-e89b-12d3-a456-426614174000")
	f.Add("00000000-0000-0000-0000-000000000000")
	f.Add("random-text-with-hyphens")

	f.Fuzz(func(t *testing.T, orig string) {
		// MUST NOT PANIC
		_ = d.Scan(orig)
	})
}

func TestUUIDDetector_Log(t *testing.T) {
	d := NewUUIDDetector()
	log := `
Processing events:
- request 123e4567-e89b-12d3-a456-426614174000
- retry a0e9825f-561a-4835-b2a8-2921c2618991
- ignore malformed 123e4567e89b12d3a456426614174000
`
	if got := len(d.Scan(log)); got != 2 {
		t.Fatalf("expected 2 UUIDs, got %d", got)
	}
}

func BenchmarkUUIDDetector_LongText(b *testing.B) {
	d := NewUUIDDetector()
	payload := `
TraceIDs:
123e4567-e89b-12d3-a456-426614174000
aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee
deadbeef-dead-beef-dead-beefdeadbeef
`
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.Scan(payload)
	}
}
