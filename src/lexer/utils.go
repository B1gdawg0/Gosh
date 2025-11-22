package lexer

import "fmt"

func RemoveNumericSuffix(literal string) string{
	if literal[len(literal)-1] == 'f' || literal[len(literal)-1] == 'F' ||
       literal[len(literal)-1] == 'd' || literal[len(literal)-1] == 'D' ||
	   literal[len(literal)-1] == 'L' || literal[len(literal)-1] == 'l' {
		fmt.Println(literal[:len(literal)-1])
        return literal[:len(literal)-1]
    }
	fmt.Println(literal)
	return literal
}