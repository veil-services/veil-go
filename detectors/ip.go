package detectors

type IPDetector struct{}

func (d *IPDetector) Name() string {
	return "global_ipv4"
}

func (d *IPDetector) Scan(input string) []Match {
	var results []Match
	for i := 0; i < len(input); i++ {
		if !isDigitChar(input[i]) {
			continue
		}
		if start, end, ok := matchIPv4(input, i); ok {
			results = append(results, Match{
				StartIndex: start,
				EndIndex:   end,
				Value:      input[start:end],
				Type:       TypeIP,
				Score:      1.0,
			})
			i = end - 1
		}
	}
	return results
}

func NewIPDetector() Detector {
	return &IPDetector{}
}

func matchIPv4(s string, start int) (int, int, bool) {
	if start > 0 {
		prev := s[start-1]
		if isDigitChar(prev) || prev == '.' || prev == '-' {
			return 0, 0, false
		}
	}

	n := len(s)
	idx := start

	for octet := 0; octet < 4; octet++ {
		if idx >= n || !isDigitChar(s[idx]) {
			return 0, 0, false
		}
		val := 0
		digits := 0
		firstChar := s[idx]
		for idx < n && isDigitChar(s[idx]) {
			val = val*10 + int(s[idx]-'0')
			digits++
			if val > 255 || digits > 3 {
				return 0, 0, false
			}
			idx++
		}
		if digits == 0 {
			return 0, 0, false
		}
		if digits > 1 && firstChar == '0' {
			return 0, 0, false
		}
		if octet < 3 {
			if idx >= n || s[idx] != '.' {
				return 0, 0, false
			}
			idx++
		}
	}

	if idx < n {
		if isDigitChar(s[idx]) || s[idx] == '.' || s[idx] == '-' || isLetter(s[idx]) {
			return 0, 0, false
		}
	}

	return start, idx, true
}
