package detectors

type CNPJDetector struct{}

func (d *CNPJDetector) Name() string {
	return "br_cnpj"
}

func (d *CNPJDetector) Scan(input string) []Match {
	var results []Match
	var buffer [14]byte

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

		for j < len(input) && count < 14 {
			c := input[j]
			switch {
			case isDigitChar(c):
				buffer[count] = c
				count++
			case isCNPJSeparator(c):
				if !validCNPJSeparatorPosition(c, count) {
					goto nextCandidate
				}
			case c == ' ':
				if count == 0 {
					goto nextCandidate
				}
			default:
				goto evaluate
			}
			j++
		}

	evaluate:
		if count == 14 {
			if j < len(input) && isDigitChar(input[j]) {
				goto nextCandidate
			}
			if isValidCNPJBytes(buffer[:]) {
				results = append(results, Match{
					StartIndex: start,
					EndIndex:   j,
					Value:      input[start:j],
					Type:       TypeCNPJ,
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

func NewCNPJDetector() Detector {
	return &CNPJDetector{}
}

func validCNPJSeparatorPosition(sep byte, digitCount int) bool {
	switch sep {
	case '.':
		return digitCount == 2 || digitCount == 5
	case '/':
		return digitCount == 8
	case '-':
		return digitCount == 12
	default:
		return false
	}
}

// isValidCNPJ validates CNPJ using Mod11 with specific weights and byte arithmetic.
func isValidCNPJBytes(cnpj []byte) bool {
	// Check equal digits
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

	// Weights
	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	// First Digit
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

	// Second Digit
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
