package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	lx "github.com/B1gdawg0/Gosh/src/lexer"
	"github.com/B1gdawg0/Gosh/src/parsing"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: gosh <file.gosh> [--debug=true]")
		os.Exit(1)
	}

	filePath := os.Args[1]

	debug := false
	if len(os.Args) >= 3 {
		if os.Args[2] == "--debug=true" {
			debug = true
		}
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("[Error]: can't find source file")
		os.Exit(1)
	}

	out := transpileWithPath(string(bytes), filePath)
	if debug {
		fmt.Println(out)
		os.Exit(0)
	}

	tmpFile, err := os.CreateTemp("", "gosh_*.go")
	if err != nil {
		fmt.Printf("[Error]: Failed to create pre-process go file: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(out); err != nil {
		fmt.Printf("[Error]: Failed to write string into pre-process go file: %v\n", err)
		os.Exit(1)
	}
	tmpFile.Close()

	cmd := exec.Command("go", "run", tmpFile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("[Error] Run time error: %v\n", err)
		os.Exit(1)
	}
}

func transpile(in string) string {
	return transpileWithPath(in, "")
}

func transpileWithPath(in string, basePath string) string {
	lexer := lx.NewLexer(in)
	var out strings.Builder
	userImports := []string{}
	goshImports := []string{}

	for {
		tok := lexer.Tokenize()
		if tok.Type == lx.NATIVE_IMPORT {
			var importParts []string
			for {
				tok = lexer.Tokenize()
				if tok.Type == lx.SEMI || tok.Type == lx.EOF {
					break
				}
				if tok.Type == lx.IDENT {
					importParts = append(importParts, tok.Literal)
				} else if tok.Type == lx.DIV && len(importParts) > 0 {
					importParts = append(importParts, "/")
					nextTok := lexer.Tokenize()
					if nextTok.Type == lx.IDENT {
						importParts = append(importParts, nextTok.Literal)
					} else {
						lexer.CheckPointThis(nextTok)
						break
					}
				} else {
					lexer.CheckPointThis(tok)
					break
				}
			}
			if len(importParts) > 0 {
				importName := strings.Join(importParts, "")
				userImports = append(userImports, importName)
			}
		} else if tok.Type == lx.IMPORT {
			for {
				tok = lexer.Tokenize()
				if tok.Type == lx.SEMI || tok.Type == lx.EOF {
					break
				}
				importName := tok.Literal
				goshImports = append(goshImports, importName)
			}
		} else {
			lexer.CheckPointThis(tok)
			break
		}
	}

	importedClasses := []*parsing.ClassDecl{}
	for _, goshImport := range goshImports {
		var goshFilePath string
		var fileName string
		if strings.HasSuffix(goshImport, ".gosh") {
			goshFilePath = goshImport
			fileName = goshImport
		} else {
			fileName = goshImport + ".gosh"
			goshFilePath = fileName
		}

		if basePath != "" {
			baseDir := basePath
			if lastSlash := strings.LastIndex(basePath, "/"); lastSlash != -1 {
				baseDir = basePath[:lastSlash+1]
			} else if lastSlash := strings.LastIndex(basePath, "\\"); lastSlash != -1 {
				baseDir = basePath[:lastSlash+1]
			} else {
				baseDir = "./"
			}
			goshFilePath = baseDir + goshFilePath
		}

		found := false
		searchDirs := []string{}

		if basePath != "" {
			baseDir := basePath
			if lastSlash := strings.LastIndex(basePath, "/"); lastSlash != -1 {
				baseDir = basePath[:lastSlash+1]
			} else if lastSlash := strings.LastIndex(basePath, "\\"); lastSlash != -1 {
				baseDir = basePath[:lastSlash+1]
			} else {
				baseDir = "./"
			}
			searchDirs = append(searchDirs, baseDir)
		}

		searchDirs = append(searchDirs,
			"src/examples/class/",
			"src/examples/",
			"./",
		)

		for _, dir := range searchDirs {
			searchPath := dir + fileName
			if _, err := os.Stat(searchPath); err == nil {
				goshFilePath = searchPath
				found = true
				break
			}
		}

		if !found {
			for _, dir := range searchDirs {
				entries, err := os.ReadDir(dir)
				if err != nil {
					continue
				}
				for _, entry := range entries {
					if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".gosh") {
						entryName := strings.ToLower(strings.TrimSuffix(entry.Name(), ".gosh"))
						importNameLower := strings.ToLower(goshImport)
						if strings.Contains(entryName, importNameLower) || strings.Contains(importNameLower, entryName) {
							goshFilePath = dir + entry.Name()
							found = true
							break
						}
					}
				}
				if found {
					break
				}
			}
		}

		if !found {
			fmt.Fprintf(os.Stderr, "[Warning] Could not find Gosh file: %s (skipping import)\n", goshImport)
			continue
		}

		importBytes, err := os.ReadFile(goshFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[Warning] Could not read Gosh file: %s (skipping import)\n", goshFilePath)
			continue
		}

		importLexer := lx.NewLexer(string(importBytes))
		for {
			importTok := importLexer.Tokenize()
			if importTok.Type == lx.IMPORT || importTok.Type == lx.NATIVE_IMPORT {
				for {
					importTok = importLexer.Tokenize()
					if importTok.Type == lx.SEMI || importTok.Type == lx.EOF {
						break
					}
				}
			} else {
				importLexer.CheckPointThis(importTok)
				break
			}
		}

		for {
			classTok := importLexer.Tokenize()
			if classTok.Type == lx.EOF {
				break
			}
			if classTok.Type == lx.CLASS {
				importedClass := parsing.ParseClass(importLexer)
				importedClasses = append(importedClasses, importedClass)
			} else {
				break
			}
		}
	}

	out.WriteString("package main\n")
	out.WriteString("import (\n")
	for _, im := range userImports {
		out.WriteString(fmt.Sprintf("\t\"%s\"\n", im))
	}
	out.WriteString(")\n")

	allClasses := append(importedClasses, []*parsing.ClassDecl{}...)

	classMap := make(map[string]*parsing.ClassDecl)
	for _, class := range importedClasses {
		classMap[class.Name] = class
	}

	parsing.SetClassRegistry(classMap)

	classes := []*parsing.ClassDecl{}
	for {
		tok := lexer.Tokenize()
		if tok.Type == lx.EOF {
			break
		}
		if tok.Type == lx.CLASS {
			class := parsing.ParseClass(lexer)
			classes = append(classes, class)
			classMap[class.Name] = class
			parsing.SetClassRegistry(classMap)
		} else {
			lexer.CheckPointThis(tok)
			break
		}
	}

	allClasses = append(allClasses, classes...)

	for _, class := range allClasses {
		out.WriteString(parsing.TranspileClass(class))
	}

	out.WriteString("func main() {\n")

	for {
		tok := lexer.Tokenize()
		if tok.Type == lx.EOF {
			break
		}

		switch tok.Type {
		case lx.TYPE_INT, lx.TYPE_LONG, lx.TYPE_FLOAT, lx.TYPE_DOUBLE,
			lx.TYPE_BYTE, lx.TYPE_STRING, lx.TYPE_BOOLEAN:
			leftTok, isArray, expr := parsing.GetVarAndExpr(lexer)

			if isArray {
				goType := fmt.Sprintf("[]%s", primitiveGoType(tok.Type))
				goRhs := parsing.TranspileExprWithType(expr, goType)
				out.WriteString(fmt.Sprintf("\tvar %s %s = %s\n", leftTok.Literal, goType, goRhs))
				continue
			}

			goRhs := parsing.TranspileExpr(expr)

			switch tok.Type {
			case lx.TYPE_INT:
				out.WriteString(fmt.Sprintf("\tvar %s int = %s\n", leftTok.Literal, goRhs))
			case lx.TYPE_LONG:
				out.WriteString(fmt.Sprintf("\tvar %s int64 = %s\n", leftTok.Literal, lx.RemoveNumericSuffix(goRhs)))
			case lx.TYPE_FLOAT:
				out.WriteString(fmt.Sprintf("\tvar %s float32 = %s\n", leftTok.Literal, lx.RemoveNumericSuffix(goRhs)))
			case lx.TYPE_DOUBLE:
				out.WriteString(fmt.Sprintf("\tvar %s float64 = %s\n", leftTok.Literal, lx.RemoveNumericSuffix(goRhs)))
			case lx.TYPE_BYTE:
				out.WriteString(fmt.Sprintf("\tvar %s byte = %s\n", leftTok.Literal, goRhs))
			case lx.TYPE_STRING:
				out.WriteString(fmt.Sprintf("\tvar %s string = %s\n", leftTok.Literal, goRhs))
			case lx.TYPE_BOOLEAN:
				out.WriteString(fmt.Sprintf("\tvar %s bool = %s\n", leftTok.Literal, goRhs))
			}
		case lx.IF:
			ifExpr := parsing.ParseIfExpr(lexer)
			goCode := parsing.TranspileExpr(ifExpr)
			out.WriteString(goCode + "\n")
		case lx.IDENT:
			next := lexer.Tokenize()
			var pendingIndexTok *lx.Token

			if next.Type == lx.LBRACKET {
				bracketClose := lexer.Tokenize()
				if bracketClose.Type == lx.RBRACKET {
					className := tok.Literal
					varName := lexer.Tokenize()
					if varName.Type != lx.IDENT {
						panic(fmt.Sprintf("[Error] Expected variable name after '[]' at line %d", varName.Line))
					}
					eq := lexer.Tokenize()
					if eq.Type != lx.ASSIGN {
						panic(fmt.Sprintf("[Error] Expected '=' after variable at line %d", eq.Line))
					}
					expr := parsing.ParseExpr(lexer)
					typeStr := "[]*" + className
					goRhs := parsing.TranspileExprWithType(expr, typeStr)
					semi := lexer.Tokenize()
					if semi.Type != lx.SEMI {
						panic(fmt.Sprintf("[Error] Expected ';' after expression at line %d", semi.Line))
					}
					out.WriteString(fmt.Sprintf("\tvar %s %s = %s\n", varName.Literal, typeStr, goRhs))
					continue
				}
				tokenCopy := bracketClose
				pendingIndexTok = &tokenCopy
			}

			if isCompoundAssignToken(next.Type) {
				expr := parsing.ParseExpr(lexer)
				goRhs := parsing.TranspileExpr(expr)
				semi := lexer.Tokenize()
				if semi.Type != lx.SEMI {
					panic(fmt.Sprintf("[Error] Expected ';' after expression at line %d", semi.Line))
				}
				out.WriteString(fmt.Sprintf("\t%s %s %s\n", tok.Literal, next.Literal, goRhs))
				continue
			} else if next.Type == lx.IDENT {
				className := tok.Literal
				varName := next.Literal
				eq := lexer.Tokenize()
				if eq.Type != lx.ASSIGN {
					panic(fmt.Sprintf("[Error] Expected '=' after variable at line %d", eq.Line))
				}
				expr := parsing.ParseExpr(lexer)
				goRhs := parsing.TranspileExpr(expr)
				semi := lexer.Tokenize()
				if semi.Type != lx.SEMI {
					panic(fmt.Sprintf("[Error] Expected ';' after expression at line %d", semi.Line))
				}
				var typeStr string
				if _, isNewExpr := expr.(*parsing.NewExpr); isNewExpr {
					typeStr = "*" + className
				} else if strings.HasPrefix(goRhs, "&") || goRhs == "nil" {
					typeStr = "*" + className
				} else {
					typeStr = className
				}
				out.WriteString(fmt.Sprintf("\tvar %s %s = %s\n", varName, typeStr, goRhs))
				continue
			} else if next.Type == lx.DOT || next.Type == lx.LBRACKET {
				objExpr := &parsing.IdentifierExpr{Name: tok.Literal}
				currentExpr := parsing.Expr(objExpr)

				if pendingIndexTok != nil {
					lexer.CheckPointThis(*pendingIndexTok)
					pendingIndexTok = nil
				}

				peek := next

			memberLoop:
				for {
					switch peek.Type {
					case lx.DOT:
						memberTok := lexer.Tokenize()
						if memberTok.Type != lx.IDENT {
							panic(fmt.Sprintf("[Error] Expected member name after '.' at line %d", memberTok.Line))
						}
						currentExpr = &parsing.MemberAccessExpr{
							Object: currentExpr,
							Member: memberTok.Literal,
						}
						peek = lexer.Tokenize()
					case lx.LBRACKET:
						indexExpr := parsing.ParseExpr(lexer)
						closeTok := lexer.Tokenize()
						if closeTok.Type != lx.RBRACKET {
							panic(fmt.Sprintf("[Error] Expected ']' after '[' at line %d", closeTok.Line))
						}
						currentExpr = &parsing.IndexExpr{
							Collection: currentExpr,
							Index:      indexExpr,
						}
						peek = lexer.Tokenize()
					case lx.ASSIGN:
						expr := parsing.ParseExpr(lexer)
						goRhs := parsing.TranspileExpr(expr)
						semi := lexer.Tokenize()
						if semi.Type != lx.SEMI {
							panic(fmt.Sprintf("[Error] Expected ';' after assignment at line %d", semi.Line))
						}
						goLhs := parsing.TranspileExpr(currentExpr)
						out.WriteString(fmt.Sprintf("\t%s = %s\n", goLhs, goRhs))
						break memberLoop
					case lx.PLUS_ASSIGN, lx.MINUS_ASSIGN, lx.MULT_ASSIGN, lx.DIV_ASSIGN:
						expr := parsing.ParseExpr(lexer)
						goRhs := parsing.TranspileExpr(expr)
						semi := lexer.Tokenize()
						if semi.Type != lx.SEMI {
							panic(fmt.Sprintf("[Error] Expected ';' after assignment at line %d", semi.Line))
						}
						goLhs := parsing.TranspileExpr(currentExpr)
						out.WriteString(fmt.Sprintf("\t%s %s %s\n", goLhs, peek.Literal, goRhs))
						break memberLoop
					case lx.LPAREN:
						args := []parsing.Expr{}
						for {
							argTok := lexer.Tokenize()
							if argTok.Type == lx.RPAREN {
								break
							}
							lexer.CheckPointThis(argTok)
							arg := parsing.ParseExpr(lexer)
							args = append(args, arg)
							nextArg := lexer.Tokenize()
							if nextArg.Type == lx.RPAREN {
								break
							} else if nextArg.Type != lx.COMMA {
								panic(fmt.Sprintf("[Error] Expected ',' or ')' in argument list at line %d", nextArg.Line))
							}
						}
						callExpr := &parsing.CallExpr{
							Callee: currentExpr,
							Args:   args,
						}
						nextTok := lexer.Tokenize()
						if nextTok.Type != lx.SEMI {
							lexer.CheckPointThis(nextTok)
						}
						goExpr := parsing.TranspileExpr(callExpr)
						out.WriteString(fmt.Sprintf("\t%s\n", goExpr))
						break memberLoop
					default:
						panic(fmt.Sprintf("[Error] Expected '=', '(', '.', or '[' after member access at line %d", peek.Line))
					}
				}
				continue
			} else {
				lexer.CheckPointThis(next)
				lexer.CheckPointThis(tok)
				expr := parsing.ParseExpr(lexer)
				semi := lexer.Tokenize()
				if semi.Type != lx.SEMI {
					panic(fmt.Sprintf("[Error] Expected ';' after expression at line %d", semi.Line))
				}
				goExpr := parsing.TranspileExpr(expr)
				out.WriteString(fmt.Sprintf("\t%s\n", goExpr))
			}
		case lx.NATIVE:
			out.WriteString("\t" + strings.TrimSpace(tok.Literal) + "\n")
		}
	}

	out.WriteString("}\n")
	return out.String()
}

func isCompoundAssignToken(t lx.TokenType) bool {
	switch t {
	case lx.PLUS_ASSIGN, lx.MINUS_ASSIGN, lx.MULT_ASSIGN, lx.DIV_ASSIGN:
		return true
	default:
		return false
	}
}

func primitiveGoType(t lx.TokenType) string {
	switch t {
	case lx.TYPE_INT:
		return "int"
	case lx.TYPE_LONG:
		return "int64"
	case lx.TYPE_FLOAT:
		return "float32"
	case lx.TYPE_DOUBLE:
		return "float64"
	case lx.TYPE_BYTE:
		return "byte"
	case lx.TYPE_STRING:
		return "string"
	case lx.TYPE_BOOLEAN:
		return "bool"
	default:
		return "interface{}"
	}
}
