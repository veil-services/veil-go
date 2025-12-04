package detectors

import (
	"regexp"
	"strings"
)

type CreditCardDetector struct{}

// Regex genérica para capturar sequências que parecem cartões (13 a 19 dígitos)
var ccRegex = regexp.MustCompile(`(?:\d{4}[- ]?){3}\d{1,7}|\d{13,19}`)

func (d *CreditCardDetector) Name() string {
	return "global_credit_card"
}

func (d *CreditCardDetector) Scan(input string) []Match {
	matches := ccRegex.FindAllStringIndex(input, -1)
	if len(matches) == 0 {
		return nil
	}

	results := make([]Match, 0, len(matches))

	for _, loc := range matches {
		start, end := loc[0], loc[1]
		val := input[start:end]

		// Limpar separadores
		cleanVal := strings.ReplaceAll(val, "-", "")
		cleanVal = strings.ReplaceAll(cleanVal, " ", "")

		length := len(cleanVal)
		if length < 13 || length > 19 {
			continue
		}

		if isValidLuhn(cleanVal) {
			results = append(results, Match{
				StartIndex: start,
				EndIndex:   end,
				Value:      val,
				Type:       "CREDIT_CARD",
				Score:      1.0,
			})
		}
	}
	return results
}

func NewCreditCardDetector() Detector {
	return &CreditCardDetector{}
}

// isValidLuhn implementa o algoritmo de Luhn usando aritmética de bytes.
func isValidLuhn(number string) bool {
	sum := 0
	alternate := false

	// Itera da direita para a esquerda
	for i := len(number) - 1; i >= 0; i-- {
		n := int(number[i] - '0')
		
		// Se não for dígito, deveria falhar (mas a regex já garante dígitos na maioria dos casos)
		// Em caso de sujeira não filtrada:
		if n < 0 || n > 9 {
			return false
		}

		if alternate {
			n *= 2
			if n > 9 {
				n = (n % 10) + 1
			}
		}
		sum += n
		alternate = !alternate
	}

	return sum%10 == 0
}
