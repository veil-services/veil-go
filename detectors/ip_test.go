package detectors

import "testing"

func TestIPDetector(t *testing.T) {
	d := NewIPDetector()

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Localhost", "Server at 127.0.0.1", 1},
		{"Private Range", "Hit 10.0.0.42 today", 1},
		{"Public Range", "Client 8.8.8.8 requested DNS", 1},
		{"Invalid Segment", "Address 256.0.0.1 invalid", 0},
		{"Loopback Fragment", "Version 1.2.3", 0},
		{"Short Fragment", "Log shows 192.168", 0},
		{"Mixed Text", "IDs 10.0.0.1a and 10.0.0.2b", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(d.Scan(tt.input)); got != tt.expected {
				t.Errorf("expected %d matches, got %d", tt.expected, got)
			}
		})
	}
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

