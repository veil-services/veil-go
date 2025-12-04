package detectors

import "testing"

func TestCNPJDetector(t *testing.T) {
	d := NewCNPJDetector()

	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Valid BB", "Empresa 00.000.000/0001-91 cadastrada.", 1},
		{"Valid Google", "Filial 06.990.590/0001-23", 1},
		{"Valid Plain", "06990590000123", 1}, // same as Google CNPJ without punctuation
		{"Valid Embedded", "Texto XX00.000.000/0001-91YY", 1},
		{"Invalid Checksum", "00.000.000/0001-00", 0},
		{"Short Sequence", "12.345", 0},
		{"Noise", "Use REF/0000-000", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(d.Scan(tt.input)); got != tt.expected {
				t.Errorf("expected %d matches, got %d", tt.expected, got)
			}
		})
	}
}

func TestCNPJDetector_Report(t *testing.T) {
	d := NewCNPJDetector()
	text := `
Empresas:
- Matriz 00.000.000/0001-91
- Subsidiária 06.990.590/0001-23
- Fornecedor 12.345.678/0001-00 (inválido, ignorar)
`
	if got := len(d.Scan(text)); got != 2 {
		t.Fatalf("expected 2 valid CNPJs, got %d", got)
	}
}

func BenchmarkCNPJDetector_LongText(b *testing.B) {
	d := NewCNPJDetector()
	payload := `
Dados fiscais:
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

