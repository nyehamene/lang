package parser

import (
	"errors"

	"github.com/nyehamene/lang/csv/token"
	"github.com/nyehamene/lang/csv/tokenizer"
)

type Parser struct {
	line      int
	tokenizer tokenizer.Tokenizer
}

type Error int8

type Record []token.Token

const (
	ErrMissingFieldSeparator Error = iota
)

var errMsg = map[Error]string{
	ErrMissingFieldSeparator: "missing field separator COMMA or CRLF",
}

var NULL = token.Token{
	Value: "",
	Token: token.Null,
}

func New(t tokenizer.Tokenizer) *Parser {
	return &Parser{
		tokenizer: t,
	}
}

func (p *Parser) ParseAll() ([]Record, error) {
	var records []Record
	for {
		record, err := p.Record()
		if errors.Is(err, tokenizer.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func (p *Parser) Record() (Record, error) {
	var mark = p.tokenizer.Mark()
	var record Record
	field, err := p.field()

	if err != nil {
		return nil, err
	}

	if field.Token == token.Comma {
		record = append(record, NULL)
		p.tokenizer.Reset(mark)
	} else {
		record = append(record, field)
	}

loop:
	for p.match(token.Comma) {
		mark := p.tokenizer.Mark()
		f, err := p.field()

		if errors.Is(err, tokenizer.EOF) {
			record = append(record, NULL)
			break
		}
		if err != nil {
			return nil, err
		}

		switch f.Token {
		case token.Comma:
			record = append(record, NULL)
			p.tokenizer.Reset(mark)
		case token.Newline:
			record = append(record, NULL)
			p.tokenizer.Reset(mark)
			break loop
		default:
			record = append(record, f)
		}
	}

	if !p.match(token.Newline) && !p.isAtEnd() {
		return nil, ErrMissingFieldSeparator
	}

	return record, nil
}

func (p *Parser) field() (token.Token, error) {
	return p.tokenizer.Tokenize()
}

func (p *Parser) isAtEnd() bool {
	return p.tokenizer.IsAtEnd()
}

func (p *Parser) peek(token token.Type) bool {
	mark := p.tokenizer.Mark()

	defer func() {
		p.tokenizer.Reset(mark)
	}()

	current, err := p.tokenizer.Tokenize()
	if err != nil {
		return false
	}
	return current.Token == token
}

func (p *Parser) match(token token.Type) bool {
	if !p.peek(token) {
		return false
	}
	_, _ = p.field()
	return true
}

func (pe Error) Error() string {
	msg, ok := errMsg[pe]
	if !ok {
		return "unexpected parse error"
	}
	return msg
}
