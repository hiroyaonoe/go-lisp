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
	e := eval.NewEnv(nil)
	new := true
	for {
		if new {
			fmt.Print(" > ")
		} else {
			fmt.Print("*> ")
		}
		if !s.Scan() {
			break
		}
		tokens, err := l.ReadString(s.Text())
		if err != nil {
			if errors.Is(err, lexer.EOF) {
				new = false
				continue
			}
			fmt.Println(err)
			new = true
			continue
		}
		ast, err := p.Parse(tokens)
		if err != nil {
			if errors.Is(err, parser.EOF) {
				new = false
				continue
			}
			fmt.Println(err)
			new = true
			continue
		}
		for _, node := range ast {
			ret, err := e.Eval(node)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("val:", ret)
		}
		new = true
	}
	return nil
}

func main() {
	err := readEvalPrint()
	fmt.Println(err)
}
