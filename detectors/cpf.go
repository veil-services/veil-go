package detectors

import (
	"regexp"
	"strings"
)

type CPFDetector struct{}

// Regex: \d{3}\.?\d{3}\.?\d{3}-?\d{2}
var cpfRegex = regexp.MustCompile(`\d{3}\.?\d{3}\.?\d{3}-?\d{2}`)

func (d *CPFDetector) Name() string {
	return "br_cpf"
}

func (d *CPFDetector) Scan(input string) []Match {
	matches := cpfRegex.FindAllStringIndex(input, -1)
	if len(matches) == 0 {
		return nil
	}

	// Pre-alocar slice para evitar resizes (micro-otimização)
	results := make([]Match, 0, len(matches))

	for _, loc := range matches {
		start, end := loc[0], loc[1]
		val := input[start:end]

		// Limpeza otimizada (evita múltiplas passadas se possível, mas strings.ReplaceAll é rápido o suficiente para v1)
		cleanVal := strings.ReplaceAll(val, ".", "")
		cleanVal = strings.ReplaceAll(cleanVal, "-", "")

		if len(cleanVal) != 11 {
			continue
		}

		if isValidCPF(cleanVal) {
			results = append(results, Match{
				StartIndex: start,
				EndIndex:   end,
				Value:      val,
				Type:       "CPF",
				Score:      1.0,
			})
		}
	}
	return results
}

func NewCPFDetector() Detector {
	return &CPFDetector{}
}

// isValidCPF implementa o algoritmo de Módulo 11 usando aritmética de bytes para zero alocação.
func isValidCPF(cpf string) bool {
	// Verifica se todos os dígitos são iguais
	// Otimização: Comparar byte a byte é muito rápido
	first := cpf[0]
	allEquals := true
	for i := 1; i < 11; i++ {
		if cpf[i] != first {
			allEquals = false
			break
		}
	}
	if allEquals {
		return false
	}

	// Conversão rápida de byte para int: '0' é 48 na tabela ASCII.
	// Então (char - '0') dá o valor numérico.

	// Primeiro Dígito
	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(cpf[i]-'0') * (10 - i)
	}
	remainder := sum % 11
	digit1 := 0
	if remainder >= 2 {
		digit1 = 11 - remainder
	}

	if int(cpf[9]-'0') != digit1 {
		return false
	}

	// Segundo Dígito
	sum = 0
	for i := 0; i < 10; i++ {
		sum += int(cpf[i]-'0') * (11 - i)
	}
	remainder = sum % 11
	digit2 := 0
	if remainder >= 2 {
		digit2 = 11 - remainder
	}

	return int(cpf[10]-'0') == digit2
}
