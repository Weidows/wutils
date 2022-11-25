package math

// Pow calculates n to the mth power. Since the result is an int, it is assumed that m is a positive power
func Pow(n, m int) (result int) {
	if n == 1 || m == 0 {
		return 1
	}
	result = n
	for i := 2; i <= m; i++ {
		result *= n
	}
	return
}

func Pow10(n int) int {
	return Pow(10, n)
}

// 最后转 int 时有可能溢出
//func MathPow(n, m int) int {
//	return int(math.Pow(float64(n), float64(m)))
//}
