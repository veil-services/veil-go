package detectors

import "testing"

func TestUUIDDetector(t *testing.T) {
	d := NewUUIDDetector()

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"UUID v4", "Trace 123e4567-e89b-12d3-a456-426614174000", 1},
		{"Uppercase", "Trace 123E4567-E89B-12D3-A456-426614174000", 1},
		{"Surrounded By Text", "before 123e4567-e89b-12d3-a456-426614174000 after", 1},
		{"Multiple UUIDs", "IDs 123e4567-e89b-12d3-a456-426614174000 and a0e9825f-561a-4835-b2a8-2921c2618991", 2},
		{"Missing Section", "123e4567-e89b-12d3-a456-42661417400", 0},
		{"No Hyphens", "123e4567e89b12d3a456426614174000", 0},
		{"With Neighboring Hex", "foo123e4567-e89b-12d3-a456-426614174000bar", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(d.Scan(tt.input)); got != tt.expected {
				t.Errorf("expected %d matches, got %d", tt.expected, got)
			}
		})
	}
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

