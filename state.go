package lexer

import (
	"unicode"
)

type (
	// StateFn represents the state of the scanner as a function that returns the next state.
	StateFn func(*TLexer) StateFn

	// ValidatorFn represents a function that is used to check whether a specific rune matches certain rules.
	ValidatorFn func(rune) bool
)

var (
	states_fn map[string]StateFn
)

func init() {
	states_fn = map[string]StateFn{
		"?": lexHolder,
		".": lexPeriod, //new

		",":  lexComma,
		";":  lexStetementEnd,
		":":  lexColon,
		`"`:  lexString,
		`'`:  lexString,
		"(":  lexLParen,
		")":  lexRParen,
		"[":  lexLBrack,
		"]":  lexRBrack,
		"{":  lexLBrace,
		"}":  lexRBrace,
		"%s": lexHolder,
		"/*": lexMultiLineComment,
		"//": lexSingleLineComment,
		"--": lexSingleLineComment,

		//"\\": lexStringExcapeComment,
	}
}

func RegisterState(lit string, fn StateFn) {
	states_fn[lit] = fn
}

func state_filter2(lit string) StateFn {
	return states_fn[lit]
}

func state_filter1(lit string) StateFn {
	return states_fn[lit]
}

func lexPeriod(lexer *TLexer) StateFn {
	lexer.Next()
	lexer.Emit(PERIOD)
	return lexWhitespace
}

func lexComma(lexer *TLexer) StateFn {
	lexer.Next()
	lexer.Emit(COMMA)
	return lexWhitespace
}

func lexStetementEnd(lexer *TLexer) StateFn {
	lexer.Next()
	lexer.Emit(SEMICOLON)
	return lexWhitespace
}

func lexColon(lexer *TLexer) StateFn {
	lexer.Next()
	lexer.Emit(COLON)
	return lexWhitespace
}

func lexLParen(lexer *TLexer) StateFn {
	lexer.Next()
	lexer.Emit(LPAREN)
	return lexWhitespace
}

func lexRParen(lexer *TLexer) StateFn {
	lexer.Next()
	lexer.Emit(RPAREN)
	return lexWhitespace
}

func lexLBrack(lexer *TLexer) StateFn {
	lexer.Next()
	lexer.Emit(LBRACK)
	return lexWhitespace
}

func lexRBrack(lexer *TLexer) StateFn {
	lexer.Next()
	lexer.Emit(RBRACK)
	return lexWhitespace
}

func lexLBrace(lexer *TLexer) StateFn {
	lexer.Next()
	lexer.Emit(LBRACE)
	return lexWhitespace
}

func lexRBrace(lexer *TLexer) StateFn {
	lexer.Next()
	lexer.Emit(RBRACE)
	return lexWhitespace
}
func lexWhitespace(lexer *TLexer) StateFn {
	lexer.AcceptWhile(isWhitespace)
	if lexer.bufferPos > 0 {
		lexer.Emit(SAPCE)
	}

	next := lexer.Peek()
	nextTwo := lexer.PeekNext(2)

	// 以下优先两个字符的
	if next == -1 {
		lexer.Emit(EOF)
		return nil
	} else if sf := states_fn[nextTwo]; sf != nil {
		return sf
	} else if sf := states_fn[string(next)]; sf != nil {
		return sf
	}

	switch {
	case next == EOF:
		lexer.Emit(EOF)
		return nil

	case isOperator(next):
		return lexOperator

	case ('0' <= next && next <= '9'):
		return lexNumber

	case IsAlphaNumericRune(next) || next == '`':
		return lexIdentifierOrKeyword

	default:
		return lexUnknown
		//lexer.Errorf("don't know what to do with '%s'", nextTwo)
		//return nil
	}
}

func lexSingleLineComment(lexer *TLexer) StateFn {
	lexer.AcceptUntil(isEndOfLine)
	lexer.Emit(TokenSingleLineComment)
	return lexWhitespace
}

