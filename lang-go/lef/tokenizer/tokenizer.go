package tokenizer

import (
	"github.com/nyehamene/lang/lef/token"
)

type Tokenizer struct {
	line    int
	current int
	source  string
}

type Reader interface {
	Read(buf []byte) (int, error)
}

type Error int

const (
	EOF Error = iota
	ErrUnterminatedValue
	ErrUnterminatedRegex
	ErrEmptyValue
	ErrEmptyRegex
	ErrInvalidEscape
)

var tokenErrMsg = map[Error]string{
	EOF:                  "<EOF>",
	ErrUnterminatedValue: "<Unterminated value>",
	ErrUnterminatedRegex: "<Unterminated regex>",
	ErrEmptyRegex:        "<Empty regex>",
	ErrEmptyValue:        "<Empty value>",
	ErrInvalidEscape:     "<Invalid Escape>",
}

var NO_TOKEN token.Token

func New(s string) *Tokenizer {
	return &Tokenizer{
		source: s,
	}
}

func (t *Tokenizer) TokenizeAll() ([]token.Token, error) {
	var tokens []token.Token
	for !t.IsAtEnd() {
		token, err := t.Tokenize()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func (t *Tokenizer) Tokenize() (_ token.Token, err error) {
	start := t.current

	defer func() {
		if err != nil {
			t.current = start
		}
	}()

	if t.IsAtEnd() {
		return NO_TOKEN, EOF
	}

	for t.peek(' ') {
		_ = t.advance()
	}

	var token token.Token

	switch current := t.advance(); current {
	case '/':
		regex, err := t.regex()
		if err != nil {
			return NO_TOKEN, err
		}
		token = regex
	case '"':
		quoted, err := t.value()
		if err != nil {
			return NO_TOKEN, err
		}
		token = quoted
	default:
		token = t.identifier(current)
	}

	return token, nil
}

func (t *Tokenizer) regex() (_ token.Token, err error) {
	start := t.current

	defer func() {
		if err != nil {
			t.current = start
		}
	}()

	var regex string

	for !t.peek('/') && !t.IsAtEnd() {
		current := t.advance()
		regex += string(current)
	}

	if !t.match('/') {
		return NO_TOKEN, ErrUnterminatedRegex
	}
	if regex == "" {
		return NO_TOKEN, ErrEmptyRegex
	}

	return token.Token{Value: regex, Kind: token.Regex}, nil
}

func (t *Tokenizer) escape() (_ string, err error) {
	start := t.current

	defer func() {
		if err != nil {
			t.current = start
		}
	}()

	switch current := t.advance(); current {
	case '"':
		return "\"", nil
	case 'n':
		return "\n", nil
	case 'r':
		return "", nil
	default:
		return "", ErrInvalidEscape
	}
}

func (t *Tokenizer) value() (_ token.Token, err error) {
	start := t.current

	defer func() {
		if err != nil {
			t.current = start
		}
	}()

	var quoted string

	for !t.peek('"') && !t.IsAtEnd() {
		switch current := t.advance(); current {
		case '\\':
			escape, err := t.escape()
			if err != nil {
				return NO_TOKEN, err
			}
			quoted += string(escape)
		default:
			quoted += string(current)
		}
	}

	if !t.match('"') {
		return NO_TOKEN, ErrUnterminatedValue
	}
	if quoted == "" {
		return NO_TOKEN, ErrEmptyValue
	}

	return token.Token{Value: quoted, Kind: token.Value}, nil
}

func (t *Tokenizer) identifier(start byte) token.Token {
	var identifier = string(start)
	for !t.IsAtEnd() && !t.peek(' ') {
		current := t.advance()
		identifier += string(current)
	}
	switch identifier {
	case "do":
		return token.Token{Value: identifier, Kind: token.Do}
	case "end":
		return token.Token{Value: identifier, Kind: token.End}
	}
	return token.Token{Value: identifier, Kind: token.Identifier}
}

func (t *Tokenizer) Mark() int {
	return t.current
}

func (t *Tokenizer) Reset(i int) {
	t.current = i
}

func (t *Tokenizer) IsAtEnd() bool {
	return t.current >= len(t.source)
}

func (t *Tokenizer) peek(b byte) bool {
	if t.IsAtEnd() {
		return false
	}
	current := t.source[t.current]
	return current == b
}

func (t *Tokenizer) advance() byte {
	current := t.source[t.current]
	t.current += 1
	return current
}

func (t *Tokenizer) match(b byte) bool {
	if !t.peek(b) {
		return false
	}
	_ = t.advance()
	return true
}

func (te Error) Error() string {
	msg, ok := tokenErrMsg[te]
	if !ok {
		return "unexpected tokenizer error"
	}
	return msg
}
