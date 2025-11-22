package lexer

import (
	"fmt"
	"os"
	"strconv"
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
	case '(':
		tok = Token{Type: LPAREN, Literal: string(lx.ch), Line: lx.line}
	case ')':
		tok = Token{Type: RPAREN, Literal: string(lx.ch), Line: lx.line}
	case '+':
		tok = Token{Type: PLUS, Literal: string(lx.ch), Line: lx.line}
	case '-':
		tok = Token{Type: MINUS, Literal: string(lx.ch), Line: lx.line}
	case '*':
		tok = Token{Type: MULT, Literal: string(lx.ch), Line: lx.line}
	case '/':
		tok = Token{Type: DIV, Literal: string(lx.ch), Line: lx.line}
	case '"':
		tok.Literal = lx.readContentInsideDoubleQuote()
		tok.Type = STRING
		return tok
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
			literal := lx.readNumeric()
			tok.Type = DetectNumberType(literal)
			tok.Literal = literal
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
	first_i := lx.pos+1;
	if lx.ch != '"'{
		fmt.Println("[Error] invalid string format at line: ", lx.line)
		os.Exit(1)
	}
	lx.moveCursorToRight()
	for lx.ch != '"' && lx.ch != 0 {
		lx.moveCursorToRight()
	}
	lx.moveCursorToRight()
	return lx.input[first_i:lx.pos-1]
}

func (lx *Lexer) readIdentifier() string{
	first_i := lx.pos;
	for isLetter(lx.ch) || isDigit(lx.ch){
		lx.moveCursorToRight();
	}
	return lx.input[first_i:lx.pos]
}

func (lx *Lexer) readNumeric() string {
    start := lx.pos

    hasDot := false
    hasExp := false

    for {
        ch := lx.ch

        if isDigit(ch) {
            lx.moveCursorToRight()
            continue
        }

        if ch == '.' && !hasDot && !hasExp {
            hasDot = true
            lx.moveCursorToRight()
            continue
        }

        if (ch == 'e' || ch == 'E') && !hasExp {
            hasExp = true
            lx.moveCursorToRight()

            if lx.ch == '+' || lx.ch == '-' {
                lx.moveCursorToRight()
            }
            continue
        }

        if ch == 'L' || ch == 'l' || ch == 'F' || ch == 'f' || ch == 'D' || ch == 'd' {
            lx.moveCursorToRight()
            break
        }

        break
    }

    return lx.input[start:lx.pos]
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

func DetectNumberType(lit string) TokenType {
    if len(lit) == 0 {
        return ILLEGAL
    }

    // check suffix
    last := lit[len(lit)-1]
    hasLongSuffix := last == 'L' || last == 'l'
    hasFloatSuffix := last == 'F' || last == 'f'
    hasDoubleSuffix := last == 'D' || last == 'd'

    // remove suffix
    num := lit
    if hasLongSuffix || hasFloatSuffix || hasDoubleSuffix {
        num = lit[:len(lit)-1]
    }

    hasDot := false
    hasExp := false
    expIndex := -1

    for i := 0; i < len(num); i++ {
        c := num[i]

        switch {
        case c >= '0' && c <= '9':
            continue
		
        case c == '.':
            if hasDot || hasExp {
				// why u want many dot fam?
                return ILLEGAL
            }
            hasDot = true

        case c == 'e' || c == 'E':
            if hasExp || i == 0 { 
				// why u need many e also fam?
                return ILLEGAL
            }
            hasExp = true
            expIndex = i

            if i+1 >= len(num) {
                return ILLEGAL
            }
            if num[i+1] == '+' || num[i+1] == '-' {
                if i+2 >= len(num) {
                    return ILLEGAL
                }
            }

        case c == '+' || c == '-':
            if i != expIndex+1 {
                return ILLEGAL
            }

        default:
            return ILLEGAL
        }
    }

    if hasFloatSuffix {
        return FLOAT
    }
    if hasDoubleSuffix {
        return DOUBLE
    }
    if hasLongSuffix {
        return LONG
    }

    // check double
    if hasDot {
        return DOUBLE
    }
    if hasExp {
        return DOUBLE
    }

    // check byte
    if n, err := strconv.Atoi(num); err == nil && n >= 0 && n <= 255 {
        return BYTE
    }

    return INT
}
