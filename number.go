package main

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func lcm(a, b int) int {
	return a / gcd(a, b) * b
}

func abs(a int) int {
	if a >= 0 {
		return a
	}
	return -a
}
