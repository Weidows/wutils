<!--
 * @?: *********************************************************************
 * @Author: Weidows
 * @Date: 2022-08-30 14:51:11
 * @LastEditors: Weidows
 * @LastEditTime: 2023-03-31 18:30:02
 * @FilePath: \Weidows\Golang\README.md
 * @Description:
 * @!: *********************************************************************
-->
# cmd

以下 cmd-cli 已在 release 打包多平台二进制版本

## common-starter

格式转换启动, 比如某些程序只能启动 .exe, 通过此程序中转启动其他格式的, 比如 .bat

```shell
go install github.com/Weidows/wutils/cmd/common-starter@latest
```

***

```console
> common-starter ./test.bat
```

<a>![分割线](https://cdn.jsdelivr.net/gh/Weidows/Weidows/image/divider.png)</a>

## subdir-extracter

```shell
go install github.com/Weidows/wutils/cmd/subdir-extracter@latest
```

主要功能为解散一级目录

- [x] 支持递归文件夹
- [x] hash 校验

如下为解散前后, 测试文件在 `cmd/subdir-extracter/1`

```
D:\DESKTOP\1
│  2.3.txt
│
├─2.1
│  │  3.1.txt
│  │  3.2.txt
│  │
│  └─3.3
└─2.2
    │  3.1.txt
    │  3.2.txt
    │
    └─2.2
```

```
D:\DESKTOP\1
│  2.2-3.2.txt
│  2.3.txt
│  3.1.txt
│  3.2.txt
│
├─2.2
└─3.3
```

<a>![分割线](https://cdn.jsdelivr.net/gh/Weidows/Weidows/image/divider.png)</a>

## dsg

Disk sleep guard, 防止硬盘睡眠

外接HDD频繁启停甚是头疼, 后台挂着dsg让它一直怠速跑着

```shell
go install github.com/Weidows/wutils/cmd/dsg@latest
```

***

```console
> .\dsg.exe
please start with params like: 'dsg.exe E: 30'
      1. disk
      2. delay seconds


> dsg D:
10 / 30 [------------------->_______________________________________] 33.33%
```

<a>![分割线](https://cdn.jsdelivr.net/gh/Weidows/Weidows/image/divider.png)</a>

## gmm

Golang package Mirror Manager

- [x] 结果排序
- [x] 协程加速

***

```shell
go install github.com/Weidows/wutils/cmd/gmm@latest
```

***

```console
> gmm test
proxy
      125ms   huawei
      178ms   baidu
      219ms   aliyun
      338ms   proxy-cn
      476ms   default
      612ms   proxy-io
      623ms   tencent
sumdb
      433ms   google
      451ms   default
      743ms   sumdb-io
```

```console
╰─ gmm proxy huawei
Proxy use huawei https://repo.huaweicloud.com/repository/goproxy

╰─ gmm sumdb default
Sumdb use default https://sum.golang.org
```