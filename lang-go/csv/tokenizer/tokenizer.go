package tokenizer

import (
	"errors"

	"github.com/nyehamene/lang/csv/token"
)

type Tokenizer struct {
	current int
	line    int
	source  string
}

type Error int8

const (
	EOF Error = iota
	ErrInvalidNewline
	ErrRangeOutOfBound
	ErrUnterminatedString
	ErrUnexpectedDQuote
)

type Reader interface {
	Read(b []byte) (int, error)
}

var NO_TOKEN token.Token

var errMsg = map[Error]string{
	EOF:                   "<EOF>",
	ErrInvalidNewline:     "<Invalid newline>",
	ErrRangeOutOfBound:    "<Range out of bound>",
	ErrUnterminatedString: "<Unterminated string>",
	ErrUnexpectedDQuote:   "<Unexpected double quote>",
}

func New(s string) *Tokenizer {
	return &Tokenizer{
		source: s,
	}
}

func (t *Tokenizer) TokenizeAll() ([]token.Token, error) {
	var fields []token.Token

	for {
		field, err := t.Tokenize()

		if errors.Is(err, EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		fields = append(fields, field)
	}

	return fields, nil
}

func (t *Tokenizer) Tokenize() (token.Token, error) {
	if t.IsAtEnd() {
		return NO_TOKEN, EOF
	}

	b := t.advance()

	switch b {
	case '"':
		return t.escapedValue()
	case '\r':
		if !t.peekEither('\n') {
			return token.Token{}, ErrInvalidNewline
		}
		_ = t.advance()
		t.line += 1
		return token.Token{Value: "\n", Token: token.Newline}, nil
	case '\n':
		t.line += 1
		return token.Token{Value: "\n", Token: token.Newline}, nil
	case ',':
		return token.Token{Value: ",", Token: token.Comma}, nil
	default:
		return t.value(string(b))
	}
}

func (t *Tokenizer) Mark() int {
	return t.current
}

func (t *Tokenizer) Reset(mark int) {
	t.current = mark
}

func (t *Tokenizer) IsAtEnd() bool {
	return t.current >= len(t.source)
}

func (t *Tokenizer) peekEither(bs ...byte) bool {
	if t.IsAtEnd() {
		return false
	}
	c := t.source[t.current]
	for _, n := range bs {
		if c == n {
			return true
		}
	}
	return false
}

func (t *Tokenizer) advance() byte {
	b := t.source[t.current]
	t.current += 1
	return b
}

func (t *Tokenizer) value(start string) (token.Token, error) {
	value := start

	for !t.IsAtEnd() && !t.peekEither(',', '"', '\n') {
		b := t.advance()
		value += string(b)
	}

	return token.Token{Value: value, Token: token.Value}, nil
}

func (t *Tokenizer) escapedValue() (token.Token, error) {
	var value string

	if t.IsAtEnd() {
		return token.Token{}, ErrUnterminatedString
	}

	isClosed := false
loop:
	for !t.IsAtEnd() {
		b := t.advance()

		switch b {
		case '"':
			if t.peekEither('"') {
				b := t.advance()
				value += string(b)
				continue
			}
			isClosed = true
			break loop
		default:
			value += string(b)
		}
	}

	if t.IsAtEnd() && !isClosed {
		return token.Token{}, ErrUnterminatedString
	}

	return token.Token{Value: value, Token: token.EscapedValue}, nil
}

func (te Error) Error() string {
	msg, ok := errMsg[te]
	if !ok {
		return "unexpected tokenizer error"
	}
	return msg
}
