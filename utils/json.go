package utils

import (
	"encoding/json"
	"fmt"
)

func Marshal(v any) string {
	marshal, err := json.Marshal(v)
	if err != nil {
		fmt.Println(err)
	}
	return string(marshal)
}
