package os

import (
	"github.com/Weidows/wutils/utils/collection"
	"testing"
)

func TestGetEnumWindowsInfo(t *testing.T) {
	collection.ForEach(GetEnumWindowsInfo(&EnumWindowsFilter{
		IgnoreNoTitled:  true,
		IgnoreInvisible: true,
	}), func(i int, v *EnumWindowsResult) {
		logger.Println(v)
	})
}

func TestSetWindowOpacity(t *testing.T) {
	collection.ForEach(FindWindows("Visual Studio Code"), func(i int, v *EnumWindowsResult) {
		logger.Println(SetWindowOpacity(v.Handle, 200))
	})
}

func TestFindWindows(t *testing.T) {
	collection.ForEach(FindWindows("Visual Studio Code"), func(i int, v *EnumWindowsResult) {
		logger.Println(v)
	})
}
