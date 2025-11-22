package parsing

import (
	"fmt"
	"strconv"

	lx "github.com/B1gdawg0/Gosh/src/lexer"
)

var precedence = map[lx.TokenType]int{
	lx.PLUS:  1,
	lx.MINUS: 1,
	lx.MULT:  2,
	lx.DIV:   2,
}

type Expr interface{}

type NumericExpr struct {
	Raw string
}

type IdentifierExpr struct {
	Name string
}

type StringExpr struct {
	Value string
}

type BooleanExpr struct {
	Value bool
}
type BinaryExpr struct {
	Left  Expr
	Ops   lx.TokenType
	Right Expr
}

type VarDecl struct {
	Name  string
	Type  lx.TokenType
	Value Expr
}

func ParseExpr(lexer *lx.Lexer) Expr {
	return parseBinaryExpr(lexer, 0)
}

func parseBinaryExpr(lexer *lx.Lexer, minPrec int) Expr {
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
	case lx.INT, lx.BYTE, lx.FLOAT, lx.DOUBLE, lx.LONG:
		return &NumericExpr{Raw: tok.Literal}
	case lx.BOOLEAN:
		b, _ := strconv.ParseBool(tok.Literal)
		return &BooleanExpr{Value: b}
	case lx.STRING:
		return &StringExpr{Value: tok.Literal}
	case lx.IDENT:
		return &IdentifierExpr{Name: tok.Literal}
	case lx.LPAREN:
		e := ParseExpr(lexer)
		if next := lexer.Tokenize(); next.Type != lx.RPAREN {
			panic(fmt.Sprintf("[Error] Expected ')' at line %d", tok.Line))
		}
		return e
	default:
		panic(fmt.Sprintf("[Error] Unexpected token '%s' at line %d", tok.Literal, tok.Line))
	}
}

func TranspileExpr(e Expr) string {
	switch v := e.(type) {

	case *NumericExpr:
		return v.Raw

	case *StringExpr:
		return fmt.Sprintf("\"%s\"", v.Value)

	case *BooleanExpr:
		if v.Value {
			return "true"
		}
		return "false"

	case *IdentifierExpr:
		return v.Name

	case *BinaryExpr:
		return fmt.Sprintf("(%s %s %s)",
			TranspileExpr(v.Left),
			v.Ops,
			TranspileExpr(v.Right),
		)

	default:
		panic("[Error] Unknown expr type in TranspileExpr")
	}
}

func GetVarAndExpr(lexer *lx.Lexer) (*lx.Token, Expr) {
	left := lexer.Tokenize()

	eq := lexer.Tokenize()
	if eq.Type != lx.ASSIGN {
		panic(fmt.Sprintf("[Error] Expected '=' after variable at line %d", left.Line))
	}

	expr := ParseExpr(lexer)

	semi := lexer.Tokenize()
	if semi.Type != lx.SEMI {
		panic(fmt.Sprintf("[Error] Expected ';' after expression at line %d", semi.Line))
	}

	return &left, expr
}
