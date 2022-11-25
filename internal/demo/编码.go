/*
 * @?: *********************************************************************
 * @Author: Weidows
 * @Date: 2022-02-12 00:19:04
 * @LastEditors: Weidows
 * @LastEditTime: 2022-05-25 22:28:43
 * @FilePath: \Golang\src\learn\编码.go
 * @Description:
 * @!: *********************************************************************
 */

package main

import (
	"bytes"
	"fmt"
	"math/big"
)

func main() {
	hello := Base58Encoding("hello")
	fmt.Println(hello)
	fmt.Println(Base58Decoding(hello))
}

var b58 = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

// base58编码
func Base58Encoding(src string) string {
	// he : 104 101 --> 104*256 + 101 = 26725

	// 26725 / 58 = 15 16 17

	// 1.ascii码对应的值

	src_byte := []byte(src)

	// 转成十进制
	i := big.NewInt(0).SetBytes(src_byte)

	var mod_slice []byte
	// 循环取余
	//for i.Cmp(big.NewInt(0)) != 0 {
	for i.Cmp(big.NewInt(0)) > 0 {
		mod := big.NewInt(0)
		i58 := big.NewInt(58)
		// 取余
		i.DivMod(i, i58, mod)
		// 将余数添加到数组中
		mod_slice = append(mod_slice, b58[mod.Int64()])
	}

	// 把0使用字节'1'代替
	for _, s := range src_byte {
		if s != 0 {
			break
		}
		mod_slice = append(mod_slice, byte('1'))
	}

	// 反转byte数组

	// 方法一: 单变量->中间
	for i := 0; i < len(mod_slice)/2; i++ {
		mod_slice[i], mod_slice[len(mod_slice)-1-i] = mod_slice[len(mod_slice)-1-i], mod_slice[i]
	}

	// 方法二: 双变量->中间
	// for i, j := 0, len(mod_slice)-1; i < j; i, j = i+1, j-1 {
	// 	mod_slice[i], mod_slice[j] = mod_slice[j], mod_slice[i]
	// }

	//fmt.Println(mod_slice)
	return string(mod_slice)
}

// base58解码
func Base58Decoding(src string) string {

	// 转成byte数组
	src_byte := []byte(src)
	//fmt.Println(src_byte)

	// 这里得到的是十进制
	ret := big.NewInt(0)
	for _, b := range src_byte {
		//fmt.Println(b)
		i := bytes.IndexByte(b58, b)
		//fmt.Println(i)
		//乘回去
		ret.Mul(ret, big.NewInt(58))
		// 相加
		ret.Add(ret, big.NewInt(int64(i)))

	}

	return string(ret.Bytes())
}
