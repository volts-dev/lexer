package lexer

import (
	"fmt"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	lex, err := NewLexer(strings.NewReader("1.1+1=2"))
	if err != nil {

	}

	for {
		token, ok := <-lex.Tokens
		if !ok {
			break
		}

		fmt.Println(PrintToken(token))
	}
	fmt.Println("complete")
}
