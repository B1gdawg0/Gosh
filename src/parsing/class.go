package parsing

import (
	"fmt"
	"strings"

	lx "github.com/B1gdawg0/Gosh/src/lexer"
)

type ClassDecl struct {
	Name    string
	Fields  []Field
	Methods []Method
}

type Field struct {
	Name     string
	Type     lx.TokenType
	TypeName string
}

type Method struct {
	Name       string
	ReturnType lx.TokenType
	Params     []Param
	Body       []Stmt
	BodyCode   string
}

type Param struct {
	Name     string
	Type     lx.TokenType
	TypeName string
}

type Stmt interface{}

type NewExpr struct {
	ClassName string
	Args      []Expr
}

type MemberAccessExpr struct {
	Object Expr
	Member string
}

type CallExpr struct {
	Callee Expr
	Args   []Expr
}

func ParseClass(lexer *lx.Lexer) *ClassDecl {
	classNameTok := lexer.Tokenize()
	if classNameTok.Type != lx.IDENT {
		panic(fmt.Sprintf("[Error] Expected class name at line %d", classNameTok.Line))
	}

	class := &ClassDecl{
		Name:    classNameTok.Literal,
		Fields:  []Field{},
		Methods: []Method{},
	}

	lbrace := lexer.Tokenize()
	if lbrace.Type != lx.LBRACE {
		panic(fmt.Sprintf("[Error] Expected '{' after class name at line %d", lbrace.Line))
	}

	for {
		tok := lexer.Tokenize()
		if tok.Type == lx.RBRACE {
			break
		}
		if tok.Type == lx.EOF {
			panic(fmt.Sprintf("[Error] Unexpected EOF in class declaration at line %d", tok.Line))
		}

		if tok.Type == lx.TYPE_INT || tok.Type == lx.TYPE_STRING || tok.Type == lx.TYPE_BOOLEAN ||
			tok.Type == lx.TYPE_FLOAT || tok.Type == lx.TYPE_DOUBLE || tok.Type == lx.TYPE_LONG || tok.Type == lx.TYPE_BYTE {
			nameTok := lexer.Tokenize()
			if nameTok.Type != lx.IDENT {
				panic(fmt.Sprintf("[Error] Expected identifier after type at line %d", nameTok.Line))
			}

			peek := lexer.Tokenize()

			if peek.Type == lx.LPAREN {
				method := ParseMethodWithTokens(lexer, tok.Type, nameTok, peek)
				class.Methods = append(class.Methods, *method)
			} else {
				field := &Field{
					Name: nameTok.Literal,
					Type: tok.Type,
				}
				if peek.Type != lx.SEMI {
					panic(fmt.Sprintf("[Error] Expected ';' after field declaration at line %d", peek.Line))
				}
				class.Fields = append(class.Fields, *field)
			}
		} else if tok.Type == lx.IDENT {
			peek := lexer.Tokenize()

			if peek.Type == lx.DOT {
				typePart2 := lexer.Tokenize()
				if typePart2.Type != lx.IDENT {
					panic(fmt.Sprintf("[Error] Expected type name after '.' at line %d", typePart2.Line))
				}
				qualifiedType := tok.Literal + "." + typePart2.Literal
				nextPeek := lexer.Tokenize()

				if nextPeek.Type == lx.LBRACKET {
					bracketClose := lexer.Tokenize()
					if bracketClose.Type != lx.RBRACKET {
						panic(fmt.Sprintf("[Error] Expected ']' after '[' at line %d", bracketClose.Line))
					}
					nameTok := lexer.Tokenize()
					if nameTok.Type != lx.IDENT {
						panic(fmt.Sprintf("[Error] Expected field name after '[]' at line %d", nameTok.Line))
					}
					semi := lexer.Tokenize()
					if semi.Type != lx.SEMI {
						panic(fmt.Sprintf("[Error] Expected ';' after array field at line %d", semi.Line))
					}
					field := &Field{
						Name:     nameTok.Literal,
						Type:     lx.IDENT,
						TypeName: "[]" + qualifiedType,
					}
					class.Fields = append(class.Fields, *field)
				} else if nextPeek.Type == lx.IDENT {
					nameTok := nextPeek
					nextNextPeek := lexer.Tokenize()

					if nextNextPeek.Type == lx.LPAREN {
						lexer.CheckPointThis(nextNextPeek)
						lexer.CheckPointThis(nameTok)
						lexer.CheckPointThis(nextPeek)
						lexer.CheckPointThis(typePart2)
						lexer.CheckPointThis(peek)
						lexer.CheckPointThis(tok)
						method := ParseMethod(lexer, lx.IDENT)
						class.Methods = append(class.Methods, *method)
					} else if nextNextPeek.Type == lx.SEMI {
						field := &Field{
							Name:     nameTok.Literal,
							Type:     lx.IDENT,
							TypeName: qualifiedType,
						}
						class.Fields = append(class.Fields, *field)
					} else {
						panic(fmt.Sprintf("[Error] Expected ';' or '(' after identifier at line %d", nextNextPeek.Line))
					}
				} else if nextPeek.Type == lx.SEMI {
					panic(fmt.Sprintf("[Error] Expected field name after qualified type at line %d", nextPeek.Line))
				} else {
					panic(fmt.Sprintf("[Error] Expected identifier or '[' after qualified type at line %d", nextPeek.Line))
				}
			} else if peek.Type == lx.LBRACKET {
				bracketClose := lexer.Tokenize()
				if bracketClose.Type != lx.RBRACKET {
					panic(fmt.Sprintf("[Error] Expected ']' after '[' at line %d", bracketClose.Line))
				}
				nameTok := lexer.Tokenize()
				if nameTok.Type != lx.IDENT {
					panic(fmt.Sprintf("[Error] Expected field name after '[]' at line %d", nameTok.Line))
				}
				semi := lexer.Tokenize()
				if semi.Type != lx.SEMI {
					panic(fmt.Sprintf("[Error] Expected ';' after array field at line %d", semi.Line))
				}
				field := &Field{
					Name:     nameTok.Literal,
					Type:     lx.IDENT,
					TypeName: "[]" + tok.Literal,
				}
				class.Fields = append(class.Fields, *field)
			} else if peek.Type == lx.IDENT {
				nameTok := peek
				nextPeek := lexer.Tokenize()

				if nextPeek.Type == lx.LPAREN {
					lexer.CheckPointThis(nextPeek)
					lexer.CheckPointThis(nameTok)
					lexer.CheckPointThis(tok)
					method := ParseMethod(lexer, lx.IDENT)
					class.Methods = append(class.Methods, *method)
				} else if nextPeek.Type == lx.SEMI {
					field := &Field{
						Name:     nameTok.Literal,
						Type:     lx.IDENT,
						TypeName: tok.Literal,
					}
					class.Fields = append(class.Fields, *field)
				} else {
					panic(fmt.Sprintf("[Error] Expected ';' or '(' after identifier at line %d", nextPeek.Line))
				}
			} else {
				panic(fmt.Sprintf("[Error] Expected identifier, '[', '.', or '(' after class type at line %d", peek.Line))
			}
		} else if tok.Type == lx.VOID {
			method := ParseMethod(lexer, lx.IDENT)
			class.Methods = append(class.Methods, *method)
		} else {
			panic(fmt.Sprintf("[Error] Unexpected token '%s' (type: %s) in class body at line %d", tok.Literal, tok.Type, tok.Line))
		}
	}

	return class
}

