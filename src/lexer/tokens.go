package lexer

type TokenType string

const (
	IF			 TokenType = "IF"
	ELSEIF			 TokenType = "ELSEIF"
	ELSE			 TokenType = "ELSE"

	TYPE_INT     TokenType = "TYPE_INT"
	TYPE_STRING  TokenType = "TYPE_STRING"
	TYPE_BOOLEAN TokenType = "TYPE_BOOLEAN"
	TYPE_FLOAT   TokenType = "TYPE_FLOAT"
	TYPE_DOUBLE  TokenType = "TYPE_DOUBLE"
	TYPE_LONG    TokenType = "TYPE_LONG"
	TYPE_BYTE    TokenType = "TYPE_BYTE"

	INT     TokenType = "INT"
	IDENT   TokenType = "IDENT"
	STRING  TokenType = "STRING"
	FLOAT   TokenType = "FLOAT"
	DOUBLE  TokenType = "DOUBLE"
	LONG    TokenType = "LONG"
	BYTE    TokenType = "BYTE"
	BOOLEAN TokenType = "BOOLEAN"

	ASSIGN TokenType = "="
	PLUS   TokenType = "+"
	MINUS  TokenType = "-"
	MULT   TokenType = "*"
	DIV    TokenType = "/"

	EQ     TokenType = "=="
	NEQ    TokenType = "!="
	LT     TokenType = "<" 
	GT     TokenType = ">"
	LTE    TokenType = "<="
	GTE    TokenType = ">="

	AND    TokenType = "&&"
	OR     TokenType = "||"
	NOT    TokenType = "!" 

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

	COMMENT       TokenType = "COMMENT"
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
	"float":         TYPE_FLOAT,
	"double":        TYPE_DOUBLE,
	"long":          TYPE_LONG,
	"byte":          TYPE_BYTE,
	"if":			 IF,
	"else if":       ELSEIF,
	"else":          ELSE,
	"//":			 COMMENT,
}

var replaceWords = map[string]string{
	"null":            "nil",
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func LookupReplaceWord(key string) string {
	if tok, ok := replaceWords[key]; ok {
		return tok
	}
	return key
}