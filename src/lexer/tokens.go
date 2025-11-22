package lexer

type TokenType string;

const (
	TYPE_INT TokenType = "TYPE_INT"
	TYPE_STRING TokenType = "TYPE_STRING"
	TYPE_BOOLEAN TokenType = "TYPE_BOOLEAN"
	TYPE_FLOAT TokenType = "TYPE_FLOAT"
	TYPE_DOUBLE TokenType = "TYPE_DOUBLE"
	TYPE_LONG TokenType = "TYPE_LONG"
	TYPE_BYTE TokenType = "TYPE_BYTE"

	INT TokenType = "INT"
	IDENT TokenType = "IDENT"
	STRING TokenType = "STRING"
	FLOAT TokenType = "FLOAT"
	DOUBLE TokenType = "DOUBLE"
	LONG TokenType = "LONG"
	BYTE   TokenType = "BYTE"
	BOOLEAN TokenType = "BOOLEAN"

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
	"float": TYPE_FLOAT,
	"double": TYPE_DOUBLE,
	"long": TYPE_LONG,
	"byte": TYPE_BYTE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}