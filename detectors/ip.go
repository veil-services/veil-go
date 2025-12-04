package detectors

import (
	"net"
	"regexp"
)

type IPDetector struct{}

// Regex IPv4 simples (grupos de 1-3 dígitos)
// A validação real de range (0-255) será feita no código para precisão.
var ipRegex = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)

func (d *IPDetector) Name() string {
	return "global_ipv4"
}

func (d *IPDetector) Scan(input string) []Match {
	matches := ipRegex.FindAllStringIndex(input, -1)
	if len(matches) == 0 {
		return nil
	}

	results := make([]Match, 0, len(matches))

	for _, loc := range matches {
		start, end := loc[0], loc[1]
		val := input[start:end]

		// Validação Real: net.ParseIP do Go é a melhor fonte da verdade
		// Ele valida ranges (0-255) e formato corretamente.
		if net.ParseIP(val) != nil {
			results = append(results, Match{
				StartIndex: start,
				EndIndex:   end,
				Value:      val,
				Type:       "IP",
				Score:      1.0,
			})
		}
	}
	return results
}

func NewIPDetector() Detector {
	return &IPDetector{}
}
