package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/smalldevshima/go-monkey/evaluator"
	"github.com/smalldevshima/go-monkey/lexer"
	"github.com/smalldevshima/go-monkey/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	writer := bufio.NewWriter(out)

	for {
		writer.WriteString(PROMPT)
		writer.Flush()
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(writer, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			writer.WriteString(evaluated.Inspect())
		} else {
			writer.WriteString("I cannot evaluate that, sorry!")
		}
		writer.WriteString("\n")
	}
}

func printParserErrors(out *bufio.Writer, errors []string) {
	out.WriteString(fmt.Sprintf("parser has %d errors:\n", len(errors)))
	for i, msg := range errors {
		if i >= 10 {
			out.WriteString("(omitting more errors)\n")
			break
		}
		out.WriteString(fmt.Sprintf("%3d: %s\n", i+1, msg))
	}
}
