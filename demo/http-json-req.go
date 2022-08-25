package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func main() {
	token := "123456"

	jsonBody, _ := json.Marshal(map[string]any{
		"grantType":    "refresh_token",
		"refreshToken": token,
	})
	req := bytes.NewBuffer(jsonBody)

	//resp, _ := http.Post("", "", buff)
	resp := req

	var m any
	_ = json.NewDecoder(resp).Decode(&m)
	//	map[grantType:refresh_token refreshToken:123456]
	fmt.Println(m)
}
