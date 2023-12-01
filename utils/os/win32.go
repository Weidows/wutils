package os

// github.com/lxn/win

import (
	"github.com/Weidows/wutils/utils/collection"
	"github.com/Weidows/wutils/utils/grammar"
	"github.com/zzl/go-win32api/v2/win32"
	"golang.org/x/sys/windows"
	"os"
	"strings"
	"syscall"
)

var (
	enumWindowsFilter  *EnumWindowsFilter
	_enumWindowsResult []*EnumWindowsResult
)

type EnumWindowsFilter struct {
	IgnoreNoTitled  bool
	IgnoreInvisible bool
}

type EnumWindowsResult struct {
	Handle  win32.HWND
	Title   string
	Opacity byte
}

func GetEnumWindowsInfo(filter *EnumWindowsFilter) []*EnumWindowsResult {
	_enumWindowsResult = nil
	enumWindowsFilter = filter
	EnumWindows(EnumWindowsProc)
	return _enumWindowsResult
}

type TypeEnumWindowsProc func(hwnd win32.HWND, lparam win32.LPARAM) uintptr

// EnumWindows 一进多出, 最好是输入全局变量往里面push, 而不是return
func EnumWindows(callbackFn TypeEnumWindowsProc) {
	if res, err := win32.EnumWindows(syscall.NewCallback(callbackFn), 0); res == 0 && err.NilOrError() != nil {
		logger.Error(err.Error())
	}
}

// EnumWindowsProc 在这里处理每个窗口
func EnumWindowsProc(hwnd win32.HWND, lparam win32.LPARAM) uintptr {
	info := EnumWindowsResult{
		Handle: hwnd,
	}

	buffer := make([]uint16, 1024)
	_, _ = win32.GetWindowText(hwnd, &buffer[0], int32(len(buffer)))
	//[71 68 73 43 32 87 105 110 100 111 119 32 40 84 97 98 84 105 112 46 101 120 101 41 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
	//IsWindowVisible: 0窗口句柄: 131170窗口标题: GDI+ Window (TabTip.exe)
	//[83 104 101 108 108 32 72 97 110 100 119 114 105 116 105 110 103 32 67 97 110 118 97 115 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
	//IsWindowVisible: 0窗口句柄: 66444窗口标题: Shell Handwriting Canvas
	//fmt.Println(utf16.Decode(buffer))
	info.Title = strings.Trim(windows.UTF16ToString(buffer), " ")
	if enumWindowsFilter.IgnoreNoTitled && info.Title == "" {
		return 1
	}
	//这样分开性能更好些
	invisible := win32.IsWindowVisible(hwnd) == 0
	if enumWindowsFilter.IgnoreInvisible && invisible {
		return 1
	}

	//不太靠谱
	//alpha, _ := win32.GetWindowLongA(hwnd, win32.GWL_EXSTYLE)
	if res, err := win32.GetLayeredWindowAttributes(hwnd, nil, &info.Opacity, nil); res != 1 || err.NilOrError() != nil {
		//logger.Println(res, err.NilOrError(), opacity, title)
		info.Opacity = 255
	}

	_enumWindowsResult = append(_enumWindowsResult, &info)
	return 1 // 返回1继续枚举，返回0停止枚举
}

// 输出显示会带有连续空格, 是IDE的问题
func log2File(log string) {
	//IsWindowVisible: 0, 窗口句柄: 68672, 窗口标题: 更换歌单
	//IsWindowVisible: 1, 窗口句柄: 2626206, 窗口标题: wutils – win32.go
	//IsWindowVisible: 1, 窗口句柄: 267984, 窗口标题: win-opacity-js - Visual Studio Code
	//IsWindowVisible: 0, 窗口句柄: 1118748, 窗口标题: 关注创建者
	//IsWindowVisible: 1, 窗口句柄: 132658, 窗口标题: 新建标签页 和另外 34 个页面 - 个人 - Microsoft​ Edge
	//IsWindowVisible: 1, 窗口句柄: 4592372, 窗口标题: Process Hacker [WEIDOWS\Administrator]
	//IsWindowVisible: 0, 窗口句柄: 66830, 窗口标题: ArmourySwAgent
	//IsWindowVisible: 1, 窗口句柄: 3217252, 窗口标题: Blog-private (工作区) - Visual Studio Code
	file, err := os.OpenFile("win32.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(log)
	if err != nil {
		panic(err)
	}
}

// SetWindowOpacity Allow range is [50,255], in order to avoid completely invisible
func SetWindowOpacity(hwnd win32.HWND, opacity byte) bool {
	if opacity < 50 || opacity > 255 {
		return false
	}

	windowLong, err := win32.GetWindowLongA(hwnd, win32.GWL_EXSTYLE)
	if err.NilOrError() != nil {
		logger.Error(err.Error())
	}

	if _, err = win32.SetWindowLongPtrA(hwnd, win32.GWL_EXSTYLE, uintptr(windowLong)|uintptr(win32.WS_EX_LAYERED)); err.NilOrError() != nil {
		logger.Error(err.Error())
	}
	if res, err := win32.SetLayeredWindowAttributes(
		hwnd, 0, opacity, win32.LWA_ALPHA); res == 0 && err.NilOrError() != nil {
		logger.Error(err.Error())
	}
	return true
}

func FindWindows(regex string) (res []*EnumWindowsResult) {
	wins := GetEnumWindowsInfo(&EnumWindowsFilter{
		IgnoreNoTitled:  true,
		IgnoreInvisible: false,
	})
	collection.ForEach(wins, func(index int, value *EnumWindowsResult) {
		if grammar.Match(regex, value.Title) {
			res = append(res, value)
		}
	})
	return res
}
