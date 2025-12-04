package detectors

type UUIDDetector struct{}

func (d *UUIDDetector) Name() string {
	return "global_uuid"
}

func (d *UUIDDetector) Scan(input string) []Match {
	var results []Match
	for i := 0; i < len(input); i++ {
		if !isHexChar(input[i]) {
			continue
		}
		if end, ok := matchUUID(input, i); ok {
			results = append(results, Match{
				StartIndex: i,
				EndIndex:   end,
				Value:      input[i:end],
				Type:       TypeUUID,
				Score:      1.0,
			})
			i = end - 1
		}
	}
	return results
}

func NewUUIDDetector() Detector {
	return &UUIDDetector{}
}

func matchUUID(s string, start int) (int, bool) {
	sections := []int{8, 4, 4, 4, 12}
	idx := start
	n := len(s)

	// boundary: preceding char cannot be hex or '-'
	if start > 0 {
		prev := s[start-1]
		if isHexChar(prev) || prev == '-' {
			return 0, false
		}
	}

	for si, length := range sections {
		for j := 0; j < length; j++ {
			if idx >= n || !isHexChar(s[idx]) {
				return 0, false
			}
			idx++
		}
		if si < len(sections)-1 {
			if idx >= n || s[idx] != '-' {
				return 0, false
			}
			idx++
		}
	}

	return idx, true
}
