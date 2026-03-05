---
title: 🎉Docs-wutils
password: ""
tags:
  - tools
  - doc
  - wutils
katex: false
comments: true
aside: true
date: 2024-08-06 04:47:19
top_img:
cover: https://pan.weidows.tech/d/local/blog/1d36e9d50555af6fca23e5fd36246cf5490809012.jpg
---

> https://github.com/Weidows/wutils

# Docs-wutils

<!--
 * @!: *********************************************************************
 * @Author: Weidows
 * @LastEditors: Weidows
 * @Date: 2022-08-30 14:51:11
 * @LastEditTime: 2026-03-03 10:59:18
 * @FilePath: \wutils\README.md
 * @Description:
 * @:
 * ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡏⠻⣄
 * ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣀⣀⣀⣀⣀⠀⠀⠀⠀⠀⣀⡾⠀⠀⣿
 * ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣠⡶⠟⠛⠉⠉⠀⠀⠀⠀⠀⠀⠉⠉⠛⠻⠿⣤⣀⣠⡴⠋
 * ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⠟⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⣤⠂⠀⠀⠀⠀⠀⠀⠀⠀⠙⢶⣀
 * ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡾⠉⠀⠀⠀⠀⠀⠀⡴⠁⠀⠀⠀⡞⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠳⣄
 * ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣴⠋⠀⠀⠀⠀⠀⠀⢠⠏⠀⠀⠀⠀⡾⠀⠀⠀⠀⠀⠀⠀⠀⢰⠀⠀⠀⠀⠀⠀⠈⢷⡀
 * ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡿⠀⠀⠀⠀⠀⠀⠀⣠⠃⠀⠀⠀⠀⢰⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⢿⠀⠀⠀⠀⠀⠀⠀⠙⣄
 * ⠀⠀⠀⠀⠀⠀⠀⠀⠀⡿⠀⠀⠀⠀⠀⠀⠀⠀⣿⠀⠀⠀⠀⣀⣿⢠⠀⠀⠀⠀⠀⠀⠀⠀⣶⠀⣧⠀⠀⠀⠀⠀⠀⠀⠈⣆
 * ⠀⠀⠀⠀⠀⠀⠀⠀⣼⠁⠀⣰⠀⠀⠀⠀⣾⢸⣿⣀⣤⣶⡛⠁⣿⠘⡄⠀⠀⠀⠀⢀⠀⠀⣿⠀⢻⠀⠀⠀⠀⠀⠀⠀⠀⠸
 * ⠀⠀⠀⠀⠀⠀⠀⠀⣿⠀⠀⡟⠀⠀⠀⠀⣿⡈⡏⠀⠀⠀⠀⠙⡜⡄⢿⠀⠀⠀⠀⣾⣧⣾⢸⠀⢸⠀⠀⡇⠀⠀⠀⠀⠀⠀⡇
 * ⠀⠀⠀⠀⠀⠀⠀⢀⡇⠀⠀⣇⠀⠀⠀⠀⣷⣿⣿⣶⣦⣀⠀⠀⠀⠹⡘⣿⡀⠀⠀⢹⣿⠃⢸⠀⣿⠀⢠⣿⠀⠀⠀⠀⠀⣄⣧
 * ⠀⠀⠀⠀⠀⠀⠀⣼⠀⠀⠀⣿⠀⠀⠀⠀⢿⠉⠉⠉⠻⣿⣷⡀⠀⠀⠈⠁⠙⢦⣀⣸⠋⠈⣸⣼⢻⠀⣾⡿⠀⠀⠀⠀⢸⢸⣿
 * ⠀⠀⠀⠀⠀⠀⠀⣿⠀⣿⠀⠸⡄⣤⠀⠀⠀⣧⠀⠀⠀⠀⠉⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠀⠀⣿⣿⠇⠀⠀⠀⠀⣾⢸⣿
 * ⠀⠀⠀⠀⠀⠀⢰⡇⠀⣿⠀⠀⠙⠿⣿⠛⠒⠛⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣾⣿⣿⣿⣿⣦⡀⡿⠀⠀⠀⠀⠀⣿⣿⣿
 * ⠀⠀⠀⠀⠀⠀⣿⠀⣴⠋⠙⣦⠀⠀⠀⣇⠀⠀⠀⠀⠀⠀⣼⠉⠙⠳⠦⣤⡀⠀⠀⠀⠀⠀⠀⠈⠻⡿⠀⠀⠀⠀⠀⠀⣿⣿⣿
 * ⠀⠀⠀⠀⠀⠀⣿⢤⠇⣀⡀⣼⠀⠀⠀⢹⠀⠀⠀⠀⠀⢸⠁⠀⠀⠀⠀⢰⠁⠀⠀⠀⠀⠀⠀⣠⠟⠀⠀⠀⠀⠀⠀⢠⣯⣿⣿
 * ⠀⠀⠀⠀⣠⠋⣠⠟⠁⠀⠈⠛⣄⠀⠀⠀⣇⠀⠀⠀⠀⢹⠀⠀⠀⠀⢀⠟⠀⠀⠀⠀⠀⠀⢾⣁⠀⠀⠀⠀⢀⣾⣤⡏⠉⠉⣿
 * ⠀⠀⢀⣞⠤⣴⠁⠀⠀⠀⠀⠀⠀⣧⠀⠀⢿⠉⠳⣤⡀⠀⠁⠀⠠⠶⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠧⣴⣾⣿⣿⠟⡇⠀⠀⢸
 * ⣠⠚⠁⢀⠏⣿⠀⠀⠀⠀⠀⠀⠀⠸⡀⠀⠈⡏⠛⠶⣤⡉⠛⠶⢤⣤⣀⡀⠀⠀⠀⠀⠀⠀⠀⣀⡴⣿⣿⠟⠁⠀⠀⡇⠀⠀⢸
 * ⠀⠀⠀⢸⠀⠘⡄⠀⠀⠀⠀⠀⢹⠀⡿⡄⠀⣿⠀⠀⠀⠀⠙⠲⣤⡀⠙⢦⠉⠉⠉⠉⠉⠉⠁⠀⠀⣿⠀⠀⠀⠀⠀⡇⠀⠀⢸
 * ⠀⠀⠀⠘⡄⠀⣿⠀⠀⠀⠀⠀⣸⠛⡄⢻⠀⠸⡀⠀⠀⠀⠀⠀⠀⠈⠛⢦⡉⢦⡀⠀⠀⠀⠀⠀⠀⣿⠀⠀⠀⠀⢰⠇⠀⠀⢸
 * ⠀⠀⠀⠀⢿⠀⠈⠙⢦⣄⣠⠴⠃⢠⠃⠀⡇⠀⣇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⢾⣦⠀⠀⠀⠀⠀⣿⠀⠀⠀⠀⣼⠀⠀⠀⠀⡇
 * ⠀⠀⠀⠀⠀⠙⢤⠀⠀⠀⠀⠀⡴⠋⠀⠀⡇⠀⢿⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠳⣄⠀⠀⠀⡿⠀⠀⠀⠀⡿⠀⠀⠀⠀⡇
 * ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⡇⠀⠘⣿⣿⣄⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠳⠀⠀⡇⠀⠀⠀⢰⠁⠀⠀⠀⠀⡇
 * ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡾⠀⠀⠀⣿⣿⣿⣿⣷⣶⣶⣤⣄⠀⠀⠀⠀⠀⠀⠀⠀⢸⠁⠀⠀⠀⡿⠀⠀⠀⠀⠀⣷
 *
 * @?: *********************************************************************
