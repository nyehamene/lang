package csv

import (
	"errors"
	"reflect"
	"testing"

	"github.com/nyehamene/lang/csv/parser"
	"github.com/nyehamene/lang/csv/token"
	"github.com/nyehamene/lang/csv/tokenizer"
)

var testdata = []struct {
	Name    string
	Source  string
	Fields  []token.Token
	Records []parser.Record
	Err     error
}{
	{
		Name:   "first",
		Source: "x",
		Fields: []token.Token{
			{Value: "x", Token: token.Value},
		},
		Records: []parser.Record{
			{token.Token{Value: "x", Token: token.Value}},
		},
		Err: nil,
	},

	{
		Name:   "first quoted",
		Source: `"x"`,
		Fields: []token.Token{
			{Value: "x", Token: token.EscapedValue},
		},
		Records: []parser.Record{
			{token.Token{Value: "x", Token: token.EscapedValue}},
		},
		Err: nil,
	},

	{
		Name:   "second",
		Source: "x,y",
		Fields: []token.Token{
			{Value: "x", Token: token.Value},
			{Value: ",", Token: token.Comma},
			{Value: "y", Token: token.Value},
		},
		Records: []parser.Record{
			{
				token.Token{Value: "x", Token: token.Value},
				token.Token{Value: "y", Token: token.Value},
			},
		},
		Err: nil,
	},

	{
		Name:   "second quoted",
		Source: `"x","y"`,
		Fields: []token.Token{
			{Value: "x", Token: token.EscapedValue},
			{Value: ",", Token: token.Comma},
			{Value: "y", Token: token.EscapedValue},
		},
		Records: []parser.Record{
			{
				token.Token{Value: "x", Token: token.EscapedValue},
				token.Token{Value: "y", Token: token.EscapedValue},
			},
		},
		Err: nil,
	},

	{
		Name:   "third",
		Source: ",,y",
		Fields: []token.Token{
			{Value: ",", Token: token.Comma},
			{Value: ",", Token: token.Comma},
			{Value: "y", Token: token.Value},
		},
		Records: []parser.Record{
			{
				token.Token{Value: "", Token: token.Null},
				token.Token{Value: "", Token: token.Null},
				token.Token{Value: "y", Token: token.Value},
			},
		},
		Err: nil,
	},

	{
		Name:   "third quoted",
		Source: `,,"y"`,
		Fields: []token.Token{
			{Value: ",", Token: token.Comma},
			{Value: ",", Token: token.Comma},
			{Value: "y", Token: token.EscapedValue},
		},
		Records: []parser.Record{
			{
				token.Token{Value: "", Token: token.Null},
				token.Token{Value: "", Token: token.Null},
				token.Token{Value: "y", Token: token.EscapedValue},
			},
		},
		Err: nil,
	},

	{
		Name:   "fourth",
		Source: "x,y\n",
		Fields: []token.Token{
			{Value: "x", Token: token.Value},
			{Value: ",", Token: token.Comma},
			{Value: "y", Token: token.Value},
			{Value: "\n", Token: token.Newline},
		},
		Records: []parser.Record{
			{
				token.Token{Value: "x", Token: token.Value},
				token.Token{Value: "y", Token: token.Value},
			},
		},
		Err: nil,
	},

	{
		Name:   "fourth quoted",
		Source: "\"x\",\"y\"\n",
		Fields: []token.Token{
			{Value: "x", Token: token.EscapedValue},
			{Value: ",", Token: token.Comma},
			{Value: "y", Token: token.EscapedValue},
			{Value: "\n", Token: token.Newline},
		},
		Records: []parser.Record{
			{
				token.Token{Value: "x", Token: token.EscapedValue},
				token.Token{Value: "y", Token: token.EscapedValue},
			},
		},
		Err: nil,
	},

	{
		Name:   "fifth",
		Source: "x,y\na,b\n",
		Fields: []token.Token{
			{Value: "x", Token: token.Value},
			{Value: ",", Token: token.Comma},
			{Value: "y", Token: token.Value},
			{Value: "\n", Token: token.Newline},
			{Value: "a", Token: token.Value},
			{Value: ",", Token: token.Comma},
			{Value: "b", Token: token.Value},
			{Value: "\n", Token: token.Newline},
		},
		Records: []parser.Record{
			{
				token.Token{Value: "x", Token: token.Value},
				token.Token{Value: "y", Token: token.Value},
			},
			{
				token.Token{Value: "a", Token: token.Value},
				token.Token{Value: "b", Token: token.Value},
			},
		},
		Err: nil,
	},

	{
		Name:   "fifth quoted",
		Source: "\"x\",\"y\"\n\"a\",\"b\"\n",
		Fields: []token.Token{
			{Value: "x", Token: token.EscapedValue},
			{Value: ",", Token: token.Comma},
			{Value: "y", Token: token.EscapedValue},
			{Value: "\n", Token: token.Newline},
			{Value: "a", Token: token.EscapedValue},
			{Value: ",", Token: token.Comma},
			{Value: "b", Token: token.EscapedValue},
			{Value: "\n", Token: token.Newline},
		},
		Records: []parser.Record{
			{
				token.Token{Value: "x", Token: token.EscapedValue},
				token.Token{Value: "y", Token: token.EscapedValue},
			},
			{
				token.Token{Value: "a", Token: token.EscapedValue},
				token.Token{Value: "b", Token: token.EscapedValue},
			},
		},
		Err: nil,
	},

	{
		Name:   "sixth",
		Source: "x,\na,b",
		Fields: []token.Token{
			{Value: "x", Token: token.Value},
			{Value: ",", Token: token.Comma},
			{Value: "\n", Token: token.Newline},
			{Value: "a", Token: token.Value},
			{Value: ",", Token: token.Comma},
			{Value: "b", Token: token.Value},
		},
		Records: []parser.Record{
			{
				token.Token{Value: "x", Token: token.Value},
				token.Token{Token: token.Null},
			},
			{
				token.Token{Value: "a", Token: token.Value},
				token.Token{Value: "b", Token: token.Value},
			},
		},
		Err: nil,
	},

	{
		Name:   "sixth quoted",
		Source: "\"x\",\n\"a\",\"b\"",
		Fields: []token.Token{
			{Value: "x", Token: token.EscapedValue},
			{Value: ",", Token: token.Comma},
			{Value: "\n", Token: token.Newline},
			{Value: "a", Token: token.EscapedValue},
			{Value: ",", Token: token.Comma},
			{Value: "b", Token: token.EscapedValue},
		},
		Records: []parser.Record{
			{
				token.Token{Value: "x", Token: token.EscapedValue},
				token.Token{Token: token.Null},
			},
			{
				token.Token{Value: "a", Token: token.EscapedValue},
				token.Token{Value: "b", Token: token.EscapedValue},
			},
		},
		Err: nil,
	},

	{
		Name:   "seventh",
		Source: "x,y\na,",
		Fields: []token.Token{
			{Value: "x", Token: token.Value},
			{Value: ",", Token: token.Comma},
			{Value: "y", Token: token.Value},
			{Value: "\n", Token: token.Newline},
			{Value: "a", Token: token.Value},
			{Value: ",", Token: token.Comma},
		},
		Records: []parser.Record{
			{
				token.Token{Value: "x", Token: token.Value},
				token.Token{Value: "y", Token: token.Value},
			},
			{
				token.Token{Value: "a", Token: token.Value},
				token.Token{Token: token.Null},
			},
		},
		Err: nil,
	},

	{
		Name:   "seventh quoted",
		Source: "\"x\",\"y\"\n\"a\",",
		Fields: []token.Token{
			{Value: "x", Token: token.EscapedValue},
			{Value: ",", Token: token.Comma},
			{Value: "y", Token: token.EscapedValue},
			{Value: "\n", Token: token.Newline},
			{Value: "a", Token: token.EscapedValue},
			{Value: ",", Token: token.Comma},
		},
		Records: []parser.Record{
			{
				token.Token{Value: "x", Token: token.EscapedValue},
				token.Token{Value: "y", Token: token.EscapedValue},
			},
			{
				token.Token{Value: "a", Token: token.EscapedValue},
				token.Token{Token: token.Null},
			},
		},
		Err: nil,
	},

	{
		Name:   "eight",
		Source: "x,y\na,\n",
		Fields: []token.Token{
			{Value: "x", Token: token.Value},
			{Value: ",", Token: token.Comma},
			{Value: "y", Token: token.Value},
			{Value: "\n", Token: token.Newline},
			{Value: "a", Token: token.Value},
			{Value: ",", Token: token.Comma},
			{Value: "\n", Token: token.Newline},
		},
		Records: []parser.Record{
			{
				token.Token{Value: "x", Token: token.Value},
				token.Token{Value: "y", Token: token.Value},
			},
			{
				token.Token{Value: "a", Token: token.Value},
				token.Token{Token: token.Null},
			},
		},
		Err: nil,
	},

	{
		Name:   "eight quoted",
		Source: "\"x\",\"y\"\n\"a\",\n",
		Fields: []token.Token{
			{Value: "x", Token: token.EscapedValue},
			{Value: ",", Token: token.Comma},
			{Value: "y", Token: token.EscapedValue},
			{Value: "\n", Token: token.Newline},
			{Value: "a", Token: token.EscapedValue},
			{Value: ",", Token: token.Comma},
			{Value: "\n", Token: token.Newline},
		},
		Records: []parser.Record{
			{
				token.Token{Value: "x", Token: token.EscapedValue},
				token.Token{Value: "y", Token: token.EscapedValue},
			},
			{
				token.Token{Value: "a", Token: token.EscapedValue},
				token.Token{Token: token.Null},
			},
		},
		Err: nil,
	},

	{
		Name:   "ninth",
		Source: "x,\"y\"\na,b\n",
		Fields: []token.Token{
			{Value: "x", Token: token.Value},
			{Value: ",", Token: token.Comma},
			{Value: "y", Token: token.EscapedValue},
			{Value: "\n", Token: token.Newline},
			{Value: "a", Token: token.Value},
			{Value: ",", Token: token.Comma},
			{Value: "b", Token: token.Value},
			{Value: "\n", Token: token.Newline},
		},
		Records: []parser.Record{
			{
				token.Token{Value: "x", Token: token.Value},
				token.Token{Value: "y", Token: token.EscapedValue},
			},
			{
				token.Token{Value: "a", Token: token.Value},
				token.Token{Value: "b", Token: token.Value},
			},
		},
		Err: nil,
	},

	{
		Name:   "CRLF",
		Source: "\r\n",
		Fields: []token.Token{
			{Value: "\n", Token: token.Newline},
		},
		Records: []parser.Record{
			{token.Token{Value: "\n", Token: token.Newline}},
		},
		Err: nil,
	},

	{
		Name:   "LF",
		Source: "\n",
		Fields: []token.Token{
			{Value: "\n", Token: token.Newline},
		},
		Records: []parser.Record{
			{token.Token{Value: "\n", Token: token.Newline}},
		},
		Err: nil,
	},
}

