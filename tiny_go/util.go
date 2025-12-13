package tinygo

func isDigitOrAlpha(b byte) bool {
	return isDigit(b) || isAlpha(b)
}
func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

func isAlpha(b byte) bool {
	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z')
}
