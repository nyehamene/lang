package csv_test

import (
	"errors"
	"reflect"

	"testing"

	"github.com/nyehamene/lang/csv/parser"
	"github.com/nyehamene/lang/csv/token"
	"github.com/nyehamene/lang/csv/tokenizer"
)

var Rfc = []struct {
	Name    string
	Source  string
	Fields  []token.Token
	Records []parser.Record
	Err     error
}{
	{
		Name:   "rule1",
		Source: "aaa,bbb,ccc\nzzz,yyy,xxx\n",
		Err:    nil,
		Fields: []token.Token{
			{Value: "aaa", Token: token.Value},
			{Token: token.Comma},
			{Value: "bbb", Token: token.Value},
			{Token: token.Comma},
			{Value: "ccc", Token: token.Value},
			{Token: token.Newline},
			{Value: "zzz", Token: token.Value},
			{Token: token.Comma},
			{Value: "yyy", Token: token.Value},
			{Token: token.Comma},
			{Value: "xxx", Token: token.Value},
			{Token: token.Newline},
		},
		Records: []parser.Record{
			{
				{Value: "aaa", Token: token.Value},
				{Value: "bbb", Token: token.Value},
				{Value: "ccc", Token: token.Value},
			},
			{
				{Value: "zzz", Token: token.Value},
				{Value: "yyy", Token: token.Value},
				{Value: "xxx", Token: token.Value},
			},
		},
	},

	{
		Name:   "rule2",
		Source: "aaa,bbb,ccc\nzzz,yyy,xxx",
		Err:    nil,
		Fields: []token.Token{
			{Value: "aaa", Token: token.Value},
			{Token: token.Comma},
			{Value: "bbb", Token: token.Value},
			{Token: token.Comma},
			{Value: "ccc", Token: token.Value},
			{Token: token.Newline},
			{Value: "zzz", Token: token.Value},
			{Token: token.Comma},
			{Value: "yyy", Token: token.Value},
			{Token: token.Comma},
			{Value: "xxx", Token: token.Value},
		},
		Records: []parser.Record{
			{
				{Value: "aaa", Token: token.Value},
				{Value: "bbb", Token: token.Value},
				{Value: "ccc", Token: token.Value},
			},
			{
				{Value: "zzz", Token: token.Value},
				{Value: "yyy", Token: token.Value},
				{Value: "xxx", Token: token.Value},
			},
		},
	},

	{
		Name:   "rule4.1",
		Source: "aaa, bbb,ccc\nzzz,yyy,xxx\n",
		Err:    nil,
		Fields: []token.Token{
			{Value: "aaa", Token: token.Value},
			{Token: token.Comma},
			{Value: " bbb", Token: token.Value},
			{Token: token.Comma},
			{Value: "ccc", Token: token.Value},
			{Token: token.Newline},
			{Value: "zzz", Token: token.Value},
			{Token: token.Comma},
			{Value: "yyy", Token: token.Value},
			{Token: token.Comma},
			{Value: "xxx", Token: token.Value},
			{Token: token.Newline},
		},
		Records: []parser.Record{
			{
				{Value: "aaa", Token: token.Value},
				{Value: " bbb", Token: token.Value},
				{Value: "ccc", Token: token.Value},
			},
			{
				{Value: "zzz", Token: token.Value},
				{Value: "yyy", Token: token.Value},
				{Value: "xxx", Token: token.Value},
			}},
	},

	{
		Name:   "rule5",
		Source: "\"aaa\",\"bbb\",\"ccc\"\nzzz,yyy,xxx",
		Err:    nil,
		Fields: []token.Token{
			{Value: "aaa", Token: token.EscapedValue},
			{Token: token.Comma},
			{Value: "bbb", Token: token.EscapedValue},
			{Token: token.Comma},
			{Value: "ccc", Token: token.EscapedValue},
			{Token: token.Newline},
			{Value: "zzz", Token: token.Value},
			{Token: token.Comma},
			{Value: "yyy", Token: token.Value},
			{Token: token.Comma},
			{Value: "xxx", Token: token.Value},
		},
		Records: []parser.Record{
			{
				{Value: "aaa", Token: token.EscapedValue},
				{Value: "bbb", Token: token.EscapedValue},
				{Value: "ccc", Token: token.EscapedValue},
			},
			{
				{Value: "zzz", Token: token.Value},
				{Value: "yyy", Token: token.Value},
				{Value: "xxx", Token: token.Value},
			},
		},
	},

	{
		Name:   "rule6",
		Source: "\"aaa\",\"b\nbb\",\"ccc\"\nzzz,yyy,xxx\n",
		Err:    nil,
		Fields: []token.Token{
			{Value: "aaa", Token: token.EscapedValue},
			{Token: token.Comma},
			{Value: "b\nbb", Token: token.EscapedValue},
			{Token: token.Comma},
			{Value: "ccc", Token: token.EscapedValue},
			{Token: token.Newline},
			{Value: "zzz", Token: token.Value},
			{Token: token.Comma},
			{Value: "yyy", Token: token.Value},
			{Token: token.Comma},
			{Value: "xxx", Token: token.Value},
			{Token: token.Newline},
		},
		Records: []parser.Record{
			{
				{Value: "aaa", Token: token.EscapedValue},
				{Value: "b\nbb", Token: token.EscapedValue},
				{Value: "ccc", Token: token.EscapedValue},
			},
			{
				{Value: "zzz", Token: token.Value},
				{Value: "yyy", Token: token.Value},
				{Value: "xxx", Token: token.Value},
			},
		},
	},

	{
		Name:   "rule7",
		Source: "\"aaa\",\"b\"\"bb\",\"ccc\"\nzzz,yyy,xxx",
		Err:    nil,
		Fields: []token.Token{
			{Value: "aaa", Token: token.EscapedValue},
			{Token: token.Comma},
			{Value: "b\"bb", Token: token.EscapedValue},
			{Token: token.Comma},
			{Value: "ccc", Token: token.EscapedValue},
			{Token: token.Newline},
			{Value: "zzz", Token: token.Value},
			{Token: token.Comma},
			{Value: "yyy", Token: token.Value},
			{Token: token.Comma},
			{Value: "xxx", Token: token.Value},
		},
		Records: []parser.Record{
			{
				{Value: "aaa", Token: token.EscapedValue},
				{Value: "b\"bb", Token: token.EscapedValue},
				{Value: "ccc", Token: token.EscapedValue},
			},
			{
				{Value: "zzz", Token: token.Value},
				{Value: "yyy", Token: token.Value},
				{Value: "xxx", Token: token.Value},
			},
		},
	},
}

