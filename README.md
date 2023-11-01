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

> Some `demos` and `utils` in learning & developing golang.
>
> 仿照 https://github.com/vuejs/vuepress/tree/master/packages/docs/docs

<a>![分割线](https://cdn.jsdelivr.net/gh/Weidows/Weidows/image/divider.png)</a>



# Utils

一些工具方法, 主要为了解决标准库中棘手/经常复用但没有的情况

# common-starter

格式转换启动, 比如某些程序只能启动 .exe, 通过此程序中转启动其他格式的, 比如 .bat

```shell
go install github.com/Weidows/wutils/cmd/common-starter@latest
```

***

```console
> common-starter ./test.bat
```

# subdir-extracter

```shell
go install github.com/Weidows/wutils/cmd/subdir-extracter@latest
```

解散一级目录

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

# dsg

Disk sleep guard, 防止硬盘睡眠

```shell
go install github.com/Weidows/wutils/cmd/dsg
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

#  gmm

Golang package Mirror Manager

- [x] 结果排序
- [x] 协程加速

***

```shell
go install github.com/Weidows/wutils/cmd/gmm
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

# jpu

Jetbrains Portable Upgrader

```shell
go install github.com/Weidows/wutils/cmd/jpu
```

```shell
echo "Default to use $SCOOP/persist/jetbrains-toolbox/apps"
set JPU_PATH=D:/Scoop/persist/jetbrains-toolbox/apps && jpu.exe
```

***

通过改配置实现 Portable 效果

```
- PyCharm-P
  - ch-0
      - 223.8214.51
          - bin
              - idea.properties
      - 223.8214.51.plugins
      - 223.8617.48
      - 223.8617.48.plugins
  - config
  - system
- Goland
- datagrip
```

在 `IDE/bin/idea.properties` 顶部添加

```properties
idea.config.path=${idea.home.path}/../../config
idea.system.path=${idea.home.path}/../../system
```

由于路径中含带版本号, 用脚本不易操作, 所以用 go 写