func ParseField(lexer *lx.Lexer, fieldType lx.TokenType) *Field {
	nameTok := lexer.Tokenize()
	if nameTok.Type != lx.IDENT {
		panic(fmt.Sprintf("[Error] Expected field name at line %d", nameTok.Line))
	}

	semi := lexer.Tokenize()
	if semi.Type != lx.SEMI {
		panic(fmt.Sprintf("[Error] Expected ';' after field declaration at line %d", semi.Line))
	}

	return &Field{
		Name: nameTok.Literal,
		Type: fieldType,
	}
}

func ParseMethod(lexer *lx.Lexer, returnType lx.TokenType) *Method {
	nameTok := lexer.Tokenize()
	if nameTok.Type != lx.IDENT {
		if returnType == lx.IDENT {
			lexer.CheckPointThis(nameTok)
			return nil
		}
		panic(fmt.Sprintf("[Error] Expected method name at line %d", nameTok.Line))
	}

	lparen := lexer.Tokenize()
	if lparen.Type != lx.LPAREN {
		if returnType == lx.IDENT {
			lexer.CheckPointThis(lparen)
			lexer.CheckPointThis(nameTok)
			return nil
		}
		panic(fmt.Sprintf("[Error] Expected '(' after method name at line %d", lparen.Line))
	}

	return ParseMethodWithTokens(lexer, returnType, nameTok, lparen)
}

