package parser

import (
	"errors"

	"github.com/nyehamene/lang/lef/token"
	"github.com/nyehamene/lang/lef/tokenizer"
)

type Parser struct {
	tokenizer tokenizer.Tokenizer
	line      int
}

type Rule struct {
	Name  string
	Value []token.Token
}

type Error int

const (
	ErrExpectedKeywordDo Error = iota
	ErrExpectedKeywordEnd
	ErrExpectedIdentifier
)

var parserErrMsg = map[Error]string{
	ErrExpectedKeywordDo:  "expected keyword 'do'",
	ErrExpectedKeywordEnd: "expected keyword 'end'",
	ErrExpectedIdentifier: "expected identifier",
}

var norule = Rule{}

func New(t tokenizer.Tokenizer) *Parser {
	return &Parser{
		tokenizer: t,
	}
}

func (t *Parser) Parse() (Rule, error) {
	return t.rule()
}

func (p *Parser) rule() (_ Rule, err error) {
	start := p.mark()

	defer func() {
		if err != nil {
			p.reset(start)
		}
	}()

	if !p.peek(token.Identifier) {
		return norule, ErrExpectedIdentifier
	}

	name, err := p.advance()

	if err != nil {
		return norule, err
	}

	if !p.expect(token.Do) {
		return norule, ErrExpectedKeywordDo
	}

	var value []token.Token
	for !p.peek(token.End) && !p.isAtEnd() {
		token, err := p.advance()
		if err != nil {
			return norule, err
		}
		value = append(value, token)
	}

	if !p.expect(token.End) {
		return norule, ErrExpectedKeywordEnd
	}

	return Rule{Name: name.Value, Value: value}, nil
}

func (p *Parser) expect(t token.Type) bool {
	if !p.peek(t) {
		return false
	}
	token, _ := p.advance()
	return token.Kind == t
}

func (p *Parser) advance() (token.Token, error) {
	return p.tokenizer.Tokenize()
}

func (p *Parser) peek(t token.Type) bool {
	start := p.mark()

	defer func() {
		p.reset(start)
	}()

	if p.isAtEnd() {
		return false
	}
	token, err := p.tokenizer.Tokenize()
	if err != nil {
		return false
	}
	return token.Kind == t
}

func (p *Parser) isAtEnd() bool {
	return p.tokenizer.IsAtEnd()
}

func (p *Parser) mark() int {
	return p.tokenizer.Mark()
}

func (p *Parser) reset(i int) {
	p.tokenizer.Reset(i)
}

func (pe Error) Error() string {
	msg, ok := parserErrMsg[ErrExpectedKeywordDo]
	if !ok {
		return "unexpected parser error"
	}
	return msg
}

func (p *Parser) error(msg string) error {
	return errors.New(msg)
}
