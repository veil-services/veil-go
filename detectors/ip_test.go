package detectors

import (
	"sync"
	"testing"
)

func TestIPDetector(t *testing.T) {
	d := NewIPDetector()

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		// Valid Cases
		{"Localhost", "Server at 127.0.0.1", 1},
		{"Private Range", "Hit 10.0.0.42 today", 1},
		{"Public Range", "Client 8.8.8.8 requested DNS", 1},
		{"High Range", "Router 192.168.255.254", 1},

		// Invalid Cases (Logic/Range)
		{"Invalid Segment > 255", "Address 256.0.0.1 invalid", 0},
		{"Negative Segment", "Address -1.0.0.1 invalid", 0}, // Parser likely ignores dash
		{"Leading Zeros", "Address 010.001.001.001", 0},     // Typically invalid or octal, usually ignored in strict parsing
		{"Too Many Segments", "1.2.3.4.5", 0},               // Or maybe 1 if it catches the first 4? Usually 0 for strict boundary.

		// Noise & Formats (Software versions, etc)
		{"Loopback Fragment", "Version 1.2.3 released", 0},
		{"Short Fragment", "Log shows 192.168", 0},
		{"Mixed Text", "IDs 10.0.0.1a and 10.0.0.2b", 0}, // Invalid boundary (attached letter)
		{"Date Confusion", "Date 2023.01.01", 0},         // Looks like IP but isn't

		// Boundaries
		{"Boundary Start", "192.168.1.1 is the IP", 1},
		{"Boundary End", "The IP is 10.0.0.1", 1},
		{"Unicode Noise", "IP: 8.8.4.4 🚀", 1},
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
func TestIPDetector_Concurrency(t *testing.T) {
	d := NewIPDetector()
	payload := "Thread safe test for 192.168.1.1 running."
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
// Run with: go test -fuzz=FuzzIP -fuzztime=10s
func FuzzIPDetector(f *testing.F) {
	d := NewIPDetector()

	// Seed corpus
	f.Add("192.168.0.1")
	f.Add("10.0.0.1")
	f.Add("255.255.255.255")
	f.Add("Version 1.0.0")

	f.Fuzz(func(t *testing.T, orig string) {
		// MUST NOT PANIC
		_ = d.Scan(orig)
	})
}

func TestIPDetector_LogSnippet(t *testing.T) {
	d := NewIPDetector()
	log := `
		[INFO] client 172.16.5.10 connected
		[INFO] forwarding to upstream 203.0.113.5
		[DEBUG] health-check 10.1.1.1 -> 10.1.1.254
		[WARN] rejected spoof 999.10.20.30
	`
	matches := d.Scan(log)
	// Expecting: 172.16.5.10, 203.0.113.5, 10.1.1.1, 10.1.1.254
	// 999.10.20.30 should be ignored (invalid segment)
	if len(matches) != 4 {
		t.Fatalf("expected 4 IPs, got %d", len(matches))
	}
}

func BenchmarkIPDetector_Log(b *testing.B) {
	d := NewIPDetector()
	payload := `
		2024-01-01 00:00:01 ACCEPT src=10.0.0.5 dst=192.168.0.10
		2024-01-01 00:00:02 ACCEPT src=10.0.0.6 dst=192.168.0.11
		2024-01-01 00:00:03 DROP src=203.0.113.50 dst=10.0.0.8
		2024-01-01 00:00:04 DROP src=198.51.100.2 dst=10.0.0.9
		2024-01-01 00:00:05 ACCEPT src=172.16.10.20 dst=10.0.0.10
	`
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.Scan(payload)
	}
}