func ParseMethodWithTokens(lexer *lx.Lexer, returnType lx.TokenType, nameTok lx.Token, lparen lx.Token) *Method {
	if lparen.Type != lx.LPAREN {
		panic(fmt.Sprintf("[Error] Expected '(' after method name at line %d", lparen.Line))
	}

	method := &Method{
		Name:       nameTok.Literal,
		ReturnType: returnType,
		Params:     []Param{},
		Body:       []Stmt{},
	}

	for {
		tok := lexer.Tokenize()
		if tok.Type == lx.RPAREN {
			break
		}
		if tok.Type == lx.EOF {
			panic(fmt.Sprintf("[Error] Unexpected EOF in method parameters at line %d", tok.Line))
		}

		if tok.Type == lx.TYPE_INT || tok.Type == lx.TYPE_STRING || tok.Type == lx.TYPE_BOOLEAN ||
			tok.Type == lx.TYPE_FLOAT || tok.Type == lx.TYPE_DOUBLE || tok.Type == lx.TYPE_LONG || tok.Type == lx.TYPE_BYTE {
			paramName := lexer.Tokenize()
			if paramName.Type != lx.IDENT {
				panic(fmt.Sprintf("[Error] Expected parameter name at line %d", paramName.Line))
			}
			method.Params = append(method.Params, Param{
				Name: paramName.Literal,
				Type: tok.Type,
			})

			next := lexer.Tokenize()
			if next.Type == lx.COMMA {
				continue
			} else if next.Type == lx.RPAREN {
				break
			} else {
				panic(fmt.Sprintf("[Error] Expected ',' or ')' after parameter at line %d", next.Line))
			}
		} else if tok.Type == lx.IDENT {
			peek := lexer.Tokenize()
			var paramTypeName string
			if peek.Type == lx.DOT {
				typePart2 := lexer.Tokenize()
				if typePart2.Type != lx.IDENT {
					panic(fmt.Sprintf("[Error] Expected type name after '.' at line %d", typePart2.Line))
				}
				qualifiedType := tok.Literal + "." + typePart2.Literal
				nextPeek := lexer.Tokenize()

				if nextPeek.Type == lx.LBRACKET {
					bracketClose := lexer.Tokenize()
					if bracketClose.Type != lx.RBRACKET {
						panic(fmt.Sprintf("[Error] Expected ']' after '[' at line %d", bracketClose.Line))
					}
					paramTypeName = "[]" + qualifiedType
					paramName := lexer.Tokenize()
					if paramName.Type != lx.IDENT {
						panic(fmt.Sprintf("[Error] Expected parameter name after '[]' at line %d", paramName.Line))
					}
					method.Params = append(method.Params, Param{
						Name:     paramName.Literal,
						Type:     lx.IDENT,
						TypeName: paramTypeName,
					})
				} else if nextPeek.Type == lx.IDENT {
					paramName := nextPeek
					paramTypeName = qualifiedType
					method.Params = append(method.Params, Param{
						Name:     paramName.Literal,
						Type:     lx.IDENT,
						TypeName: paramTypeName,
					})
				} else {
					panic(fmt.Sprintf("[Error] Expected parameter name after qualified type at line %d", nextPeek.Line))
				}
			} else if peek.Type == lx.LBRACKET {
				bracketClose := lexer.Tokenize()
				if bracketClose.Type != lx.RBRACKET {
					panic(fmt.Sprintf("[Error] Expected ']' after '[' at line %d", bracketClose.Line))
				}
				paramTypeName = "[]" + tok.Literal
				paramName := lexer.Tokenize()
				if paramName.Type != lx.IDENT {
					panic(fmt.Sprintf("[Error] Expected parameter name after '[]' at line %d", paramName.Line))
				}
				method.Params = append(method.Params, Param{
					Name:     paramName.Literal,
					Type:     lx.IDENT,
					TypeName: paramTypeName,
				})
			} else {
				lexer.CheckPointThis(peek)
				paramName := lexer.Tokenize()
				if paramName.Type != lx.IDENT {
					panic(fmt.Sprintf("[Error] Expected parameter name at line %d", paramName.Line))
				}
				paramTypeName = tok.Literal
				method.Params = append(method.Params, Param{
					Name:     paramName.Literal,
					Type:     lx.IDENT,
					TypeName: paramTypeName,
				})
			}

			next := lexer.Tokenize()
			if next.Type == lx.COMMA {
				continue
			} else if next.Type == lx.RPAREN {
				break
			} else {
				panic(fmt.Sprintf("[Error] Expected ',' or ')' after parameter at line %d", next.Line))
			}
		} else {
			if tok.Type != lx.RPAREN {
				lexer.CheckPointThis(tok)
			}
			break
		}
	}

	lbrace := lexer.Tokenize()
	if lbrace.Type != lx.LBRACE {
		panic(fmt.Sprintf("[Error] Expected '{' for method body, got '%s' (type: %s) at line %d", lbrace.Literal, lbrace.Type, lbrace.Line))
	}

	braceCount := 1
	var bodyLines []string
	for braceCount > 0 {
		tok := lexer.Tokenize()
		if tok.Type == lx.EOF {
			panic(fmt.Sprintf("[Error] Unexpected EOF in method body at line %d", tok.Line))
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
			continue
		} else {
			lexer.CheckPointThis(tok)
			stmt := TranspileMethodStatement(lexer)
			if stmt != "" {
				bodyLines = append(bodyLines, stmt)
			}
		}
	}
	method.BodyCode = strings.Join(bodyLines, "\n\t")

	return method
}

