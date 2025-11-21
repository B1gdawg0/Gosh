package lexer

import (
	"fmt"
	"os"
)

type Lexer struct {
	input   string
	pos     int
	readPos int
	ch      byte
	line    int

	lastToken *Token
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1}
	l.moveCursorToRight()
	return l
}

func (lx *Lexer) Tokenize() Token {
	if lx.lastToken != nil{
		t := *lx.lastToken
		lx.lastToken = nil
		return t
	}
	
	var tok Token
	lx.skipWhiteSpace()
	tok.Line = lx.line

	switch lx.ch {
	case '=':
		tok = Token{Type: ASSIGN, Literal: string(lx.ch), Line: lx.line}
	case ';':
		tok = Token{Type: SEMI, Literal: string(lx.ch), Line: lx.line}
	case '"':
		if isLetter(lx.ch){
			tok.Literal = lx.readContentInsideDoubleQuote()
			tok.Type = STRING
			return tok
		}
	case ':':
		if lx.checkRightChar() == ':' {
        lx.moveCursorToRight()
        lx.moveCursorToRight()

        start := lx.pos
        for lx.ch != 0 && lx.ch != '\n' && lx.ch != ';' {
            lx.moveCursorToRight()
        }
        literal := lx.input[start:lx.pos]
        tok = Token{Type: NATIVE, Literal: literal, Line: lx.line}
        return tok
    }
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(lx.ch) {
			tok.Literal = lx.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(lx.ch) {
			tok.Type = INT
			tok.Literal = lx.readNumber()
			return tok
		} else {
			tok = Token{Type: ILLEGAL, Literal: string(lx.ch), Line: lx.line}
		}
	}
	lx.moveCursorToRight()
	return tok
}

func (lx *Lexer) moveCursorToRight(){
	if(lx.readPos >= len(lx.input)){
		lx.ch = 0
	}else{
		lx.ch = lx.input[lx.readPos]
	}
	lx.pos = lx.readPos
	lx.readPos++
}

func (lx *Lexer) checkRightChar() byte{
	if(lx.readPos >= len(lx.input)){
		return 0
	}
	return  lx.input[lx.readPos]
}

func (lx *Lexer) skipWhiteSpace() {
	for lx.ch == ' ' || lx.ch == '\t' || lx.ch == '\n' || lx.ch == '\r'{
		if(lx.ch == '\n'){
			lx.line++
		}
		lx.moveCursorToRight()
	}
}

func (lx *Lexer) readContentInsideDoubleQuote() string{
	first_i := lx.pos;
	if lx.ch != '"'{
		fmt.Println("[Error] invalid string format at line: ", lx.line)
		os.Exit(1)
	}
	for lx.ch != '"'&&(isLetter(lx.ch) || isDigit(lx.ch)){
		lx.moveCursorToRight();
	}
	return lx.input[first_i:lx.pos]
}

func (lx *Lexer) readIdentifier() string{
	first_i := lx.pos;
	for isLetter(lx.ch) || isDigit(lx.ch){
		lx.moveCursorToRight();
	}
	return lx.input[first_i:lx.pos]
}

func (lx *Lexer) readNumber() string{
	first_n := lx.pos
	for isDigit(lx.ch){
		lx.moveCursorToRight();
	}
	return  lx.input[first_n: lx.pos];
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (lx *Lexer) CheckPointThis(tok Token){
	lx.lastToken = &tok
}