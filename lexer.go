package lexer

import (
	"bufio"
	//	"bufio"
	"fmt"
	"io"
	"strings"
	//	"unicode"
	"unicode/utf8"
)

type (
	TLexer struct {
		state             StateFn       // the next lexing function to enter
		input             io.RuneReader // the input source
		inputCurrentStart int           // start position of this item
		buffer            []rune        // a slice of runes that contains the currently lexed item
		bufferPos         int           // the current position in the buffer
		Tokens            chan TToken   // channel of scanned Tokens
	}
)

// lex creates a new scanner for the input string.
func NewLexer(input io.Reader) (*TLexer, *Error) {
	lexer := &TLexer{
		input:  bufio.NewReader(input),
		buffer: make([]rune, 0, 10),
		Tokens: make(chan TToken),
	}

	go lexer.Run()
	return lexer, nil
}

// run runs the state machine for the Lexer.
func (self *TLexer) Run() {
	for state := lexWhitespace; state != nil; {
		state = state(self)
	}

	close(self.Tokens)
}

// next() returns the next rune in the input.
func (self *TLexer) Next() rune {
	if self.bufferPos < len(self.buffer) {
		res := self.buffer[self.bufferPos]
		self.bufferPos++
		return res
	}

	r, _, err := self.input.ReadRune()
	if err == io.EOF {
		r = EOF
	} else if err != nil {
		panic(err)
	}

	self.buffer = append(self.buffer, r)
	self.bufferPos++
	return r
}

// peek() returns but does not consume the next rune in the input.
func (self *TLexer) Peek() rune {
	if self.bufferPos < len(self.buffer) {
		return self.buffer[self.bufferPos]
	}

	r, _, err := self.input.ReadRune()
	if err == io.EOF {
		r = EOF
	} else if err != nil {
		panic(err)
	}

	self.buffer = append(self.buffer, r)
	return r
}

// peek() returns but does not consume the next few runes in the input.
func (self *TLexer) PeekNext(length int) string {
	lenDiff := self.bufferPos + length - len(self.buffer)
	if lenDiff > 0 {
		for i := 0; i < lenDiff; i++ {
			r, _, err := self.input.ReadRune()
			if err == io.EOF {
				r = EOF
			} else if err != nil {
				panic(err)
			}

			self.buffer = append(self.buffer, r)
		}
	}

	return string(self.buffer[self.bufferPos : self.bufferPos+length])
}

// backup steps back one rune
func (self *TLexer) Backup() {
	self.BackupWith(1)
}

// backup steps back many runes
func (self *TLexer) BackupWith(length int) {
	if self.bufferPos < length {
		panic(fmt.Errorf("lexer: trying to backup with %d when the buffer position is %d", length, self.bufferPos))
	}

	self.bufferPos -= length
}

// emit passes an Item back to the client.
func (self *TLexer) Emit(t int) {
	self.Tokens <- TToken{Type: t, Pos: self.inputCurrentStart, Val: string(self.buffer[:self.bufferPos])}
	self.Ignore()
}

// ignore skips over the pending input before this point.
func (self *TLexer) Ignore() {
	itemByteLen := 0
	for i := 0; i < self.bufferPos; i++ {
		itemByteLen += utf8.RuneLen(self.buffer[i])
	}

	self.inputCurrentStart += itemByteLen
	self.buffer = self.buffer[self.bufferPos:] //TODO: check for memory leaks, maybe copy remaining items into a new slice?
	self.bufferPos = 0
}

// accept consumes the next rune if it's from the valid set.
func (self *TLexer) Accept(valid string) int {
	r := self.Next()
	if strings.IndexRune(valid, r) >= 0 {
		return 1
	}
	self.Backup()
	return 0
}

// acceptWhile consumes runes while the specified condition is true
func (self *TLexer) AcceptWhile(fn ValidatorFn) int {
	r := self.Next()
	count := 0
	for fn(r) {
		r = self.Next()
		count++
	}
	self.Backup()
	return count
}

// acceptUntil consumes runes until the specified contidtion is met
func (self *TLexer) AcceptUntil(fn ValidatorFn) int {
	r := self.Next()
	count := 0
	for !fn(r) && r != EOF {
		r = self.Next()
		count++
	}
	self.Backup()
	return count
}

// errorf returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating self.nextItem.
func (self *TLexer) Errorf(format string, args ...interface{}) StateFn {
	self.Tokens <- TToken{Type: TokenError, Pos: self.inputCurrentStart, Val: fmt.Sprintf(format, args...)}
	return nil
}

// nextItem returns the next Item from the input.
func (self *TLexer) NextToken() TToken {
	return <-self.Tokens
}