func lexMultiLineComment(lexer *TLexer) StateFn {
	lexer.Next()
	lexer.Next()
	for {
		lexer.AcceptUntil(func(r rune) bool { return r == '*' })
		if lexer.PeekNext(2) == "*/" {
			lexer.Next()
			lexer.Next()
			lexer.Emit(TokenMultiLineComment)
			return lexWhitespace
		}

		if lexer.Peek() == EOF {
			lexer.Errorf("reached EOF when looking for comment end")
			return nil
		}

		lexer.Next()
	}
}

func lexStringExcapeComment(lexer *TLexer) StateFn {
	return lexWhitespace
}

func lexUnknown(lexer *TLexer) StateFn {
	lexer.Next()
	lexer.Emit(UNKNOWN)
	return lexWhitespace
}
func lexOperator(lexer *TLexer) StateFn {
	lexer.AcceptWhile(isOperator)
	lexer.Emit(OPERATOR)
	return lexWhitespace
}

func lexNumber(lexer *TLexer) StateFn {
	count := 0
	count += lexer.AcceptWhile(unicode.IsDigit)
	if lexer.Accept(".") > 0 {
		count += 1 + lexer.AcceptWhile(unicode.IsDigit)
		lexer.Emit(FLOAT)
	} else if lexer.Accept("eE") > 0 {
		count += 1 + lexer.Accept("+-")
		count += lexer.AcceptWhile(unicode.IsDigit)
		lexer.Emit(IMAG)
	} else if IsAlphaNumericRune(lexer.Peek()) {
		// We were lexing an identifier all along - backup and pass the ball
		lexer.BackupWith(count)
		return lexIdentifierOrKeyword
	} else {
		lexer.Emit(NUMBER)
	}

	return lexWhitespace
}

func lexString(lexer *TLexer) StateFn {
	quote := lexer.Next()
	lexer.Emit(QUOTES)

	lexer.Ignore()
	for {
		n := lexer.Next()

		// TODO 使用更好的回溯方法
		// 未构成STRING
		if n == EOF {
			lexer.Backup()
			lexer.Emit(UNKNOWN)
			return lexWhitespace
			//return lexer.Errorf("unterminated quoted string")
		}

		if n == '\\' {
			//TODO: fix possible problems with NO_BACKSLASH_ESCAPES mode
			if lexer.Peek() == EOF {
				return lexer.Errorf("unterminated quoted string")
			}
			lexer.Next()
		}

		// 构成STRING
		if n == quote {
			if lexer.Peek() == quote {
				lexer.Next()
			} else {

				lexer.Backup() //回退quote
				lexer.Emit(STRING)
				lexer.Next() //
				lexer.Emit(QUOTES)

				lexer.Ignore()
				return lexWhitespace
			}
		}
	}

}

func lexIdentifierOrKeyword(lexer *TLexer) StateFn {
	for {
		s := lexer.Next()
		if s == '`' {
			for {
				n := lexer.Next()
				if n == EOF {
					return lexer.Errorf("unterminated quoted string")
				} else if n == '`' {
					if lexer.Peek() == '`' {
						lexer.Next()
					} else {
						break
					}
				}
			}
			lexer.Emit(IDENT)
		} else if IsAlphaNumericRune(s) {
			lexer.AcceptWhile(IsAlphaNumericRune)

			//TODO: check whether token is a keyword or an identifier
			lexer.Emit(IDENT)
		}

		lexer.AcceptWhile(isWhitespace)
		if lexer.bufferPos > 0 {
			lexer.Emit(SAPCE)
		}

		if lexer.Peek() != '.' {
			break
		}

		lexer.Next()
		lexer.Emit(PERIOD)
	}

	return lexWhitespace
}

func lexHolder(lexer *TLexer) StateFn {
	lexer.AcceptWhile(isHolder)
	lexer.Emit(HOLDER)
	return lexWhitespace
}
