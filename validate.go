package lexer

import (
	"unicode"
)

// isSpace reports whether r is a whitespace character (space or end of line).
func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n' || r == EOF
}

// TODO : ! ~
// isOperator reports whether r is an operator.
func isOperator(r rune) bool {
	return r == '+' || r == '-' || r == '*' || r == '/' || r == '=' || r == '>' || r == '<' || r == '~' || r == '|' || r == '^' || r == '&' || r == '%' || r == '!'
}

func isHolder(r rune) bool {
	return r == '?' || r == '%' || r == 's'
}

func IsAlphaNumericRune(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
