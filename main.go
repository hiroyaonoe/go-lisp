package main

import (
	"bufio"
	"errors"
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
	e := eval.NewEnv()
	for {
		fmt.Print("> ")
		if !s.Scan() {
			break
		}
		tokens := l.ReadString(s.Text())
		ast, err := p.Parse(tokens)
		if err != nil {
			if errors.Is(err, parser.ErrNeedNextTokens) {
				continue
			}
			return err
		}
		value, err := e.Eval(ast)
		if err != nil {
			return err
		}
		fmt.Println(value)
	}
	return nil
}

func main() {
	err := readEvalPrint()
	fmt.Println(err)
}
