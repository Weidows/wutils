package reflect

// TypeofT
//
// https://studygolang.com/articles/14269
// 这个看似简单的方法主要是用来处理 T 类型的变量
// 因为没办法直接对其 t.(type), 所以在此方法通过 any 类型参数中转
// 使其可以在 fn 中被 switch t.(type)
func TypeofT(v any, fn func(any)) {
	fn(v)
}
