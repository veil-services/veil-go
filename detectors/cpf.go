package detectors

import "regexp"

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

	// Pre-allocate slice to avoid resizes (micro-optimization)
	results := make([]Match, 0, len(matches))

	for _, loc := range matches {
		start, end := loc[0], loc[1]
		val := input[start:end]

		var digits [11]byte
		pos := 0
		for k := start; k < end && pos < 11; k++ {
			if isDigitChar(input[k]) {
				digits[pos] = input[k]
				pos++
			}
		}

		if pos != 11 {
			continue
		}

		if isValidCPFBytes(digits[:]) {
			results = append(results, Match{
				StartIndex: start,
				EndIndex:   end,
				Value:      val,
				Type:       TypeCPF,
				Score:      1.0,
			})
		}
	}
	return results
}

func NewCPFDetector() Detector {
	return &CPFDetector{}
}

// isValidCPF implements the Mod11 algorithm using byte arithmetic for zero-allocation.
func isValidCPFBytes(cpf []byte) bool {
	// Check if all digits are equal (e.g., 111.111.111-11)
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

	// Fast conversion from byte to int: '0' is 48 in ASCII.
	// So (char - '0') gives the numeric value.

	// First Digit
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

	// Second Digit
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