-->

{% pullquote mindmap mindmap-sm %}

- [Docs-wutils](#docs-wutils)
  - [install](#install)
  - [Cmd](#cmd)
    - [wutils](#wutils)
      - [parallel](#parallel)
      - [ol-Opacity\_Listener](#ol-opacity_listener)
        - [特点-feature](#特点-feature)
        - [教程-tutorial](#教程-tutorial)
      - [dsg-Disk\_sleep\_guard](#dsg-disk_sleep_guard)
      - [diff](#diff)
      - [zip](#zip)
        - [crack](#crack)
      - [media](#media)
      - [subdir-extracter](#subdir-extracter)
        - [usage](#usage)
        - [example](#example)
      - [gmm](#gmm)
      - [buffer](#buffer)
        - [功能特性](#功能特性)
        - [用法](#用法)
        - [选项](#选项)
        - [配置文件](#配置文件)
        - [注意事项](#注意事项)
  - [Pkg](#pkg)
    - [zip](#zip-1)
  - [Utils](#utils)
  - [借物表](#借物表)

{% endpullquote %}

<a>![分割线](https://pan.weidows.tech/d/local/img/divider.png)</a>

## install

```shell
go install github.com/Weidows/wutils/cmd/wutils@master

# or use scoop
scoop install wutils

# or, just download from release:
https://github.com/Weidows/wutils/releases
```

> 默认配置 | default config: https://github.com/Weidows/wutils/tree/master/config \
> 如果是 scoop / release 安装, 配置会在压缩包内

<a>![分割线](https://pan.weidows.tech/d/local/img/divider.png)</a>

## Cmd

一些命令行程序, 基本都是我自己需要用的, 没找到符合需求的就简单写一个, 针对 `服务 (service)`

### wutils

- [x] 运行时配置热更新

```console
╰─ 👉wutils
Failed to find configuration ./config/cmd/wutils.yml
NAME:
   wutils - Documents(使用指南) at here:
            https://blog.weidows.tech/post/lang/golang/wutils

USAGE:
   wutils [global options] command [command options]                                                                                                                        AUTHOR: Weidows <ceo@weidows.tech>

COMMANDS:
   config        print config file
   diff          diff - Differential set between two files
                 文件对比工具, 但不是 Git-diff 那种
                 是用来求 '行-差集' 的工具
                 输入为两个特定名称的文件: './inputA.txt', './inputB.txt'
   parallel, pl  并行+后台执行任务 (配置取自wutils.yml)
   dsg           Disk sleep guard
                 防止硬盘睡眠 (每隔一段自定义的时间, 往指定盘里写一个时间戳)
                 外接 HDD 频繁启停甚是头疼, 后台让它怠速跑着, 免得起起停停增加损坏率
   ol            Opacity Listener
                 后台持续运行, 并每隔指定时间扫一次运行的窗口
                 把指定窗口设置opacity, 使其透明化 (same as BLend)
   zip           some actions to operate zip/7z files
   media         some actions to operate image or video files
   buffer        Buffer filesystem - 基于 Dokan 的 IO 缓冲虚拟文件系统
   help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

#### parallel

wutils 部分子程序设计为根据 [配置文件](https://github.com/Weidows/wutils/blob/master/config/cmd/wutils.yml) 持续运行的

所以单独出来一个配置项控制子程序后台一起运行

---

#### ol-Opacity_Listener

后台持续运行并控制窗口透明化的程序, 如图:

![1702005541764](https://pan.weidows.tech/d/local/blog/1702005541764.jpg)

只在 windows 平台使用和测试过, 其他平台可能会有 bug

可用于类比的软件是 [BLend](https://zhutix.com/tools/touming-blend/), 那个软件太老了, 总是出一些 bug, 所以自己写了一个

##### 特点-feature

- [x] 现在大多数能搜到的工具比较手动, 每个新开的窗口都要再手动设置一次, wutils 只需要改一次配置, 开多少窗口都能立即起效
- [x] 参数可调, 同一程序的透明度统一控制, 不同程序的透明度分开控制
- [x] 配置化, 免得每次启动调来调去 (config/app.yml)

匹配规则是按标题来的, 所以按着上面改就行, 挺直观的

- title 是 regex 字符串

  `^` 是从开头匹配, `$` 是结尾

  比如我的 Edge 浏览器, 通过 `wutils ol list` 命令找到它的标题是这么长 `Weidows/wutils: Some demos and utils in learning \u0026 developing golang. 和另外 154 个页面 - 个人 - Microsoft​ Edge`

  我想让 wutils 匹配以 "Weidows" 开头, 以 "- Microsoft​ Edge" 结尾的窗口, 那应该填 `^Weidows.*- Microsoft​ Edge$`

- opacity 是透明度, `(0,255]`

  一般设置在 200~240 这个范围比较舒服

##### 教程-tutorial

{% mmedia "bilibili" "bvid:BV1d94y1j7JC" %}

---

#### diff

自然, 可以通过 Git 和一些类似的工具实现 '行-差异' 的做法, 但是它们并不能输出, 作为差集

test case at [here](https://github.com/Weidows/wutils/tree/master/cmd/wutils/diff/test)

```console
> ./wutils diff
================== Missing in A ==================
onlyB1
onlyB2

================== Missing in B:==================
onlyA1
onlyA2
```

---

#### zip

##### crack

简而易懂, 破解压缩包

- 自动在 `~/.config/wutils/password-dict.txt` 创建默认密码字典
- 使用协程高速处理
  - `>=1000 test/s`
  - `~50%` CPU usage
  - `nMB ~ nGB` RAM usage
- supporting
  - .zip
  - .7z
  - .rar
  - .tar.gz
  - .tar.bz2
  - 分卷文件 (.zip, .z01, .z02 ....)
  - 分卷 7z 文件 (.7z.001, .7z.002 ...)

---

#### media

场景: 手机拍出来的图片/视频会同一堆在 `DCIM/Camera` 里, 有点太多了难以分组

此程序作用为归类图片/视频, 默认参数对 `12h时间内` & `方圆1km` 归为一组, 减少手动分组的麻烦

只接收一个参数, 图片文件夹的路径, 会在其内部生成 `output`, 并把分组后的文件复制进去 (注意大小)

```
wutils media group F:/Pictures/@Collections/DCIM/Camera
```

---

#### subdir-extracter

主要功能为解散一级目录

- [x] 支持递归文件夹
- [x] 重复文件进行 hash 校验 (前缀重命名法不会有误删, 删除的是完全一致的文件)
- [ ] 提取为 Lib 到 utils 可供调用

##### usage

```
subdir-extracter 0 ./1
```

params:

1. mode
   autoCheck = "0"
   overwrite = "1"
   skip = "2"
2. path
   input the root-dir-path where you need to extracter subdirs

---

##### example

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

适用场景很单一: 比如一大批图包/数据集, 嵌套了很多层让人不舒服 (n > 10 时手动挪出来就很头疼了)

也没有对应工具可以用, 所以就写了这个

---

#### gmm

Golang package Mirror Manager

- [x] 结果排序
- [x] 协程加速

---

#### buffer

Buffer filesystem - 基于 Dokan 的 IO 缓冲虚拟文件系统

用于减少 HDD IOPS，对小文件随机读写场景特别有效。

##### 功能特性

- 基于 Dokan (Windows FUSE) 的透明缓冲
- 写缓冲（write-behind）：减少写操作次数
- 读缓冲（LRU + read-ahead）：热点数据缓存和预取
- 可配置的缓冲策略：monitoring, defrag, download, migration, balanced
- 支持配置文件管理和热更新

##### 用法

```bash
# 挂载缓冲盘 X:，指向 D:/data
wutils buffer mount X: --source D:/data

# 使用自定义内存限制（100MB）
wutils buffer mount X: --source D:/data --memory-limit 104857600

# 使用自定义刷新间隔（30秒）
wutils buffer mount X: --source D:/data --flush-interval 30

# 查看状态
wutils buffer status

# 卸载缓冲盘
wutils buffer unmount
```

##### 选项

- `-s, --source` (必需): 源路径，需要缓冲的目录
- `--memory-limit`: 内存限制，单位字节（默认: 67108864，64MB）
- `--flush-interval`: 刷新间隔，单位秒（默认: 10）
- `--strategy`: 缓冲策略
  - `monitoring`: 监控模式，低延迟
  - `defrag`: 碎片整理模式，大批量写
  - `download`: 下载模式，顺序写
  - `migration`: 迁移模式，批量读写
  - `balanced`: 平衡模式（默认）
- `--enable-read-cache`: 启用读缓存（默认: false）
- `--enable-write-buffer`: 启用写缓冲（默认: true）

##### 配置文件

在 `~/.config/wutils/app.yml` 中配置：

```yaml
cmd:
  buffer:
    enable: false
    memory_limit: 67108864  # 64MB
    flush_interval: 10      # 10秒
    strategy: balanced      # 策略
```

##### 注意事项

- 需要安装 Dokan 驱动（https://github.com/dokan-dev/dokany）
- 仅支持 Windows 平台
- 建议配合 SSD 使用以获得最佳性能

<a>![分割线](https://pan.weidows.tech/d/local/img/divider.png)</a>

## Pkg

一般是 New 出来用的包, 针对 `对象 (object)`

### zip

用于破解压缩文件的包, 上面 cmd 中有调用

<a>![分割线](https://pan.weidows.tech/d/local/img/divider.png)</a>

## Utils

工具类函数, 多为应对 code 时常用却又棘手的情况, 针对 `类型 (type)`

比较偏向 go 的注释即文档做法, utils 里有各种工具库, 可以先装一下, 开发时说不定起手就用到了

`现存函数大大大概率不会删/改名`, base 大致搭好了, 会有 deprecated / break change

## 借物表

暂无.
