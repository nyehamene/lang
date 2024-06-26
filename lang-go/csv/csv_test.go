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
		Name: "value",
		Source: `a,a,a
,,
a,,
a,a,
,,a
,a,a
`,
		Fields: []token.Token{
			{Value: "a", Token: token.Value},
			{Value: ",", Token: token.Comma},
			{Value: "a", Token: token.Value},
			{Value: ",", Token: token.Comma},
			{Value: "a", Token: token.Value},
			{Value: "\n", Token: token.Newline},
			{Value: ",", Token: token.Comma},
			{Value: ",", Token: token.Comma},
			{Value: "\n", Token: token.Newline},
			{Value: "a", Token: token.Value},
			{Value: ",", Token: token.Comma},
			{Value: ",", Token: token.Comma},
			{Value: "\n", Token: token.Newline},
			{Value: "a", Token: token.Value},
			{Value: ",", Token: token.Comma},
			{Value: "a", Token: token.Value},
			{Value: ",", Token: token.Comma},
			{Value: "\n", Token: token.Newline},
			{Value: ",", Token: token.Comma},
			{Value: ",", Token: token.Comma},
			{Value: "a", Token: token.Value},
			{Value: "\n", Token: token.Newline},
			{Value: ",", Token: token.Comma},
			{Value: "a", Token: token.Value},
			{Value: ",", Token: token.Comma},
			{Value: "a", Token: token.Value},
			{Value: "\n", Token: token.Newline},
		},

		Records: []parser.Record{
			{
				{Value: "a", Token: token.Value},
				{Value: "a", Token: token.Value},
				{Value: "a", Token: token.Value},
			},

			{
				{Value: "", Token: token.Null},
				{Value: "", Token: token.Null},
				{Value: "", Token: token.Null},
			},

			{

				{Value: "a", Token: token.Value},
				{Value: "", Token: token.Null},
				{Value: "", Token: token.Null},
			},

			{
				{Value: "a", Token: token.Value},
				{Value: "a", Token: token.Value},
				{Value: "", Token: token.Null},
			},

			{
				{Value: "", Token: token.Null},
				{Value: "", Token: token.Null},
				{Value: "a", Token: token.Value},
			},

			{
				{Value: "", Token: token.Null},
				{Value: "a", Token: token.Value},
				{Value: "a", Token: token.Value},
			},
		},
		Err: nil,
	},

	{
		Name: "value",
		Source: `"a","a","a"
,,
"a",,
"a","a",
,,"a"
,"a","a"
`,
		Fields: []token.Token{
			{Value: "a", Token: token.EscapedValue},
			{Value: ",", Token: token.Comma},
			{Value: "a", Token: token.EscapedValue},
			{Value: ",", Token: token.Comma},
			{Value: "a", Token: token.EscapedValue},
			{Value: "\n", Token: token.Newline},
			{Value: ",", Token: token.Comma},
			{Value: ",", Token: token.Comma},
			{Value: "\n", Token: token.Newline},
			{Value: "a", Token: token.EscapedValue},
			{Value: ",", Token: token.Comma},
			{Value: ",", Token: token.Comma},
			{Value: "\n", Token: token.Newline},
			{Value: "a", Token: token.EscapedValue},
			{Value: ",", Token: token.Comma},
			{Value: "a", Token: token.EscapedValue},
			{Value: ",", Token: token.Comma},
			{Value: "\n", Token: token.Newline},
			{Value: ",", Token: token.Comma},
			{Value: ",", Token: token.Comma},
			{Value: "a", Token: token.EscapedValue},
			{Value: "\n", Token: token.Newline},
			{Value: ",", Token: token.Comma},
			{Value: "a", Token: token.EscapedValue},
			{Value: ",", Token: token.Comma},
			{Value: "a", Token: token.EscapedValue},
			{Value: "\n", Token: token.Newline},
		},

		Records: []parser.Record{
			{
				{Value: "a", Token: token.EscapedValue},
				{Value: "a", Token: token.EscapedValue},
				{Value: "a", Token: token.EscapedValue},
			},

			{
				{Value: "", Token: token.Null},
				{Value: "", Token: token.Null},
				{Value: "", Token: token.Null},
			},

			{

				{Value: "a", Token: token.EscapedValue},
				{Value: "", Token: token.Null},
				{Value: "", Token: token.Null},
			},

			{
				{Value: "a", Token: token.EscapedValue},
				{Value: "a", Token: token.EscapedValue},
				{Value: "", Token: token.Null},
			},

			{
				{Value: "", Token: token.Null},
				{Value: "", Token: token.Null},
				{Value: "a", Token: token.EscapedValue},
			},

			{
				{Value: "", Token: token.Null},
				{Value: "a", Token: token.EscapedValue},
				{Value: "a", Token: token.EscapedValue},
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
			{
				{Value: "\n", Token: token.Newline},
			},
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
			{
				{Value: "\n", Token: token.Newline},
			},
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
