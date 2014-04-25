package main

import (
	"fmt"
)

const EOF = -1

type scanner struct {
	s   []rune
	pos int
}

func newScanner(s string) scanner {
	return scanner{[]rune(s), 0}
}

func (sc *scanner) peek() rune {
	if sc.pos == len(sc.s) {
		return EOF
	}
	return sc.s[sc.pos]
}

func (sc *scanner) match(ch rune) {
	if sc.pos == len(sc.s) {
		if ch != EOF {
			sc.require(fmt.Sprintf("'%c'", ch))
		}
	} else if sc.s[sc.pos] != ch {
		if ch != EOF {
			sc.require(fmt.Sprintf("'%c'", ch))
		} else {
			sc.require("(输入结束)")
		}
	} else {
		sc.pos++
	}
}

func (sc *scanner) require(s string) {
	if sc.pos == len(sc.s) {
		panic(fmt.Sprintf("scanner: 输入过早结束。还需要%s。", s))
	} else {
		panic(fmt.Sprintf("scanner: 不应有'%c'。应为%s。", sc.s[sc.pos], s))
	}
}

func (sc *scanner) matchInt() int {
	n := 0
	for ch := sc.peek(); '0' <= ch && ch <= '9'; ch = sc.peek() {
		n = n*10 + int(ch-'0')
		sc.pos++ // match
	}
	return n
}

func (sc *scanner) matchElem() string {
	letter := sc.peek()
	name := []rune{letter}
	sc.pos++
	for letter := sc.peek(); 'a' <= letter && letter <= 'z'; letter = sc.peek() {
		name = append(name, letter)
		sc.pos++
	}
	return string(name)
}
