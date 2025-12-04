package detectors

import (
	"regexp"
)

type UUIDDetector struct{}

// Regex para UUID v4 (e variantes comuns de 32 hex chars com hifens)
// Formato: 8-4-4-4-12
// Ex: 123e4567-e89b-12d3-a456-426614174000
var uuidRegex = regexp.MustCompile(`(?i)\b[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}\b`)

func (d *UUIDDetector) Name() string {
	return "global_uuid"
}

func (d *UUIDDetector) Scan(input string) []Match {
	matches := uuidRegex.FindAllStringIndex(input, -1)
	if len(matches) == 0 {
		return nil
	}

	results := make([]Match, 0, len(matches))

	for _, loc := range matches {
		start, end := loc[0], loc[1]
		val := input[start:end]

		results = append(results, Match{
			StartIndex: start,
			EndIndex:   end,
			Value:      val,
			Type:       "UUID",
			Score:      1.0,
		})
	}
	return results
}

func NewUUIDDetector() Detector {
	return &UUIDDetector{}
}

