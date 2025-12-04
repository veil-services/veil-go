package detectors

import (
	"regexp"
	"strings"
)

type PhoneDetector struct{}

// E.164 Simplificado para detecção global
// Formato: Começa com +, seguido de 7 a 15 dígitos.
// Pode conter espaços ou traços como separadores visuais.
// Ex: +1 555 123 4567, +55-11-99999-9999
var e164Regex = regexp.MustCompile(`\+(?:[0-9][ -]?){6,14}[0-9]`)

func (d *PhoneDetector) Name() string {
	return "global_phone_e164"
}

func (d *PhoneDetector) Scan(input string) []Match {
	matches := e164Regex.FindAllStringIndex(input, -1)
	if len(matches) == 0 {
		return nil
	}

	results := make([]Match, 0, len(matches))

	for _, loc := range matches {
		start, end := loc[0], loc[1]
		val := input[start:end]

		// Limpeza para contagem de dígitos
		digitsOnly := strings.ReplaceAll(val, " ", "")
		digitsOnly = strings.ReplaceAll(digitsOnly, "-", "")
		digitsOnly = strings.ReplaceAll(digitsOnly, "+", "") // remove o prefixo para contar

		// E.164 define max 15 dígitos (sem contar o +)
		// Minimo razoável para um número internacional funcional é ~7-8 (ex: pequenos países)
		length := len(digitsOnly)
		if length < 7 || length > 15 {
			continue
		}

		results = append(results, Match{
			StartIndex: start,
			EndIndex:   end,
			Value:      val,
			Type:       "PHONE",
			Score:      1.0,
		})
	}
	return results
}

func NewPhoneDetector() Detector {
	return &PhoneDetector{}
}

