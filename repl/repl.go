package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/yassinebenaid/nishimia/eval"
	"github.com/yassinebenaid/nishimia/lexer"
	"github.com/yassinebenaid/nishimia/parser"
)

const PROMPT = ">>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	fmt.Print(PROMPT)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "exit" {
			break
		}
		lex := lexer.New(line)
		par := parser.New(lex)
		program := par.ParseProgram()

		if errs := par.Errors(); len(errs) > 0 {
			io.WriteString(out, "Parsing failed : \n")

			for i, err := range errs {
				io.WriteString(out, fmt.Sprintf("\t#%d %s\n", i, err))
			}

			io.WriteString(out, PROMPT)
			continue
		}

		evaluated := eval.Eval(program)

		if evaluated != nil {
			io.WriteString(out, "\n")
			io.WriteString(out, evaluated.Inspect())
		}

		io.WriteString(out, "\n")

		fmt.Print(PROMPT)
	}
}
