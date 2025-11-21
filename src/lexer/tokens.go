package lexer

type TokenType string;

const (
)

type Token struct{
	Type TokenType
	Literal string
}