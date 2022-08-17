package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/hiroyaonoe/go-lisp/eval"
	"github.com/hiroyaonoe/go-lisp/lexer"
	"github.com/hiroyaonoe/go-lisp/parser"
)

func readEvalPrint() error {
	s := bufio.NewScanner(os.Stdin)
	l := lexer.NewLexer()
	p := parser.NewParser()
	e := eval.NewEnv(nil)
	for {
		fmt.Print("> ")
		if !s.Scan() {
			break
		}
		tokens, err := l.ReadString(s.Text())
		if err != nil {
			fmt.Println(err)
			continue
		}
		ast, err := p.Parse(tokens)
		fmt.Println(ast)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, node := range ast {
			ret, err := e.Eval(node)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(ret)
		}
	}
	return nil
}

func main() {
	err := readEvalPrint()
	fmt.Println(err)
}
