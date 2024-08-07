package collection

import "github.com/Weidows/wutils/utils/cast"

// check the diff items between two slices
func SliceDiff[T comparable](a, b []T) (missInA []T, missInB []T) {
	missInA = cast.EmptyT[[]T]()
	missInB = cast.EmptyT[[]T]()

	if len(a) == 0 {
		missInA = b
	} else if len(b) == 0 {
		missInB = a
	} else {
		ma := Slice2Map(a, func(item T) T { return item })
		mb := Slice2Map(b, func(item T) T { return item })

		for _, v := range a {
			_, exists := mb[v]
			if !exists {
				missInB = append(missInB, v)
			}
		}

		for _, v := range b {
			_, exists := ma[v]
			if !exists {
				missInA = append(missInA, v)
			}
		}
	}
	return
}
