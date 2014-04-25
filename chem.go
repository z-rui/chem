package main

import (
	"fmt"
	"strings"
)

var (
	uniSub    = [...]rune{'₀', '₁', '₂', '₃', '₄', '₅', '₆', '₇', '₈', '₉'}
	uniSup    = [...]rune{'⁰', '¹', '²', '³', '⁴', '⁵', '⁶', '⁷', '⁸', '⁹'}
	uniCharge = map[rune]rune{'+': '⁺', '-': '⁻'}
)

func chemMat(lhs, rhs []term, elems []string) mat {
	m := make(mat, 0, len(elems))
	n := len(lhs) + len(rhs)
	for _, elem := range elems {
		row := make(matRow, 0, 2*n)
		for _, term := range lhs {
			row = append(row, term.molecule[elem])
		}
		for _, term := range rhs {
			row = append(row, -term.molecule[elem])
		}
		m = append(m, row)
	}
	return m
}

func chemTerm(coeff int, term string) string {
	if coeff == 1 {
		return term
	}
	return fmt.Sprintf("%d%s", coeff, term)
}

func chemEqnSide(coeff []int, term []term) string {
	terms := make([]string, len(term))
	for i := range terms {
		name := term[i].name
		if useUnicodeSubscript {
			name = unicodePrep(name)
		}
		terms[i] = chemTerm(coeff[i], name)
	}
	return strings.Join(terms, " + ")
}

func chemEqn(coeff []int, lhs, rhs []term) string {
	mid := len(lhs)
	l := chemEqnSide(coeff[:mid], lhs)
	r := chemEqnSide(coeff[mid:], rhs)
	return l + " = " + r
}

// unicodePrep把ASCII形式的化学式转换为含有合理Unicode字符的形式。例如
// KAl(SO4)2.12H2O -> KAl(SO₄)₂·12H₂O
// SO4[2-] -> SO₄²⁻
// e[-] -> e⁻
func unicodePrep(name string) string {
	isCoeff := true
	isSup := false
	return strings.Map(func(r rune) rune {
		switch {
		case isDigit(r):
			if !isCoeff {
				if isSup {
					return uniSup[r-'0']
				} else {
					return uniSub[r-'0']
				}
			} else {
				return r
			}
		case r == '+', r == '-':
			return uniCharge[r]
		case r == '.':
			isCoeff = true
			return '·'
		case r == '[':
			isSup = true
			if isCoeff {
				isCoeff = false
				return 'e'
			}
			return -1
		case r == ']':
			isSup = false
			return -1
		default:
			isCoeff = false
			return r
		}
	}, name)
}
