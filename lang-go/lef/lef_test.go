package lef_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/nyehamene/lang/lef/parser"
	"github.com/nyehamene/lang/lef/token"
	"github.com/nyehamene/lang/lef/tokenizer"
)

var tokenizerTestData = []struct {
	Name   string
	Source string
	Token  token.Token
	Err    error
}{
	{
		Name:   "identifier",
		Source: "xx",
		Token:  token.Token{Kind: token.Identifier, Value: "xx"},
	},

	{
		Name:   "value",
		Source: `"xx"`,
		Token:  token.Token{Kind: token.Value, Value: "xx"},
	},

	{
		Name:   "empty value",
		Source: `""`,
		Err:    tokenizer.ErrEmptyValue,
	},

	{
		Name:   "dquote",
		Source: `"\""`,
		Token:  token.Token{Kind: token.Value, Value: "\""},
	},

	{
		Name:   "lf",
		Source: `"\n"`,
		Token:  token.Token{Kind: token.Value, Value: "\n"},
	},

	{
		// ignore carrage return cr
		Name:   "cr",
		Source: `"\r"`,
		Err:    tokenizer.ErrEmptyValue,
	},

	{
		Name:   "crlf",
		Source: `"\r\n"`,
		Token:  token.Token{Kind: token.Value, Value: "\n"},
	},

	{
		Name:   "invalid escape",
		Source: `"\a"`,
		Err:    tokenizer.ErrInvalidEscape,
	},

	{
		Name:   "reqex",
		Source: "/x/",
		Token:  token.Token{Kind: token.Regex, Value: "x"},
	},

	{
		Name:   "unterminated regex",
		Source: "/",
		Err:    tokenizer.ErrUnterminatedRegex,
	},

	{
		Name:   "unterminated regex",
		Source: "/x",
		Err:    tokenizer.ErrUnterminatedRegex,
	},

	{
		Name:   "empty regex",
		Source: "//",
		Err:    tokenizer.ErrEmptyRegex,
	},

	{
		Name:   "unterminated value 1",
		Source: `"`,
		Err:    tokenizer.ErrUnterminatedValue,
	},

	{
		Name:   "unterminated value 2",
		Source: `"x`,
		Err:    tokenizer.ErrUnterminatedValue,
	},

	{
		Name:   "keyword do",
		Source: "do",
		Token:  token.Token{Kind: token.Do, Value: "do"},
	},

	{
		Name:   "keyword end",
		Source: "end",
		Token:  token.Token{Kind: token.End, Value: "end"},
	},
}

func TestTokenize(t *testing.T) {
	for _, d := range tokenizerTestData {
		t.Run(d.Name, func(t *testing.T) {
			tok := tokenizer.New(d.Source)
			token, err := tok.Tokenize()

			if !errors.Is(err, d.Err) {
				t.Error("expected", d.Err, "but got", err)
			}

			if token != d.Token {
				t.Error("expected", d.Token, "but got", token)
			}
		})
	}
}