func TranspileClass(class *ClassDecl) string {
	var result string
	result += fmt.Sprintf("type %s struct {\n", class.Name)
	for _, field := range class.Fields {
		var goType string
		if field.Type == lx.IDENT && field.TypeName != "" {
			if strings.HasPrefix(field.TypeName, "[]") {
				baseType := strings.TrimPrefix(field.TypeName, "[]")
				if strings.Contains(baseType, ".") {
					goType = "[]" + baseType
				} else {
					goType = "[]*" + baseType
				}
			} else {
				if strings.Contains(field.TypeName, ".") {
					goType = field.TypeName
				} else {
					goType = "*" + field.TypeName
				}
			}
		} else {
			goType = tokenTypeToGoType(field.Type)
		}
		result += fmt.Sprintf("\t%s %s\n", field.Name, goType)
	}
	result += "}\n\n"

	for _, method := range class.Methods {
		result += TranspileMethod(class.Name, &method)
	}

	return result
}

func TranspileMethod(className string, method *Method) string {
	var result string
	result += fmt.Sprintf("func (this *%s) %s(", className, method.Name)

	for i, param := range method.Params {
		if i > 0 {
			result += ", "
		}
		var goType string
		if param.Type == lx.IDENT && param.TypeName != "" {
			if strings.HasPrefix(param.TypeName, "[]") {
				baseType := strings.TrimPrefix(param.TypeName, "[]")
				if strings.Contains(baseType, ".") {
					goType = "[]" + baseType
				} else {
					goType = "[]*" + baseType
				}
			} else {
				if strings.Contains(param.TypeName, ".") {
					goType = param.TypeName
				} else {
					goType = "*" + param.TypeName
				}
			}
		} else {
			goType = tokenTypeToGoType(param.Type)
		}
		result += fmt.Sprintf("%s %s", param.Name, goType)
	}
	result += ")"

	if method.ReturnType != lx.IDENT {
		goType := tokenTypeToGoType(method.ReturnType)
		result += fmt.Sprintf(" %s", goType)
	}

	result += " {\n"
	if method.BodyCode != "" {
		result += "\t" + method.BodyCode + "\n"
	} else {
		result += "\t// Method body\n"
	}
	result += "}\n\n"

	return result
}

