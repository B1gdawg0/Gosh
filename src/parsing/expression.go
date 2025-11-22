package parsing

import (
	"fmt"
	"strconv"
	"strings"

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

type ArrayLiteralExpr struct {
	Elements []Expr
}

type LambdaExpr struct {
	Params []Param
	Body   string
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
	case lx.IDENT:
		expr := &IdentifierExpr{Name: tok.Literal}
		return parsePostfix(lexer, expr)
	case lx.NEW:
		return parseNewExpr(lexer)
	case lx.STRING:
		return &StringExpr{Value: tok.Literal}
	case lx.LPAREN:
		e := ParseExpr(lexer)
		if next := lexer.Tokenize(); next.Type != lx.RPAREN {
			panic(fmt.Sprintf("[Error] Expected ')' after variable at line: %d", tok.Line))
		}
		return parsePostfix(lexer, e)
	case lx.LBRACKET:
		return parseArrayLiteral(lexer)
	case lx.FUNC:
		return parseLambdaExpr(lexer)
	default:
		panic(fmt.Sprintf("[Error] Unexpected token '%s' at line: %d", tok.Literal, tok.Line))
	}
}

func parsePostfix(lexer *lx.Lexer, left Expr) Expr {
	for {
		tok := lexer.Tokenize()
		if tok.Type == lx.DOT {
			member := lexer.Tokenize()
			if member.Type != lx.IDENT {
				panic(fmt.Sprintf("[Error] Expected member name after '.' at line %d", member.Line))
			}
			left = &MemberAccessExpr{
				Object: left,
				Member: member.Literal,
			}
			next := lexer.Tokenize()
			if next.Type == lx.LPAREN {
				args := parseArguments(lexer)
				left = &CallExpr{
					Callee: left,
					Args:   args,
				}
			} else {
				lexer.CheckPointThis(next)
			}
		} else if tok.Type == lx.LPAREN {
			args := parseArguments(lexer)
			left = &CallExpr{
				Callee: left,
				Args:   args,
			}
		} else {
			lexer.CheckPointThis(tok)
			break
		}
	}
	return left
}

func parseNewExpr(lexer *lx.Lexer) Expr {
	className := lexer.Tokenize()
	if className.Type != lx.IDENT {
		panic(fmt.Sprintf("[Error] Expected class name after 'new' at line %d", className.Line))
	}

	peek := lexer.Tokenize()
	if peek.Type == lx.LPAREN {
		args := parseArguments(lexer)
		return &NewExpr{
			ClassName: className.Literal,
			Args:      args,
		}
	} else {
		lexer.CheckPointThis(peek)
		return &NewExpr{ClassName: className.Literal}
	}
}

func parseArguments(lexer *lx.Lexer) []Expr {
	args := []Expr{}
	for {
		tok := lexer.Tokenize()
		if tok.Type == lx.RPAREN {
			break
		}
		lexer.CheckPointThis(tok)
		arg := ParseExpr(lexer)
		args = append(args, arg)
		next := lexer.Tokenize()
		if next.Type == lx.RPAREN {
			break
		} else if next.Type == lx.COMMA {
			continue
		} else {
			panic(fmt.Sprintf("[Error] Expected ',' or ')' in argument list at line %d", next.Line))
		}
	}
	return args
}

func parseArrayLiteral(lexer *lx.Lexer) Expr {
	elements := []Expr{}
	for {
		tok := lexer.Tokenize()
		if tok.Type == lx.RBRACKET {
			break
		}
		lexer.CheckPointThis(tok)
		element := ParseExpr(lexer)
		elements = append(elements, element)
		next := lexer.Tokenize()
		if next.Type == lx.RBRACKET {
			break
		} else if next.Type != lx.COMMA {
			panic(fmt.Sprintf("[Error] Expected ',' or ']' in array literal at line %d", next.Line))
		}
	}
	return &ArrayLiteralExpr{Elements: elements}
}

