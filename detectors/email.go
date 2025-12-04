package detectors

import (
	"regexp"
)

type EmailDetector struct{}

// Regex simplificada conforme RFC 5322 (permissiva)
// Ref: Spec 002
var emailRegex = regexp.MustCompile(`(?i)[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}`)

func (d *EmailDetector) Name() string {
	return "email"
}

func (d *EmailDetector) Scan(input string) []Match {
	matches := emailRegex.FindAllStringIndex(input, -1)
	var results []Match

	for _, loc := range matches {
		start, end := loc[0], loc[1]
		val := input[start:end]
		results = append(results, Match{
			StartIndex: start,
			EndIndex:   end,
			Value:      val,
			Type:       "EMAIL",
			Score:      1.0, // Regex de email é estruturalmente confiável
		})
	}
	return results
}

func NewEmailDetector() Detector {
	return &EmailDetector{}
}

