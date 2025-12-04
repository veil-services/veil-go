package detectors

import (
	"regexp"
	"strings"
)

type CNPJDetector struct{}

// Regex: xx.xxx.xxx/0001-xx ou apenas números (14 dígitos)
var cnpjRegex = regexp.MustCompile(`\d{2}\.?\d{3}\.?\d{3}/?\d{4}-?\d{2}`)

func (d *CNPJDetector) Name() string {
	return "br_cnpj"
}

func (d *CNPJDetector) Scan(input string) []Match {
	matches := cnpjRegex.FindAllStringIndex(input, -1)
	if len(matches) == 0 {
		return nil
	}

	results := make([]Match, 0, len(matches))

	for _, loc := range matches {
		start, end := loc[0], loc[1]
		val := input[start:end]

		cleanVal := strings.ReplaceAll(val, ".", "")
		cleanVal = strings.ReplaceAll(cleanVal, "/", "")
		cleanVal = strings.ReplaceAll(cleanVal, "-", "")

		if len(cleanVal) != 14 {
			continue
		}

		if isValidCNPJ(cleanVal) {
			results = append(results, Match{
				StartIndex: start,
				EndIndex:   end,
				Value:      val,
				Type:       "CNPJ",
				Score:      1.0,
			})
		}
	}
	return results
}

func NewCNPJDetector() Detector {
	return &CNPJDetector{}
}

// isValidCNPJ valida o CNPJ usando Módulo 11 com pesos específicos e aritmética de bytes.
func isValidCNPJ(cnpj string) bool {
	// Verifica dígitos iguais
	first := cnpj[0]
	allEquals := true
	for i := 1; i < 14; i++ {
		if cnpj[i] != first {
			allEquals = false
			break
		}
	}
	if allEquals {
		return false
	}

	// Pesos
	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	// Primeiro Dígito
	sum := 0
	for i := 0; i < 12; i++ {
		sum += int(cnpj[i]-'0') * weights1[i]
	}
	remainder := sum % 11
	digit1 := 0
	if remainder >= 2 {
		digit1 = 11 - remainder
	}

	if int(cnpj[12]-'0') != digit1 {
		return false
	}

	// Segundo Dígito
	sum = 0
	for i := 0; i < 13; i++ {
		sum += int(cnpj[i]-'0') * weights2[i]
	}
	remainder = sum % 11
	digit2 := 0
	if remainder >= 2 {
		digit2 = 11 - remainder
	}

	return int(cnpj[13]-'0') == digit2
}
