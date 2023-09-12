package main

import (
	"fmt"
	"os/user"

	"github.com/yassinebenaid/nishimia/repl"
)

func main() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome %s , this is nishimia lang ready to get you excited ! \n", user.Username)

	repl.Start()

	fmt.Print("\nGood by !\n")
}
