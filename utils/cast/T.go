package cast

import (
	"reflect"
)

type SupportedEmptyType interface {
	any
}

// EmptyT
//
// 已知类型参数 T, 想要获取 T 类型的空值
// 用来解决: 无法将 'nil' 用作类型 T
// 注意类型参数不能填 any | interface{}, 会有无法处理的空指针异常
func EmptyT[T SupportedEmptyType]() (t T) {
	var res any

	//reflect.TypeofT(t, func(a any) {
	//	switch a.(type) {
	//	case int:
	//		res = 0
	//	}
	//})

	switch reflect.TypeOf(t).Kind() {
	case reflect.String:
		res = ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		res = reflect.Zero(reflect.TypeOf(t)).Interface()
	case reflect.Bool:
		res = false
	case reflect.Slice, reflect.Array:
		res = reflect.MakeSlice(reflect.TypeOf(t), 0, 0).Interface()
	case reflect.Chan:
		res = reflect.MakeChan(reflect.TypeOf(t), 0).Interface()
	case reflect.Map:
		res = reflect.MakeMap(reflect.TypeOf(t)).Interface()
	case reflect.Func:
		res = reflect.MakeFunc(reflect.TypeOf(t), func(args []reflect.Value) (results []reflect.Value) {
			return args
		}).Interface()
	//case reflect.Interface:
	default:
		res = nil
	}

	return res.(T)
}
