package detectors

type CPFDetector struct{}

func (d *CPFDetector) Name() string {
	return "br_cpf"
}

func (d *CPFDetector) Scan(input string) []Match {
	var results []Match
	var digits [11]byte

	for i := 0; i < len(input); i++ {
		if !isDigitChar(input[i]) {
			continue
		}
		if i > 0 && isDigitChar(input[i-1]) {
			continue
		}

		start := i
		count := 0
		j := i

		for j < len(input) && count < 11 {
			c := input[j]
			switch {
			case isDigitChar(c):
				digits[count] = c
				count++
			case isCPFSeparator(c):
				if count == 0 || !validCPFSeparatorPosition(c, count) {
					goto nextCandidate
				}
			default:
				goto evaluate
			}
			j++
		}

	evaluate:
		if count == 11 {
			if j < len(input) && isDigitChar(input[j]) {
				goto nextCandidate
			}
			if isValidCPFBytes(digits[:]) {
				results = append(results, Match{
					StartIndex: start,
					EndIndex:   j,
					Value:      input[start:j],
					Type:       TypeCPF,
					Score:      1.0,
				})
				i = j - 1
				continue
			}
		}

	nextCandidate:
	}

	return results
}

func NewCPFDetector() Detector {
	return &CPFDetector{}
}

func validCPFSeparatorPosition(sep byte, count int) bool {
	switch sep {
	case '.':
		return count == 3 || count == 6
	case '-':
		return count == 9
	default:
		return false
	}
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