func parseLambdaExpr(lexer *lx.Lexer) Expr {
	lparen := lexer.Tokenize()
	if lparen.Type != lx.LPAREN {
		panic(fmt.Sprintf("[Error] Expected '(' after 'func' at line %d", lparen.Line))
	}

	params := []Param{}
	for {
		tok := lexer.Tokenize()
		if tok.Type == lx.RPAREN {
			break
		}
		if tok.Type == lx.EOF {
			panic(fmt.Sprintf("[Error] Unexpected EOF in lambda parameters at line %d", tok.Line))
		}

		if tok.Type != lx.IDENT {
			panic(fmt.Sprintf("[Error] Expected parameter name at line %d", tok.Line))
		}

		paramName := tok.Literal
		typeTok := lexer.Tokenize()
		var paramTypeName string
		var paramType lx.TokenType
		isPointer := false

		if typeTok.Type == lx.MULT {
			isPointer = true
			typeTok = lexer.Tokenize()
		}

		if typeTok.Type == lx.TYPE_INT || typeTok.Type == lx.TYPE_STRING || typeTok.Type == lx.TYPE_BOOLEAN {
			paramType = typeTok.Type
			paramTypeName = ""
		} else if typeTok.Type == lx.IDENT {
			paramType = lx.IDENT
			peek := lexer.Tokenize()

			if peek.Type == lx.DOT {
				typePart2 := lexer.Tokenize()
				if typePart2.Type != lx.IDENT {
					panic(fmt.Sprintf("[Error] Expected type name after '.' at line %d", typePart2.Line))
				}
				qualifiedType := typeTok.Literal + "." + typePart2.Literal
				if isPointer {
					qualifiedType = "*" + qualifiedType
				}
				paramTypeName = qualifiedType
			} else if peek.Type == lx.LBRACKET {
				bracketClose := lexer.Tokenize()
				if bracketClose.Type != lx.RBRACKET {
					panic(fmt.Sprintf("[Error] Expected ']' after '[' at line %d", bracketClose.Line))
				}
				paramTypeName = "[]" + typeTok.Literal
			} else {
				lexer.CheckPointThis(peek)
				paramTypeName = typeTok.Literal
				if isPointer {
					paramTypeName = "*" + paramTypeName
				}
			}
		} else {
			panic(fmt.Sprintf("[Error] Expected type after parameter name at line %d", typeTok.Line))
		}

		params = append(params, Param{
			Name:     paramName,
			Type:     paramType,
			TypeName: paramTypeName,
		})

		next := lexer.Tokenize()
		if next.Type == lx.COMMA {
			continue
		} else if next.Type == lx.RPAREN {
			break
		} else {
			panic(fmt.Sprintf("[Error] Expected ',' or ')' after parameter at line %d", next.Line))
		}
	}

	lbrace := lexer.Tokenize()
	if lbrace.Type != lx.LBRACE {
		panic(fmt.Sprintf("[Error] Expected '{' for lambda body at line %d", lbrace.Line))
	}

	braceCount := 1
	var bodyLines []string
	for braceCount > 0 {
		tok := lexer.Tokenize()
		if tok.Type == lx.EOF {
			panic(fmt.Sprintf("[Error] Unexpected EOF in lambda body at line %d", tok.Line))
		}
		if tok.Type == lx.LBRACE {
			braceCount++
			if braceCount > 1 {
				bodyLines = append(bodyLines, "{")
			}
		} else if tok.Type == lx.RBRACE {
			braceCount--
			if braceCount > 0 {
				bodyLines = append(bodyLines, "}")
			}
		} else if tok.Type == lx.NATIVE {
			bodyLines = append(bodyLines, strings.TrimSpace(tok.Literal))
		} else if tok.Type == lx.SEMI {
			bodyLines = append(bodyLines, ";")
		} else {
			bodyLines = append(bodyLines, tok.Literal)
		}
	}

	body := strings.Join(bodyLines, " ")
	return &LambdaExpr{
		Params: params,
		Body:   body,
	}
}

var classRegistry map[string]*ClassDecl

func SetClassRegistry(classes map[string]*ClassDecl) {
	classRegistry = classes
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

	case *NewExpr:
		class, ok := classRegistry[v.ClassName]
		if !ok {
			panic(fmt.Sprintf("[Error] Unknown class: %s", v.ClassName))
		}

		// Check for a constructor method with matching parameter count
		for _, method := range class.Methods {
			if method.Name == v.ClassName && len(method.Params) == len(v.Args) {
				args := ""
				for i, arg := range v.Args {
					if i > 0 {
						args += ", "
					}
					args += TranspileExpr(arg)
				}
				return fmt.Sprintf("func() *%s { obj := &%s{}; obj.%s(%s); return obj }()", v.ClassName, v.ClassName, v.ClassName, args)
			}
		}

		// If no constructor found, use field initialization or empty struct
		if len(v.Args) > 0 {
			if len(v.Args) != len(class.Fields) {
				panic(fmt.Sprintf("[Error] Constructor for %s expects %d arguments, got %d", v.ClassName, len(class.Fields), len(v.Args)))
			}
			fields := ""
			for i, arg := range v.Args {
				if i > 0 {
					fields += ", "
				}
				fields += fmt.Sprintf("%s: %s", class.Fields[i].Name, TranspileExpr(arg))
			}
			return fmt.Sprintf("&%s{%s}", v.ClassName, fields)
		}
		return fmt.Sprintf("&%s{}", v.ClassName)
	case *MemberAccessExpr:
		return fmt.Sprintf("%s.%s", TranspileExpr(v.Object), v.Member)
	case *CallExpr:
		args := ""
		for i, arg := range v.Args {
			if i > 0 {
				args += ", "
			}
			args += TranspileExpr(arg)
		}
		return fmt.Sprintf("%s(%s)", TranspileExpr(v.Callee), args)
	case *ArrayLiteralExpr:
		elements := ""
		for i, elem := range v.Elements {
			if i > 0 {
				elements += ", "
			}
			elements += TranspileExpr(elem)
		}
		if len(v.Elements) > 0 {
			if newExpr, ok := v.Elements[0].(*NewExpr); ok {
				return fmt.Sprintf("[]*%s{%s}", newExpr.ClassName, elements)
			}
		}
		return fmt.Sprintf("[]interface{}{%s}", elements)
	case *LambdaExpr:
		params := ""
		for i, param := range v.Params {
			if i > 0 {
				params += ", "
			}
			var goType string
			if param.Type == lx.IDENT && param.TypeName != "" {
				if strings.HasPrefix(param.TypeName, "[]") {
					baseType := strings.TrimPrefix(param.TypeName, "[]")
					if strings.Contains(baseType, ".") || strings.HasPrefix(baseType, "*") {
						goType = "[]" + baseType
					} else {
						goType = "[]*" + baseType
					}
				} else {
					if strings.Contains(param.TypeName, ".") || strings.HasPrefix(param.TypeName, "*") {
						goType = param.TypeName
					} else {
						goType = "*" + param.TypeName
					}
				}
			} else {
				switch param.Type {
				case lx.TYPE_INT:
					goType = "int"
				case lx.TYPE_STRING:
					goType = "string"
				case lx.TYPE_BOOLEAN:
					goType = "bool"
				default:
					goType = "interface{}"
				}
			}
			params += fmt.Sprintf("%s %s", param.Name, goType)
		}
		return fmt.Sprintf("func(%s) { %s }", params, v.Body)
	default:
		panic(fmt.Sprintf("Unknown expr type: %T", e))
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
