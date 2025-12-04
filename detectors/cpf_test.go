package detectors

import "testing"

func TestCPFDetector(t *testing.T) {
	d := NewCPFDetector()

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Valid Formatted", "CPF do cliente: 111.444.777-35", 1},
		{"Valid Plain", "Documento: 11144477735", 1},
		{"Embedded Valid", "IDs: abc111.444.777-35xyz", 1},
		{"Second Valid", "Dados: 529.982.247-25", 1},
		{"Invalid Checksum", "CPF 123.456.789-88", 0},
		{"Repeated Digits", "111.111.111-11", 0},
		{"Short Sequence", "123.456", 0},
		{"Numeric Noise", "Order total 12345-67890", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(d.Scan(tt.input)); got != tt.expected {
				t.Errorf("expected %d matches, got %d", tt.expected, got)
			}
		})
	}
}

func TestCPFDetector_LongParagraph(t *testing.T) {
	d := NewCPFDetector()
	text := `
Cadastro:
 - Cliente 1: 111.444.777-35
 - Cliente 2: CPF 52998224725
 - Cliente 3: mascarado 000.000.000-00 (inválido, deve ser ignorado)
`
	if got := len(d.Scan(text)); got != 2 {
		t.Fatalf("expected 2 valid CPFs, got %d", got)
	}
}

func BenchmarkCPFDetector_LongText(b *testing.B) {
	d := NewCPFDetector()
	payload := `
Relatório:
Cliente A CPF 111.444.777-35
Cliente B CPF 529.982.247-25
Cliente C CPF 862.883.667-57
`
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.Scan(payload)
	}
}

