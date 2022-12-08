<!--
 * @?: *********************************************************************
 * @Author: Weidows
 * @Date: 2022-08-30 14:51:11
 * @LastEditors: Weidows
 * @LastEditTime: 2022-12-08 18:25:00
 * @FilePath: \Golang\README.md
 * @Description:
 * @!: *********************************************************************
-->

> Some `demos` and `utils` in learning & developing golang.

<a>![分割线](https://www.helloimg.com/images/2022/07/01/ZM0SoX.png)</a>

# cmd

- <details>

    <summary> dsg | Disk sleep guard, 防止硬盘睡眠 </summary>

  ```shell
  go install github.com/Weidows/Golang/cmd/dsg@latest
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

  </details>

- <details>

    <summary> gmm | Golang package Mirror Manager</summary>

  ```shell
  go install github.com/Weidows/Golang/cmd/gmm@latest
  ```

  ***

  - [x] 结果排序
  - [x] 协程加速

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

  </details>

<a>![分割线](https://www.helloimg.com/images/2022/07/01/ZM0SoX.png)</a>

# utils

一些工具方法, 主要为了解决标准库中棘手/经常复用但没有的情况
