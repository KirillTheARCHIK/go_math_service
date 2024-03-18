package mymath

func Fact(n int) int {
	if n < 0 {
		panic("factorial of n<0")
	}
	if n <= 1 {
		return 1
	}
	return Fact(n-1) * n
}