var parserTestData = []struct {
	Name     string
	Source   string
	Tokens   []token.Token
	Rule     parser.Rule
	Rules    []parser.Rule
	TokenErr error
	RuleErr  error
}{
	{
		Name:   "identifier",
		Source: "n do x end",
		Tokens: []token.Token{
			{Kind: token.Identifier, Value: "n"},
			{Kind: token.Do, Value: "do"},
			{Kind: token.Identifier, Value: "x"},
			{Kind: token.End, Value: "end"},
		},
		Rule: parser.Rule{
			Name: "n",
			Value: []token.Token{
				{Kind: token.Identifier, Value: "x"},
			},
		},
	},

	{
		Name:   "quote",
		Source: "a do \"x\" end",
		Tokens: []token.Token{
			{Kind: token.Identifier, Value: "a"},
			{Kind: token.Do, Value: "do"},
			{Kind: token.Value, Value: "x"},
			{Kind: token.End, Value: "end"},
		},
		Rule: parser.Rule{
			Name: "a",
			Value: []token.Token{
				{Kind: token.Value, Value: "x"},
			},
		},
	},

	{
		Name:   "regex",
		Source: "a do /x/ end",
		Tokens: []token.Token{
			{Kind: token.Identifier, Value: "a"},
			{Kind: token.Do, Value: "do"},
			{Kind: token.Regex, Value: "x"},
			{Kind: token.End, Value: "end"},
		},
		Rule: parser.Rule{
			Name: "a",
			Value: []token.Token{
				{Kind: token.Regex, Value: "x"},
			},
		},
	},

	{
		Name:   "mix 1",
		Source: "n do y \"x\" /z/ end",
		Tokens: []token.Token{
			{Kind: token.Identifier, Value: "n"},
			{Kind: token.Do, Value: "do"},
			{Kind: token.Identifier, Value: "y"},
			{Kind: token.Value, Value: "x"},
			{Kind: token.Regex, Value: "z"},
			{Kind: token.End, Value: "end"},
		},
		Rule: parser.Rule{
			Name: "n",
			Value: []token.Token{
				{Kind: token.Identifier, Value: "y"},
				{Kind: token.Value, Value: "x"},
				{Kind: token.Regex, Value: "z"},
			},
		},
	},

	{
		Name:   "mix 2",
		Source: "n do \"x\" /z/ y end",
		Tokens: []token.Token{
			{Kind: token.Identifier, Value: "n"},
			{Kind: token.Do, Value: "do"},
			{Kind: token.Value, Value: "x"},
			{Kind: token.Regex, Value: "z"},
			{Kind: token.Identifier, Value: "y"},
			{Kind: token.End, Value: "end"},
		},
		Rule: parser.Rule{
			Name: "n",
			Value: []token.Token{
				{Kind: token.Value, Value: "x"},
				{Kind: token.Regex, Value: "z"},
				{Kind: token.Identifier, Value: "y"},
			},
		},
	},

	{
		Name:   "error 1",
		Source: "n y end",
		Tokens: []token.Token{
			{Kind: token.Identifier, Value: "n"},
			{Kind: token.Identifier, Value: "y"},
			{Kind: token.End, Value: "end"},
		},
		RuleErr: parser.ErrExpectedKeywordDo,
	},

	{
		Name:   "error 2",
		Source: "n do y",
		Tokens: []token.Token{
			{Kind: token.Identifier, Value: "n"},
			{Kind: token.Do, Value: "do"},
			{Kind: token.Identifier, Value: "y"},
		},
		RuleErr: parser.ErrExpectedKeywordEnd,
	},

	{
		Name:   "error 3",
		Source: "do y end",
		Tokens: []token.Token{
			{Kind: token.Do, Value: "do"},
			{Kind: token.Identifier, Value: "y"},
			{Kind: token.End, Value: "end"},
		},
		RuleErr: parser.ErrExpectedIdentifier,
	},

	{
		Name:   "error 4",
		Source: "a b do y end",
		Tokens: []token.Token{
			{Kind: token.Identifier, Value: "a"},
			{Kind: token.Identifier, Value: "b"},
			{Kind: token.Do, Value: "do"},
			{Kind: token.Identifier, Value: "y"},
			{Kind: token.End, Value: "end"},
		},
		RuleErr: parser.ErrExpectedKeywordDo,
	},

	{
		Name:     "error 5",
		Source:   "a do \" end",
		TokenErr: tokenizer.ErrUnterminatedValue,
	},
}

func TestTokenizeAll(t *testing.T) {
	for _, d := range parserTestData {
		t.Run(d.Name, func(t *testing.T) {
			tok := tokenizer.New(d.Source)
			tokens, err := tok.TokenizeAll()

			if !errors.Is(err, d.TokenErr) {
				t.Error("expected", d.TokenErr, "but got", err)
			}

			if !reflect.DeepEqual(tokens, d.Tokens) {
				t.Fatal("expected", d.Tokens, "but got", tokens)
			}
		})
	}
}

func TestParse(t *testing.T) {
	for _, d := range parserTestData {
		t.Run(d.Name, func(t *testing.T) {
			tok := tokenizer.New(d.Source)
			par := parser.New(*tok)
			rule, err := par.Parse()

			if !errors.Is(err, d.TokenErr) && !errors.Is(err, d.RuleErr) {
				t.Error("expected", d.RuleErr, "or", d.TokenErr, "but got", err)
			}

			if !reflect.DeepEqual(rule, d.Rule) {
				t.Error("expected", d.Rule, "but got", rule)
			}
		})
	}
}
