package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	lx "github.com/B1gdawg0/Gosh/src/lexer"
	"github.com/B1gdawg0/Gosh/src/parsing"
)

func main(){
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

	out := transpile(string(bytes))
	if(debug){
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
    lexer := lx.NewLexer(in)
    var out strings.Builder
    userImports := []string{}
    needsBigInt := false

    tok := lexer.Tokenize()
    if tok.Type == lx.IMPORT {
        for {
            tok = lexer.Tokenize()
            if tok.Type == lx.SEMI || tok.Type == lx.EOF {
                break
            }
            userImports = append(userImports, tok.Literal)
        }
    } else {
        lexer.CheckPointThis(tok)
    }

    out.WriteString("package main\n")
    out.WriteString("import (\n")
    for _, im := range userImports {
        out.WriteString(fmt.Sprintf("\t\"%s\"\n", im))
    }
    out.WriteString(")\n")
    out.WriteString("func main() {\n")

    for {
        tok := lexer.Tokenize()
        if tok.Type == lx.EOF {
            break
        }
	
        switch tok.Type {
        case lx.TYPE_INT, lx.TYPE_LONG, lx.TYPE_FLOAT, lx.TYPE_DOUBLE,
             lx.TYPE_BYTE, lx.TYPE_STRING, lx.TYPE_BOOLEAN:

            leftTok, expr := parsing.GetVarAndExpr(lexer)
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
        case lx.NATIVE:
            out.WriteString("\t" + strings.TrimSpace(tok.Literal) + "\n")
        }
    }

    out.WriteString("}\n")
    result := out.String()
    if needsBigInt {
        result = strings.Replace(result, "import (\n", "import (\n\t\"math/big\"\n", 1)
    }

    return result
}