func TestRfc_parser(t *testing.T) {
	for _, d := range Rfc {
		t.Run(d.Name, func(t *testing.T) {
			tok := tokenizer.New(d.Source)
			par := parser.New(*tok)

			records, err := par.ParseAll()

			if !errors.Is(err, d.Err) {
				t.Error(err)
			}

			if gl, el := len(records), len(d.Records); gl != el {
				t.Error("expected", el, "fields but got", gl)
			}

			for i := 0; i < len(records); i++ {
				expected := d.Records[i]
				got := records[i]

				if !reflect.DeepEqual(expected, got) {
					t.Error("expected", expected, "but got", got)
				}
			}
		})
	}
}

func TestRfc_tokenizer(t *testing.T) {
	for _, d := range Rfc {
		t.Run(d.Name, func(t *testing.T) {
			tok := tokenizer.New(d.Source)
			fields, err := tok.TokenizeAll()

			if !errors.Is(err, d.Err) {
				t.Fatal(err)
			}

			if gl, el := len(fields), len(d.Fields); gl != el {
				t.Error("expected", el, "fields but got", gl)
			}

			for i := 0; i < len(fields); i++ {
				expected := d.Fields[i]
				got := fields[i]

				if !got.EqualIgnoreValueIfSeparator(expected) {
					t.Error("expected", expected, "but got", got)
				}
			}
		})
	}
}