func tokenTypeToGoType(t lx.TokenType) string {
	switch t {
	case lx.TYPE_INT:
		return "int"
	case lx.TYPE_STRING:
		return "string"
	case lx.TYPE_BOOLEAN:
		return "bool"
	case lx.TYPE_FLOAT:
		return "float32"
	case lx.TYPE_DOUBLE:
		return "float64"
	case lx.TYPE_LONG:
		return "int64"
	case lx.TYPE_BYTE:
		return "byte"
	default:
		return "interface{}"
	}
}

func TranspileMethodStatement(lexer *lx.Lexer) string {
	tok := lexer.Tokenize()
	if tok.Type == lx.EOF || tok.Type == lx.RBRACE || tok.Type == lx.LBRACE {
		lexer.CheckPointThis(tok)
		return ""
	}

	switch tok.Type {
	case lx.TYPE_INT, lx.TYPE_STRING, lx.TYPE_BOOLEAN, lx.TYPE_FLOAT, lx.TYPE_DOUBLE, lx.TYPE_LONG, lx.TYPE_BYTE:
		leftTok, expr := GetVarAndExpr(lexer)
		goRhs := TranspileExpr(expr)

		switch tok.Type {
		case lx.TYPE_INT:
			return fmt.Sprintf("var %s int = %s", leftTok.Literal, goRhs)
		case lx.TYPE_STRING:
			return fmt.Sprintf("var %s string = %s", leftTok.Literal, goRhs)
		case lx.TYPE_BOOLEAN:
			return fmt.Sprintf("var %s bool = %s", leftTok.Literal, goRhs)
		case lx.TYPE_FLOAT:
			return fmt.Sprintf("var %s float32 = %s", leftTok.Literal, goRhs)
		case lx.TYPE_DOUBLE:
			return fmt.Sprintf("var %s float64 = %s", leftTok.Literal, goRhs)
		case lx.TYPE_LONG:
			return fmt.Sprintf("var %s int64 = %s", leftTok.Literal, goRhs)
		case lx.TYPE_BYTE:
			return fmt.Sprintf("var %s byte = %s", leftTok.Literal, goRhs)
		}
	case lx.IDENT:
		// Check if this is a return statement
		if tok.Literal == "return" {
			expr := ParseExpr(lexer)
			semi := lexer.Tokenize()
			if semi.Type != lx.SEMI {
				panic(fmt.Sprintf("[Error] Expected ';' after return statement at line %d", semi.Line))
			}
			return fmt.Sprintf("return %s", TranspileExpr(expr))
		}

		next := lexer.Tokenize()

		if next.Type == lx.IDENT {
			className := tok.Literal
			varName := next.Literal
			eq := lexer.Tokenize()
			if eq.Type != lx.ASSIGN {
				panic(fmt.Sprintf("[Error] Expected '=' after variable at line %d", eq.Line))
			}
			expr := ParseExpr(lexer)
			goRhs := TranspileExpr(expr)
			semi := lexer.Tokenize()
			if semi.Type != lx.SEMI {
				panic(fmt.Sprintf("[Error] Expected ';' after expression at line %d", semi.Line))
			}
			var typeStr string
			if _, isNewExpr := expr.(*NewExpr); isNewExpr {
				typeStr = "*" + className
			} else if strings.HasPrefix(goRhs, "&") || goRhs == "nil" {
				typeStr = "*" + className
			} else {
				typeStr = className
			}
			return fmt.Sprintf("var %s %s = %s", varName, typeStr, goRhs)
		} else if next.Type == lx.DOT {
			objExpr := &IdentifierExpr{Name: tok.Literal}
			currentExpr := Expr(objExpr)

			for {
				memberTok := lexer.Tokenize()
				if memberTok.Type != lx.IDENT {
					panic(fmt.Sprintf("[Error] Expected member name after '.' at line %d", memberTok.Line))
				}
				currentExpr = &MemberAccessExpr{
					Object: currentExpr,
					Member: memberTok.Literal,
				}

				peek := lexer.Tokenize()
				if peek.Type == lx.DOT {
					continue
				} else if peek.Type == lx.ASSIGN {
					expr := ParseExpr(lexer)
					goRhs := TranspileExpr(expr)
					semi := lexer.Tokenize()
					if semi.Type != lx.SEMI {
						panic(fmt.Sprintf("[Error] Expected ';' after assignment at line %d", semi.Line))
					}
					goLhs := TranspileExpr(currentExpr)
					return fmt.Sprintf("%s = %s", goLhs, goRhs)
				} else if peek.Type == lx.LPAREN {
					args := []Expr{}
					for {
						argTok := lexer.Tokenize()
						if argTok.Type == lx.RPAREN {
							break
						}
						lexer.CheckPointThis(argTok)
						arg := ParseExpr(lexer)
						args = append(args, arg)
						nextArg := lexer.Tokenize()
						if nextArg.Type == lx.RPAREN {
							break
						} else if nextArg.Type != lx.COMMA {
							panic(fmt.Sprintf("[Error] Expected ',' or ')' in argument list at line %d", nextArg.Line))
						}
					}
					callExpr := &CallExpr{
						Callee: currentExpr,
						Args:   args,
					}
					nextTok := lexer.Tokenize()
					if nextTok.Type != lx.SEMI {
						lexer.CheckPointThis(nextTok)
					}
					return TranspileExpr(callExpr)
				} else {
					panic(fmt.Sprintf("[Error] Expected '=', '(', or '.' after member access at line %d", peek.Line))
				}
			}
		} else if next.Type == lx.LPAREN {
			lexer.CheckPointThis(next)
			lexer.CheckPointThis(tok)
			expr := ParseExpr(lexer)
			semi := lexer.Tokenize()
			if semi.Type != lx.SEMI {
				panic(fmt.Sprintf("[Error] Expected ';' after expression at line %d", semi.Line))
			}
			return TranspileExpr(expr)
		} else {
			lexer.CheckPointThis(next)
			lexer.CheckPointThis(tok)
			expr := ParseExpr(lexer)
			semi := lexer.Tokenize()
			if semi.Type != lx.SEMI {
				panic(fmt.Sprintf("[Error] Expected ';' after expression at line %d", semi.Line))
			}
			return TranspileExpr(expr)
		}
	default:
		panic(fmt.Sprintf("[Error] Unexpected token '%s' (type: %s) in method body at line %d", tok.Literal, tok.Type, tok.Line))
	}

	return ""
}
