package main

import (
	_ "fmt"
)

type matRow []int
type mat []matRow

func (m mat) rowSwap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m mat) rowReduce(i int) {
	vecGcd(m[i])
}

func vecGcd(r []int) {
	g := 0
	for _, v := range r {
		if v := abs(v); v != 0 {
			if g == 0 {
				g = v
			} else {
				g = gcd(g, v)
			}
		}
	}
	if g != 0 {
		// divRow
		for i := range r {
			r[i] /= g
		}
	}
}

func (m mat) rowElim(i, j, k int) {
	a := abs(m[i][k])
	b := abs(m[j][k])
	l := lcm(a, b)
	//println("elim 1", m[i], m[j])
	m.mRow(j, l/m[j][k])
	m.mRowAdd(i, j, -l/m[i][k])
	//println("elim 2", m[i], m[j])
}

func (m mat) mRow(i, n int) {
	r := m[i]
	for i := range r {
		r[i] *= n
	}
}

func (m mat) mRowAdd(i, j, n int) {
	src := m[i]
	dst := m[j]
	for i := range dst {
		dst[i] += src[i] * n
	}
}

func (m mat) rref() (rank int) {
	n := len(m[0])
	offs := 0
	for i := range m {
		//fmt.Println(m)
		m.rowReduce(i)
		k := i + offs
		for k < n && m[i][k] == 0 {
			if m[i][k] == 0 {
				for j := i + 1; j < len(m); j++ {
					if m[j][k] != 0 {
						m.rowSwap(i, j)
						break
					}
				}
			}
			if m[i][k] == 0 {
				offs++
				k++
			}
		}
		if k == n {
			break
		}
		rank++
		for j := range m {
			if i != j && m[j][k] != 0 {
				m.rowElim(i, j, k)
			}
		}
	}
	return
}

func (m mat) solveHomo() [][]int {
	rank := m.rref()
	n := len(m[0])
	dim := n - rank
	m1 := make(mat, 0, n)
	// augment
	aug := make([]int, dim)
	for i := 0; i < rank; i++ {
		m1 = append(m1, append(m[i], aug...))
	}
	isindep := make([]bool, n)
	indep := 0
	for i := 0; i < rank; i++ {
		k := i + indep
		for m[i][k] == 0 {
			//println(k, "-> indep")
			isindep[k] = true
			indep++
			k++
		}
	}
	// still need (n - indep) indep vars
	for i := n - 1; indep < dim && i >= 0; i-- {
		if !isindep[i] {
			//println(i, "-> indep")
			isindep[i] = true
			indep++
		}
	}
	//fmt.Println(isindep)
	if indep != dim {
		panic("indep != dim")
	}
	for i, v := range isindep {
		if v {
			row := make(matRow, n+dim)
			row[i] = 1
			row[n+dim-indep] = 1
			indep--
			m1 = append(m1, row)
		}
	}
	//fmt.Println(m1)
	if len(m1) != n {
		panic("len(m1) != n")
	}
	rank = m1.rref()
	if rank != n {
		panic("rank != n")
	}
	for i := range m1 {
		m1.rowReduce(i)
	}
	m1.diagReduce()
	//fmt.Println(rank, m1)
	space := make([][]int, dim)
	for i := 0; i < dim; i++ {
		xi := make([]int, n)
		for j := 0; j < n; j++ {
			xi[j] = m1[j][n+i]
		}
		vecGcd(xi)
		space[i] = xi
	}
	//fmt.Println(space)
	return space
}

func (m mat) diagReduce() {
	l := 0
	for i := range m {
		m.rowReduce(i)
		if l == 0 {
			l = abs(m[i][i])
		} else {
			l = lcm(l, abs(m[i][i]))
		}
	}
	for i := range m {
		if m[i][i] != 0 {
			m.mRow(i, l/m[i][i])
		}
	}
}
