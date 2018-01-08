package lexer

import (
	"fmt"
)

type (
	// Token represents a token or text string returned from the scanner.
	TToken struct {
		Type int    // The type of this Token.
		Val  string // The value of this Token.

		Pos  int // The starting position, in bytes, of this item in the input string.
		Line int
		Col  int
	}
)

// A set of constants for precedence-based expression parsing.
// Non-operators have lowest precedence, followed by operators
// starting with precedence 1 up to unary operators. The highest
// precedence serves as "catch-all" precedence for selector,
// indexing, and other operator and delimiter tokens.
//
const (
	LowestPrec  = 0 // non-operators
	UnaryPrec   = 6
	HighestPrec = 7
)

const (
	EOF = -1 // 必须等于-1 详细见isEndOfLine
	// base token for common
	TokenError = iota // error occurred; value is text of error
	//TokenEOF                      // end of the file
	TokenWhitespace        // a run of spaces, tabs and newlines
	TokenSingleLineComment // A comment like --
	TokenMultiLineComment  // A multiline comment like /* ... */
	TokenKeyword           // SQL language keyword like SELECT, INSERT, etc.
	//TokenIdentifier               // alphanumeric identifier or complex identifier like `a.b` and `c`.*
	TokenOperator // operators like '=', '<>', etc.
	//TokenLeftParen                // '('
	//TokenRightParen               // ')'
	//TokenComma                    // ','
	//TokenDot                      // '.'
	//TokenStetementEnd             // ';'
	TokenNumber      // simple number, including imaginary
	TokenString      // quoted string (includes quotes)
	TokenValueHolder // ?

	// The list of tokens.
	// Special tokens
	ILLEGAL
	//EOF
	COMMENT

	NUMBER   // simple number, including imaginary
	HOLDER   // ?
	OPERATOR //operators like '=', '<>', etc.
	KEYWORD  // SQL language keyword like SELECT, INSERT, etc.
	SAPCE    // a run of spaces, tabs and newlines

	literal_beg
	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	IDENT  // main // alphanumeric identifier or complex identifier like `a.b` and `c`.*
	INT    // 12345
	FLOAT  // 123.45
	IMAG   // 123.45i
	CHAR   // 'a'
	STRING // "abc"

	literal_end

	operator_beg
	// Operators and delimiters
	ADD // +
	SUB // -
	MUL // *
	QUO // /
	REM // %

	AND     // &
	OR      // |
	XOR     // ^
	SHL     // <<
	SHR     // >>
	AND_NOT // &^

	ADD_ASSIGN // +=
	SUB_ASSIGN // -=
	MUL_ASSIGN // *=
	QUO_ASSIGN // /=
	REM_ASSIGN // %=

	AND_ASSIGN     // &=
	OR_ASSIGN      // |=
	XOR_ASSIGN     // ^=
	SHL_ASSIGN     // <<=
	SHR_ASSIGN     // >>=
	AND_NOT_ASSIGN // &^=

	LAND  // &&
	LOR   // ||
	ARROW // <-
	INC   // ++
	DEC   // --

	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !

	NEQ      // !=
	LEQ      // <=
	GEQ      // >=
	DEFINE   // :=
	ELLIPSIS // ...

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :
	operator_end

	keyword_beg
	// Keywords
	BREAK
	CASE
	CHAN
	CONST
	CONTINUE

	DEFAULT
	DEFER
	ELSE
	FALLTHROUGH
	FOR

	FUNC
	GO
	GOTO
	IF
	IMPORT

	INTERFACE
	MAP
	PACKAGE
	RANGE
	RETURN

	SELECT
	STRUCT
	SWITCH
	TYPE
	VAR
	keyword_end
)

// Precedence returns the operator precedence of the binary
// operator op. If op is not a binary operator, the result
// is LowestPrecedence.
//
func Precedence(t *TToken) int {
	switch t.Val {
	case "||", "or": //LOR:
		return 1
	case "&&", "and": //LAND:
		return 2
	case "==", "!=", "<", "<=", ">", ">=", "in": //EQL, NEQ, LSS, LEQ, GTR, GEQ:
		return 3
	case "+", "-", "|", "^": //ADD, SUB, OR, XOR:
		return 4
	case "*", "/", "%", "<<", ">>", "&", "&^":
		//case MUL, QUO, REM, SHL, SHR, AND, AND_NOT:
		return 5
	}
	return LowestPrec
}

var (
	TokenNames = map[int]string{
		TokenError:      "error",
		EOF:             "EOF",
		ILLEGAL:         "ILLEGAL",
		COMMENT:         "COMMENT",
		KEYWORD:         "KEYWORD",
		OPERATOR:        "OPERATOR",
		COLON:           "COLON",
		COMMA:           "COMMA",
		LPAREN:          "left_paren",
		RPAREN:          "right_paren",
		LBRACK:          "LBRACK",
		RBRACK:          "RBRACK",
		LBRACE:          "LBRACE",
		RBRACE:          "RBRACE",
		SAPCE:           "SAPCE",
		TokenWhitespace: "space",
		//TokenStatementStart: "statement_start",
		//TokenStetementEnd: "statement_end",
		PERIOD: "PERIOD",
		HOLDER: "HOLDER",
		NUMBER: "NUMBER",
		IDENT:  "IDENT",  // main // alphanumeric identifier or complex identifier like `a.b` and `c`.*
		INT:    "INT",    // 12345
		FLOAT:  "FLOAT",  // 123.45
		IMAG:   "IMAG",   // 123.45i
		CHAR:   "CHAR",   // 'a'
		STRING: "STRING", // "abc"
	}
)

func PrintToken(t TToken) string {
	if TokenNames[t.Type] != "" {
		return fmt.Sprintf(">> %q('%q')", TokenNames[t.Type], t.Val)
	}

	return ""
}
