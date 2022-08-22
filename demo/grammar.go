/*
 * @?: *********************************************************************
 * @Author: Weidows
 * @Date: 2022-02-11 22:25:54
 * @LastEditors: Weidows
 * @LastEditTime: 2022-05-25 22:32:20
 * @FilePath: \Golang\src\learn\grammar.go
 * @Description:
 * @!: *********************************************************************
 */

package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	// test_print()
	// test_string()
	// test_var()
	// fmt.Println(test_func("a", "b"))
	// test_pointer()
	// test_struct()
	//test_array_slice_map()
	test_error()
}

func test_string() {
	var stockcode = 123
	var enddate = "2020-12-31"
	var url = "Code=%d&endDate=%s"
	var target_url = fmt.Sprintf(url, stockcode, enddate)
	fmt.Println(target_url)
} // Code=123&endDate=2020-12-31

func test_var() {
	// 可以将 var f string = "Runoob" 简写为 f := "Runoob"：
	// 只能被用在函数体内，而不可以用于全局变量的声明与赋值
	aa := "Runoob"

	/*
		和 python 很像,不需要显示声明类型，自动推断
		var vname1, vname2, vname3 = v1, v2, v3
		出现在 := 左侧的变量不应该是已经被声明过的，否则会导致编译错误
	*/
	vname1, vname2, vname3 := 1, true, "v3"

	// 这种因式分解关键字的写法一般用于声明全局变量
	var (
		vname4 int
		vname5 bool
	)

	println(aa, vname1, vname2, vname3, vname4, vname5) // Runoob 1 true v3 0 false

	// ===============================================================

	/*
		const b string = "abc"
		也可省略类型
		const b = "abc"
	*/
	const c_name1, c_name2 = "value1", "value2"

	const (
		Unknown = 0
		Female  = 1
		Male    = 2
	)

	// iota
	const (
		a = iota //0
		b        //1
		c        //2
		d = "ha" //独立值，iota += 1
		e        //"ha"   iota += 1
		f = 100  //iota +=1
		g        //100  iota +=1
		h = iota //7,恢复计数
		i        //8
	)
	fmt.Println(a, b, c, d, e, f, g, h, i) // 0 1 2 ha ha 100 100 7 8
}

// go 支持多返回值
func test_func(x, y string) (string, string) {
	return y, x
}

func test_pointer() {
	a := 20
	p := &a
	println(a, *p, p)

	// 判断空指针
	var ptr *int
	println(ptr, ptr == nil)
}

func test_struct() {
	type Books struct {
		title   string
		author  string
		subject string
		book_id int
	}

	fmt.Println(Books{"Go 语言", "www.runoob.com", "Go 语言教程", 6495407})
	// 忽略的字段为 0 或 空
	fmt.Println(Books{title: "Go 语言", author: "www.runoob.com"})

	Book1 := Books{title: "Go 语言", author: "www.runoob.com"}
	fmt.Printf("Book1: %v\n", Book1)
}

func test_array_slice_map() {
	// var arr = [5]int{1, 2, 3, 4, 5}
	// 数组大小是不能改变的
	arr := [5]int{1, 2, 3, 4, 5}
	fmt.Println("arr:", arr)

	// ====================== slice ======================
	// 长度为3, 容量为5的slice; 切片可以自动扩容
	slice1 := make([]float64, 3, 5)
	fmt.Println("slice1:", slice1, len(slice1), cap(slice1)) // slice1: [0 0 0]

	slice1 = append(slice1, 1, 2, 3, 4, 5)
	fmt.Println(slice1) // [0 0 0 1 2 3 4 5]

	//子切片
	sub_slice_1 := slice1[:3]
	sub_slice_2 := slice1[3:5] // [3,5) 左闭右开
	combined_sub_slice := append(sub_slice_1, sub_slice_2...)
	fmt.Println(combined_sub_slice) // [0 0 0 1 2]

	// ==================== map ======================
	// 仅声明
	m1 := make(map[string]int)
	// 声明时初始化
	m2 := map[string]int{
		"Sam":   11,
		"Alice": 34,
	}

	// 赋值/修改
	m1["Tom"] = 18

	fmt.Println(m2["Sam"])
}

func test_error() {
	_, err := os.Open("filename.txt")
	if err != nil {
		fmt.Print(err) //	open filename.txt: The system cannot find the file specified.
		fmt.Println(errors.New("文件不存在"))
	}

	// 类似 try-catch  当此func发生panic,会跳到此 defer func 执行
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Some error happened!", r)
		}
	}()

	// pannic
	arr := [3]int{2, 3, 4}
	index := 5
	fmt.Println(arr[index]) // Some error happened! runtime error: index out of range [5] with length 3
}
