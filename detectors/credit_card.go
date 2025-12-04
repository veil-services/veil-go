package detectors

type CreditCardDetector struct{}

func (d *CreditCardDetector) Name() string {
	return "global_credit_card"
}

func (d *CreditCardDetector) Scan(input string) []Match {
	var results []Match
	var digitsBuf [19]byte

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

		for j < len(input) && count < 19 {
			c := input[j]
			switch {
			case isDigitChar(c):
				digitsBuf[count] = c
				count++
			case (c == ' ' || c == '-') && count > 0:
				// allow separator
			default:
				goto evaluate
			}
			j++
		}

	evaluate:
		if count >= 13 && count <= 19 {
			if j < len(input) && isDigitChar(input[j]) {
				// part of longer sequence, skip
				i = j
				continue
			}
			if isValidLuhnBytes(digitsBuf[:count]) {
				results = append(results, Match{
					StartIndex: start,
					EndIndex:   j,
					Value:      input[start:j],
					Type:       TypeCreditCard,
					Score:      1.0,
				})
				i = j - 1
				continue
			}
		}
	}

	return results
}

func NewCreditCardDetector() Detector {
	return &CreditCardDetector{}
}

// isValidLuhnBytes implements the Luhn algorithm without allocations.
func isValidLuhnBytes(number []byte) bool {
	sum := 0
	alternate := false

	for i := len(number) - 1; i >= 0; i-- {
		n := int(number[i] - '0')
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
