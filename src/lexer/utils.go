package lexer

func GetLeftRightFromDefined(lx *Lexer)(*Token, *Token){
	left := lx.Tokenize()
	lx.Tokenize()
	right := lx.Tokenize()
	lx.Tokenize()
	return &left, &right
}