package token

type Token struct {
	Value string
	Token Type
}

type Type int8

const (
	Newline Type = iota
	Comma
	Null
	Value
	EscapedValue
)

func (f Token) EqualIgnoreValueIfSeparator(other Token) bool {
	matchComma := other.Token == Comma
	matchNewline := other.Token == Newline
	matchSeparator := matchComma || matchNewline
	matchToken := other.Token == f.Token
	matchValue := other.Value == f.Value
	return matchToken && matchValue ||
		matchSeparator && matchToken
}
