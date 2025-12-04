package detectors

type EmailDetector struct{}

func (d *EmailDetector) Name() string {
	return "email"
}

func (d *EmailDetector) Scan(input string) []Match {
	var results []Match
	for i := 0; i < len(input); i++ {
		if input[i] != '@' {
			continue
		}
		if start, end, ok := extractEmail(input, i); ok {
			results = append(results, Match{
				StartIndex: start,
				EndIndex:   end,
				Value:      input[start:end],
				Type:       TypeEmail,
				Score:      1.0,
			})
			i = end - 1
		}
	}
	return results
}

func NewEmailDetector() Detector {
	return &EmailDetector{}
}

func extractEmail(input string, at int) (int, int, bool) {
	if at == 0 || at == len(input)-1 {
		return 0, 0, false
	}

	// Expand left for local part
	start := at - 1
	for start >= 0 && isEmailLocalChar(input[start]) {
		start--
	}
	start++

	// Local part must exist and cannot start/end with '.'
	if start >= at || input[start] == '.' || input[at-1] == '.' {
		return 0, 0, false
	}

	for j := start + 1; j < at; j++ {
		if input[j] == '.' && input[j-1] == '.' {
			return 0, 0, false
		}
	}

	// Expand right for domain part
	end := at + 1
	for end < len(input) && isEmailDomainChar(input[end]) {
		end++
	}

	trimEnd := end
	for trimEnd > at+1 && isEmailTrailingPunct(input[trimEnd-1]) {
		trimEnd--
	}

	if trimEnd <= at+1 {
		return 0, 0, false
	}

	// TLD must be at least 2 alphabetic chars
	if tldLen := countTLD(input[at+1 : trimEnd]); tldLen < 2 {
		return 0, 0, false
	}

	// Domain cannot end with '-' or '.'
	if input[trimEnd-1] == '-' || input[trimEnd-1] == '.' {
		return 0, 0, false
	}

	if !hasDomainDot(input[at+1 : trimEnd]) {
		return 0, 0, false
	}

	// Boundary check: neighbor characters cannot be part of email syntax (letters/digits)
	if start > 0 && isEmailUnsafeBoundary(input[start-1]) {
		return 0, 0, false
	}
	if trimEnd < len(input) && isEmailUnsafeBoundary(input[trimEnd]) {
		return 0, 0, false
	}

	return start, trimEnd, true
}

func countTLD(domain string) int {
	length := 0
	for i := len(domain) - 1; i >= 0; i-- {
		c := domain[i]
		if c == '.' {
			return length
		}
		if !isLetter(c) {
			return 0
		}
		length++
	}
	return 0
}

func isEmailLocalChar(b byte) bool {
	switch {
	case b >= 'a' && b <= 'z':
		return true
	case b >= 'A' && b <= 'Z':
		return true
	case b >= '0' && b <= '9':
		return true
	case b == '.', b == '_', b == '%', b == '+', b == '-':
		return true
	default:
		return false
	}
}

func isEmailDomainChar(b byte) bool {
	switch {
	case b >= 'a' && b <= 'z':
		return true
	case b >= 'A' && b <= 'Z':
		return true
	case b >= '0' && b <= '9':
		return true
	case b == '.', b == '-':
		return true
	default:
		return false
	}
}

func isLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func isEmailTrailingPunct(b byte) bool {
	switch b {
	case '.', ',', ';', ':', '!', '?':
		return true
	default:
		return false
	}
}

func hasDomainDot(domain string) bool {
	for i := 0; i < len(domain); i++ {
		if domain[i] == '.' {
			return true
		}
	}
	return false
}

func isEmailUnsafeBoundary(b byte) bool {
	switch b {
	case '.', ',', ';', ':', '!', '?', '(', ')', '[', ']', '{', '}', '"', '\'', ' ':
		return false
	default:
		return isEmailLocalChar(b) || b == '@'
	}
}
