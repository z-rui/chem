package main

import (
	"fmt"
)

type parser struct {
	scanner
	elems map[string]bool
}

type term struct {
	name     string
	molecule map[string]int
}

type parseError struct {
	pos  int
	what string
}

func (e parseError) Error() string {
	return fmt.Sprintf("%d:%s", e.pos, e.what)
}

func newParser(s string) *parser {
	return &parser{newScanner(s), make(map[string]bool)}
}

func (p *parser) parse() (lhs, rhs []term, elems []string, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = parseError{p.pos, e.(string)}
		}
	}()
	lhs = p.poly()
	p.match('=')
	rhs = p.poly()
	p.match(EOF)
	elems = make([]string, 0, len(p.elems))
	for k := range p.elems {
		elems = append(elems, k)
	}
	return
}

func (p *parser) poly() []term {
	terms := []term{p.term()}
	for ch := p.peek(); ch == '+'; ch = p.peek() {
		p.match('+')
		terms = append(terms, p.term())
	}
	if len(terms) == 0 {
		p.require("(化学式)")
	}
	return terms
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func (p *parser) term() term {
	var m map[string]int
	start := p.pos
	switch ch := p.peek(); ch {
	case 'e':
		p.match('e')
		fallthrough
	case '[':
		m = make(map[string]int)
		break
	default:
		m = p.molecule()
		for p.peek() == '.' {
			p.match('.')
			n := 1
			if isDigit(p.peek()) {
				n = p.matchInt()
			}
			if n == 0 {
				panic("系数为0")
			}
			m1 := p.molecule()
			for k, v := range m1 {
				m[k] += v * n
			}
		}
	}
	if p.peek() == '[' {
		p.match('[')
		// 电荷数
		n := 1
		if isDigit(p.peek()) {
			n = p.matchInt()
		}
		if n == 0 {
			panic("电荷数为0")
		}
		charge := p.peek()
		switch charge {
		case '+':
			p.match('+')
		case '-':
			p.match('-')
			n = -n
		default:
			p.require("'+'或'-'")
		}
		m["+"] = n
		p.elems["+"] = true
		p.match(']')
	}
	stop := p.pos
	return term{p.name(start, stop), m}
}

func (p *parser) name(start, end int) string {
	return string(p.s[start:end])
}

func (p *parser) molecule() map[string]int {
	m := make(map[string]int)
Loop:
	for {
		ch := p.peek()
		switch {
		case ch == '(':
			p.match('(')
			m1 := p.molecule()
			p.match(')')
			n := p.matchInt()
			if n == 0 {
				panic("parser: 缺少下标或下标为0")
			}
			for k, v := range m1 {
				m[k] += v * n
			}
		case 'A' <= ch && ch <= 'Z':
			e := p.matchElem()
			n := 1
			if ch := p.peek(); isDigit(ch) {
				n = p.matchInt()
				if n == 0 {
					panic("parser: 下标为0")
				}
			}
			m[e] += n
		default:
			break Loop
		}
	}
	if len(m) == 0 {
		p.require("(元素)")
	}
	return m
}

func (p *parser) matchElem() string {
	elem := p.scanner.matchElem()
	p.elems[elem] = true
	return elem
}
