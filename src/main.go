package main

import (
	"fmt"
	"os"
)

func main(){
	bytes, err := os.ReadFile("examples/00.gosh");
	if err != nil{
		fmt.Println("error: ", err.Error());
		return ;
	}
	source := string(bytes);

	fmt.Println(source);
}