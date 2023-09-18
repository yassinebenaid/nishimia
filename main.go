package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/yassinebenaid/nishimia/eval"
	"github.com/yassinebenaid/nishimia/lexer"
	"github.com/yassinebenaid/nishimia/object"
	"github.com/yassinebenaid/nishimia/parser"
	"github.com/yassinebenaid/nishimia/repl"
)

func main() {
	if len(os.Args) < 2 {
		Interactive()
		os.Exit(0)
	}

	if len(os.Args) == 2 {
		RunFile(os.Args[1])
		os.Exit(0)
	}

	fmt.Println("undefined options", os.Args)
}

func Interactive() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome %s , this is nishimia lang ready to get you excited ! \n", user.Username)

	repl.Start(os.Stdin, os.Stdout)

	fmt.Print("\nGood by !\n")
}

func RunFile(fn string) {
	file, err := os.ReadFile(fn)
	if err != nil {
		fmt.Println("Error: \n\t", err)
	}

	env := object.NewEnvirement()
	lex := lexer.New(string(file))
	par := parser.New(lex)
	program := par.ParseProgram()
	result := eval.Eval(program, env)

	if result.Type() != object.NULL_OBJ {
		fmt.Println(result.Inspect())
	}
}
