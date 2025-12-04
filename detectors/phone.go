package detectors

type PhoneDetector struct{}

func (d *PhoneDetector) Name() string {
	return "global_phone_e164"
}

func (d *PhoneDetector) Scan(input string) []Match {
	var results []Match

	for i := 0; i < len(input); i++ {
		if input[i] != '+' {
			continue
		}

		start := i
		digits := 0
		hasInvalidSeparator := false
		j := i + 1

		for j < len(input) {
			c := input[j]

			switch {
			case c >= '0' && c <= '9':
				digits++
				if digits > 15 {
					j++
					goto nextCandidate
				}
			case isPhoneSeparator(c):
				// separators are allowed only after the first digit
				if digits == 0 {
					hasInvalidSeparator = true
				}
			default:
				goto boundaryCheck
			}
			j++
		}

	boundaryCheck:
		// ensure the match is bounded (next char can't be digit)
		if j < len(input) && input[j] >= '0' && input[j] <= '9' {
			continue
		}

		if digits >= 7 && digits <= 15 && !hasInvalidSeparator {
			results = append(results, Match{
				StartIndex: start,
				EndIndex:   j,
				Value:      input[start:j],
				Type:       TypePhone,
				Score:      1.0,
			})
		}
	nextCandidate:
		i = j - 1
	}

	return results
}

func NewPhoneDetector() Detector {
	return &PhoneDetector{}
}

func isPhoneSeparator(b byte) bool {
	return b == ' ' || b == '-' || b == '.'
}
