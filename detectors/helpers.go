package detectors

func isDigitChar(b byte) bool {
	return b >= '0' && b <= '9'
}

func isHexChar(b byte) bool {
	return (b >= '0' && b <= '9') ||
		(b >= 'a' && b <= 'f') ||
		(b >= 'A' && b <= 'F')
}

func isCNPJSeparator(b byte) bool {
	return b == '.' || b == '-' || b == '/' || b == ' '
}

func isCPFSeparator(b byte) bool {
	return b == '.' || b == '-' || b == ' '
}