func TestTokenizeAll(t *testing.T) {
	for _, d := range testdata {
		t.Run(d.Name, func(t *testing.T) {
			toks := tokenizer.New(d.Source)
			got, err := toks.TokenizeAll()

			if !errors.Is(err, d.Err) {
				t.Error("expected", d.Err, "but got", err)
			}
			if !reflect.DeepEqual(d.Fields, got) {
				t.Error("expected\n\t", d.Fields, "\nbut got\n\t", got)
			}
		})
	}
}

func TestParseAll(t *testing.T) {
	for _, d := range testdata {
		t.Run(d.Name, func(t *testing.T) {
			tok := tokenizer.New(d.Source)
			par := parser.New(*tok)
			got, err := par.ParseAll()
			if err != nil {
				t.Fatal(err)
			}
			if expected := d.Records; !reflect.DeepEqual(expected, got) {
				t.Error("expected\n\t", expected, "\nbut got\n\t", got)
			}
		})
	}
}

func TestAdvance(t *testing.T) {
	data := []struct {
		Name   string
		Source string
		Field  token.Token
		Err    error
	}{
		{
			Name:   "text",
			Source: "aaa",
			Field: token.Token{
				Value: "aaa",
				Token: token.Value,
			},
			Err: nil,
		},
		{
			Name:   "space1",
			Source: " aaa",
			Field: token.Token{
				Value: " aaa",
				Token: token.Value,
			},
			Err: nil,
		},
		{
			Name:   "space2",
			Source: "aaa ",
			Field: token.Token{
				Value: "aaa ",
				Token: token.Value,
			},
			Err: nil,
		},
		{
			Name:   "quoted",
			Source: `"aaa"`,
			Field: token.Token{
				Value: "aaa",
				Token: token.EscapedValue,
			},
			Err: nil,
		},
		{
			Name:   "quoted quote",
			Source: `""""`,
			Field: token.Token{
				Value: `"`,
				Token: token.EscapedValue,
			},
			Err: nil,
		},
		{
			Name:   "quoted comma",
			Source: `","`,
			Field: token.Token{
				Value: ",",
				Token: token.EscapedValue,
			},
			Err: nil,
		},
		{
			Name:   "quoted line feed",
			Source: "\"\n\"",
			Field: token.Token{
				Value: "\n",
				Token: token.EscapedValue,
			},
			Err: nil,
		},
		{
			Name:   "empty quote",
			Source: `""`,
			Field: token.Token{
				Value: "",
				Token: token.EscapedValue,
			},
			Err: nil,
		},
		{
			Name:   "blank quote",
			Source: `" "`,
			Field: token.Token{
				Value: " ",
				Token: token.EscapedValue,
			},
			Err: nil,
		},

		{
			Name:   "blank",
			Source: ``,
			Field:  token.Token{},
			Err:    tokenizer.EOF,
		},
		{
			Name:   "orphan quote1",
			Source: `"`,
			Field:  token.Token{},
			Err:    tokenizer.ErrUnterminatedString,
		},
		{
			Name:   "orphan quote2",
			Source: `"a`,
			Field:  token.Token{},
			Err:    tokenizer.ErrUnterminatedString,
		},
		{
			Name:   "CR",
			Source: "\r",
			Field:  token.Token{},
			Err:    tokenizer.ErrInvalidNewline,
		},
	}

	for _, d := range data {
		t.Run(d.Name, func(t *testing.T) {
			tokenizer := tokenizer.New(d.Source)
			got, err := tokenizer.Tokenize()
			if !errors.Is(err, d.Err) {
				t.Error("expected\n\t", d.Err, "\nbut got\n\t", err)
			}
			if got != d.Field {
				t.Error("expected\n\t", d.Field, "\nbut got\n\t", got)
			}
		})
	}
}
