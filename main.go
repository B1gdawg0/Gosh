package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	lx "github.com/B1gdawg0/Gosh/src/lexer"
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

func transpile(in string) string{
	lexer := lx.NewLexer(in)
	var out strings.Builder
	userImports := []string{}

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
		if tok.Type == lx.EOF{
			break
		}

		switch tok.Type {
			case lx.TYPE_INT:
				left := lexer.Tokenize()
				lexer.Tokenize()
				right := lexer.Tokenize()
				lexer.Tokenize()
				out.WriteString(fmt.Sprintf("\tvar %s int = %s\n", left.Literal, right.Literal)) 
			case lx.NATIVE:
				out.WriteString("\t" + strings.TrimSpace(tok.Literal) + "\n")
		}
	}
	out.WriteString("}\n")
	return out.String()
}