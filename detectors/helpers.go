package detectors

func isDigitChar(b byte) bool {
	return b >= '0' && b <= '9'
}

func isHexChar(b byte) bool {
	return (b >= '0' && b <= '9') ||
		(b >= 'a' && b <= 'f') ||
		(b >= 'A' && b <= 'F')
}
