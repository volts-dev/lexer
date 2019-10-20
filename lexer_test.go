package lexer

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	osfile, e := os.Open("./lexer_test_file2.js")
	if e != nil {
		return
	}
	defer osfile.Close()

	lex, err := NewLexer(osfile)
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
