package grammar

func ConditionalEqual[T any](condition bool, value1, value2 T) T {
	if condition {
		return value1
	} else {
		return value2
	}
}
