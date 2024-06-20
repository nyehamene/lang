package token

type Token struct {
	Kind  Type
	Value string
}

type Type int

const (
	Identifier Type = iota
	Value
	DQuote
	Regex
	Do
	End
)
