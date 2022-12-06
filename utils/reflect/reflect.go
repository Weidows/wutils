package reflect

// Typeof
//
// https://studygolang.com/articles/14269
func Typeof(v any, fn func(any)) {
	fn(v)
}
