package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	useUnicodeSubscript bool
	interactive         bool
)

func main() {
	flag.BoolVar(&useUnicodeSubscript, "u", false, "use Unicode subscript characters in output")
	flag.BoolVar(&interactive, "i", true, "interactive mode")
	flag.Parse()
	r := bufio.NewReader(os.Stdin)
Loop:
	for {
		if interactive {
			os.Stdout.WriteString("> ")
		}
		s, _, err := r.ReadLine()
		switch err {
		case nil:
			// nothing to do
		case io.EOF:
			break Loop
		default:
			panic(err)
		}
		lhs, rhs, elems, err := newParser(string(s)).parse()
		if err != nil {
			fmt.Println(err)
		} else {
			//fmt.Println(lhs, rhs, elems)
			m := chemMat(lhs, rhs, elems)
			//fmt.Println(m)
			space := m.solveHomo()
			for _, coeff := range space {
				//fmt.Println(coeff)
				fmt.Println(chemEqn(coeff, lhs, rhs))
			}
		}
	}
}
