package lexer

type TokenType string

const (
	TYPE_INT     TokenType = "TYPE_INT"
	TYPE_STRING  TokenType = "TYPE_STRING"
	TYPE_BOOLEAN TokenType = "TYPE_BOOLEAN"

	INT    TokenType = "INT"
	IDENT  TokenType = "IDENT"
	STRING TokenType = "STRING"

	ASSIGN TokenType = "="
	PLUS   TokenType = "+"
	MINUS  TokenType = "-"
	MULT   TokenType = "*"
	DIV    TokenType = "/"

	LPAREN   TokenType = "("
	RPAREN   TokenType = ")"
	LBRACE   TokenType = "{"
	RBRACE   TokenType = "}"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"
	COMMA    TokenType = ","
	SEMI     TokenType = ";"
	DoubleQ  TokenType = "\""

	DOT   TokenType = "."
	NEW   TokenType = "NEW"
	CLASS TokenType = "CLASS"
	VOID  TokenType = "VOID"
	FUNC  TokenType = "FUNC"

	EOF           TokenType = "EOF"
	ILLEGAL       TokenType = "ILLEGAL"
	NATIVE        TokenType = "NATIVE"
	IMPORT        TokenType = "IMPORT"
	NATIVE_IMPORT TokenType = "NATIVE_IMPORT"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

var keywords = map[string]TokenType{
	"int":           TYPE_INT,
	"import":        IMPORT,
	"native_import": NATIVE_IMPORT,
	"string":        TYPE_STRING,
	"bool":          TYPE_BOOLEAN,
	"class":         CLASS,
	"new":           NEW,
	"void":          VOID,
	"func":          FUNC,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
