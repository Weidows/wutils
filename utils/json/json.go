package json

import (
	"encoding/json"
	"fmt"
	"io"
)

func Marshal(v any) string {
	marshal, err := json.Marshal(v)
	if err != nil {
		fmt.Println(err)
	}
	return string(marshal)
}

func Decode[T any](respBody io.Reader) T {
	var res T
	err := json.NewDecoder(respBody).Decode(&res)
	if err != nil {
		fmt.Println(err)
	}
	return res
}
