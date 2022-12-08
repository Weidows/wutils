package cast

import (
	"github.com/Weidows/Golang/utils/log"
	"github.com/howcrazy/xconv"
	"strconv"
)

var (
	logger = log.GetLogger()
)

func Convert[T any](src any) (res T) {
	xconv.Convert(src, &res)
	return res
}

func ToIntSlice(i interface{}) (r []int) {
	switch v := i.(type) {
	case []int:
		r = v
	case []int32:
		r = make([]int, len(v))
		for index, value := range v {
			r[index] = int(value)
		}
	case *[]int:
		r = *v
	case []int64:
		r = make([]int, len(v))
		for index, value := range v {
			r[index] = int(value)
		}
	case []string:
		r = make([]int, len(v))
		for index, value := range v {
			r[index], _ = strconv.Atoi(value)
		}
	default:
	}
	return
}

func ToFloat64Slice(i interface{}) (r []float64) {
	switch v := i.(type) {
	case []float64:
		r = v
	case []int32:
		r = make([]float64, len(v))
		for index, value := range v {
			r[index] = float64(value)
		}
	case *[]float64:
		r = *v
	case []int64:
		r = make([]float64, len(v))
		for index, value := range v {
			r[index] = float64(value)
		}
	case []string:
		r = make([]float64, len(v))
		for index, value := range v {
			r[index], _ = strconv.ParseFloat(value, 64)
		}
	default:
	}
	return
}
