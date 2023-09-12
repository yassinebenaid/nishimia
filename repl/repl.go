package repl

import (
	"bufio"
	"fmt"
	"os"

	"github.com/yassinebenaid/nishimia/lexer"
	"github.com/yassinebenaid/nishimia/token"
)

const PROMPT = ">>> "

func Start() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print(PROMPT)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "exit" {
			break
		}
		lex := lexer.New(line)

		for tok := lex.NextToken(); tok.Type != token.EOF; tok = lex.NextToken() {
			fmt.Printf("%+s\n", tok)
		}

		fmt.Print(PROMPT)
	}
}
