package parsing

import (
	"fmt"
	"strconv"
	"strings"

	lx "github.com/B1gdawg0/Gosh/src/lexer"
)

var precedence = map[lx.TokenType]int{
	lx.OR:    1,
	lx.AND:   2,
	lx.EQ:    3,
	lx.NEQ:   3,
	lx.LT:    3,
	lx.GT:    3,
	lx.LTE:   3,
	lx.GTE:   3,
	lx.PLUS:  4,
	lx.MINUS: 4,
	lx.MULT:  5,
	lx.DIV:   5,
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

type BlockExpr struct {
	Body []Expr
}

type IfExpr struct {
	Condition Expr
	Then      Expr
	Else      Expr
}

type UnaryExpr struct {
	Op    lx.TokenType
	Right Expr
}

func ParseExpr(lexer *lx.Lexer) Expr {
	return parseBinaryExpr(lexer, 0)
}

func parseBinaryExpr(lexer *lx.Lexer, minPrec int) Expr {
	left := parsePrimary(lexer)

	for {
		op := lexer.Tokenize()

		if op.Type == lx.RPAREN || op.Type == lx.SEMI || op.Type == lx.EOF || op.Type == lx.LBRACE {
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

func ParseIfExpr(lexer *lx.Lexer) *IfExpr {
    var cond Expr

    tok := lexer.PeekToken()
    if tok.Type == lx.LPAREN {
        lexer.Tokenize()
        cond = ParseExpr(lexer)
        tok = lexer.Tokenize()
        if tok.Type != lx.RPAREN {
            panic(fmt.Sprintf("[Error] Expected ')' after if condition at line %d", tok.Line))
        }
    } else {
        cond = ParseExpr(lexer)
    }

    tok = lexer.Tokenize()
    if tok.Type != lx.LBRACE {
        panic(fmt.Sprintf("[Error] Expected '{' after if condition at line %d", tok.Line))
    }

    thenBlock := parseBlockExpr(lexer)

    lexer.SkipWhiteSpace()
    tok = lexer.PeekToken()
    if tok.Type == lx.ELSE {
        lexer.Tokenize()
        next := lexer.PeekToken()
        switch next.Type {
			case lx.IF:
				lexer.Tokenize()
				return &IfExpr{
					Condition: cond,
					Then:      thenBlock,
					Else:      ParseIfExpr(lexer),
				}
        	case lx.LBRACE:
				lexer.Tokenize()
				elseBlock := parseBlockExpr(lexer)
				return &IfExpr{
					Condition: cond,
					Then:      thenBlock,
					Else:      elseBlock,
				}
			default:
				panic(fmt.Sprintf("[Error] Expected 'if' or '{' after else at line %d", next.Line))
        }
    }

    return &IfExpr{
        Condition: cond,
        Then:      thenBlock,
        Else:      nil,
    }
}

func parseBlockExpr(lexer *lx.Lexer) Expr {
	var body []Expr
	for {
		tok := lexer.PeekToken()
		if tok.Type == lx.RBRACE {
			lexer.Tokenize()
			break
		}
		body = append(body, parseStatement(lexer))
	}
	return &BlockExpr{Body: body}
}


func parseStatement(lexer *lx.Lexer) Expr {
	tok := lexer.PeekToken()

	switch tok.Type {
	case lx.TYPE_INT, lx.TYPE_LONG, lx.TYPE_FLOAT, lx.TYPE_DOUBLE,
		lx.TYPE_BYTE, lx.TYPE_STRING, lx.TYPE_BOOLEAN:
		return parseVarDecl(lexer)
	case lx.IF:
		lexer.Tokenize()
		return ParseIfExpr(lexer)
	default:
		expr := ParseExpr(lexer)
		semi := lexer.Tokenize()
		if semi.Type != lx.SEMI {
			panic(fmt.Sprintf("[Error] Expected ';' after expression at line %d", semi.Line))
		}
		return expr
	}
}

func parseVarDecl(lexer *lx.Lexer) Expr {
	typeTok := lexer.Tokenize()
	var varType lx.TokenType = typeTok.Type

	nameTok := lexer.Tokenize()
	if nameTok.Type != lx.IDENT {
		panic(fmt.Sprintf("[Error] Expected variable name after type at line %d", nameTok.Line))
	}
	varName := nameTok.Literal
	assignTok := lexer.Tokenize()
	if assignTok.Type != lx.ASSIGN {
		panic(fmt.Sprintf("[Error] Expected '=' after variable name at line %d", assignTok.Line))
	}

	valueExpr := ParseExpr(lexer)

	semi := lexer.Tokenize()
	if semi.Type != lx.SEMI {
		panic(fmt.Sprintf("[Error] Expected ';' after expression at line %d", semi.Line))
	}

	return &VarDecl{
		Name:  varName,
		Type:  varType,
		Value: valueExpr,
	}
}

func parsePrimary(lexer *lx.Lexer) Expr {
	tok := lexer.Tokenize()
	switch tok.Type {
	case lx.INT, lx.BYTE, lx.FLOAT, lx.DOUBLE, lx.LONG:
		return &NumericExpr{Raw: tok.Literal}
	case lx.NOT:
		operand := parsePrimary(lexer)
		return &UnaryExpr{
			Op:    tok.Type,
			Right: operand,
		}
	case lx.IF:
		return ParseIfExpr(lexer)
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
	return TranspileExprWithType(e, "")
}

func TranspileExprWithType(e Expr, expectedType string) string {
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
	case *IfExpr:
		s := "if " + TranspileExpr(v.Condition) + " "
		s += TranspileExpr(v.Then)

		if v.Else != nil {
			if elseIf, ok := v.Else.(*IfExpr); ok {
				s += " else " + TranspileExpr(elseIf)
			} else {
				s += " else " + TranspileExpr(v.Else)
			}
		}
		return s

	case *BlockExpr:
		s := "{\n"
		for _, stmt := range v.Body {
			s += "\t" + TranspileExpr(stmt) + "\n"
		}
		s += "}"
		return s
	
	case *UnaryExpr:
		opStr := ""
		switch v.Op {
		case lx.NOT:
			opStr = "!"
		}
		return opStr + TranspileExpr(v.Right)

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
					// Get expected type from parameter
					var argExpectedType string
					if i < len(method.Params) {
						param := method.Params[i]
						if param.Type == lx.IDENT && param.TypeName != "" {
							if strings.HasPrefix(param.TypeName, "[]") {
								baseType := strings.TrimPrefix(param.TypeName, "[]")
								if strings.Contains(baseType, ".") || strings.HasPrefix(baseType, "*") {
									argExpectedType = "[]" + baseType
								} else {
									argExpectedType = "[]*" + baseType
								}
							} else {
								if strings.Contains(param.TypeName, ".") || strings.HasPrefix(param.TypeName, "*") {
									argExpectedType = param.TypeName
								} else {
									argExpectedType = "*" + param.TypeName
								}
							}
						}
					}
					args += TranspileExprWithType(arg, argExpectedType)
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
				// Get expected type from field
				var argExpectedType string
				if i < len(class.Fields) {
					field := class.Fields[i]
					if field.Type == lx.IDENT && field.TypeName != "" {
						if strings.HasPrefix(field.TypeName, "[]") {
							baseType := strings.TrimPrefix(field.TypeName, "[]")
							if strings.Contains(baseType, ".") || strings.HasPrefix(baseType, "*") {
								argExpectedType = "[]" + baseType
							} else {
								argExpectedType = "[]*" + baseType
							}
						} else {
							if strings.Contains(field.TypeName, ".") || strings.HasPrefix(field.TypeName, "*") {
								argExpectedType = field.TypeName
							} else {
								argExpectedType = "*" + field.TypeName
							}
						}
					}
				}
				fields += fmt.Sprintf("%s: %s", class.Fields[i].Name, TranspileExprWithType(arg, argExpectedType))
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
		// Use expected type if provided (from constructor parameter or field type)
		if expectedType != "" && strings.HasPrefix(expectedType, "[]") {
			return fmt.Sprintf("%s{%s}", expectedType, elements)
		}
		// Otherwise, try to infer from elements
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
		panic(fmt.Sprintf("[Error] Unknown expr type: %T", e))
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
