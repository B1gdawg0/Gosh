package lexer

type TokenType string;

const (
	TYPE_INT TokenType = "TYPE_INT"
	TYPE_STRING TokenType = "TYPE_STRING"
	TYPE_BOOLEAN TokenType = "TYPE_BOOLEAN"

	INT TokenType = "INT"
	IDENT TokenType = "IDENT"
	STRING TokenType = "STRING"

	ASSIGN TokenType = "="
	PLUS   TokenType = "+"
	MINUS  TokenType = "-"
	MULT   TokenType = "*"
	DIV    TokenType = "/"

	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	LBRACE TokenType = "{"
	RBRACE TokenType = "}"
	COMMA  TokenType = ","
	SEMI   TokenType = ";"
	DoubleQ TokenType = "\""

	EOF     TokenType = "EOF"
	ILLEGAL TokenType = "ILLEGAL"
	NATIVE  TokenType = "NATIVE"
	IMPORT TokenType  = "IMPORT"
)

type Token struct{
	Type TokenType
	Literal string
	Line int
}

var keywords = map[string]TokenType{
	"int": TYPE_INT,
	"import": IMPORT,
	"string": TYPE_STRING,
	"bool": TYPE_BOOLEAN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}