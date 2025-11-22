package parsing

import (
	"fmt"

	lx "github.com/B1gdawg0/Gosh/src/lexer"
)

var precedence = map[lx.TokenType]int{
    lx.PLUS:  1,
    lx.MINUS: 1,
    lx.MULT:  2,
    lx.DIV:   2,
}

type Expr interface {}

type NumericExpr struct {
	Value int
}

type IdentifierExpr struct {
	Name string
}

type StringExpr struct{
	Value string
}

type BinaryExpr struct{
	Left Expr
	Ops  lx.TokenType
	Right Expr
}

type VarDecl struct{
	Name string
	Type lx.TokenType
	Value Expr
}

func ParseExpr(lexer *lx.Lexer) Expr{
	return parseBinaryExpr(lexer, 0)
}

func parseBinaryExpr(lexer *lx.Lexer, minPrec int) Expr{
	left := parsePrimary(lexer)

    for {
        op := lexer.Tokenize()

		if op.Type == lx.RPAREN || op.Type == lx.SEMI || op.Type == lx.EOF {
            lexer.CheckPointThis(op)
            break
        }

        prec, ok := precedence[op.Type]
        if !ok || prec < minPrec {
            lexer.CheckPointThis(op)
            break
        }

        right := parseBinaryExpr(lexer, prec+1)
        left = &BinaryExpr{
            Left:  left,
            Ops:   op.Type,
            Right: right,
        }
    }
    return left
}

func parsePrimary(lexer *lx.Lexer) Expr {
    tok := lexer.Tokenize()
    switch tok.Type {
    case lx.INT:
        val := 0
        fmt.Sscanf(tok.Literal, "%d", &val)
        return &NumericExpr{Value: val}
    case lx.IDENT:
        return &IdentifierExpr{Name: tok.Literal}
	case lx.STRING:
        return &StringExpr{Value: tok.Literal}
    case lx.LPAREN:
        e := ParseExpr(lexer)
        if next := lexer.Tokenize(); next.Type != lx.RPAREN {
            panic(fmt.Sprintf("[Error] Expected ')' after variable at line: %d", tok.Line))
        }
        return e
    default:
		fmt.Println(tok.Literal)
        panic(fmt.Sprintf("[Error] Unexpected token at line: %d", tok.Line))
    }
}

func TranspileExpr(e Expr) string {
    switch v := e.(type) {
    case *NumericExpr:
        return fmt.Sprintf("%d", v.Value)
	case *StringExpr:
        return fmt.Sprintf("\"%s\"", v.Value)
    case *IdentifierExpr:
        return v.Name
    case *BinaryExpr:
        return fmt.Sprintf("(%s %s %s)",
            TranspileExpr(v.Left),
            v.Ops,
            TranspileExpr(v.Right))
    default:
        panic("Unknown expr type")
    }
}

func GetVarAndExpr(lexer *lx.Lexer) (*lx.Token, Expr) {
    left := lexer.Tokenize()
    eq := lexer.Tokenize()
    if eq.Type !=  lx.ASSIGN{
        panic(fmt.Sprintf("[Error] Expected '=' after variable at line %d", left.Line))
    }

    expr := ParseExpr(lexer)

    semi := lexer.Tokenize() 
    if semi.Type != lx.SEMI {
        panic(fmt.Sprintf("[Error] Expected ';' after expression at line %d", semi.Line))
    }
    return &left, expr
}
