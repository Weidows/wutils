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

  </details>

- <details>

    <summary> jpu | Jetbrains Portable Upgrader</summary>

  ```shell
  go install github.com/Weidows/Golang/cmd/jpu@latest
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

  </details>

<a>![分割线](https://www.helloimg.com/images/2022/07/01/ZM0SoX.png)</a>

# utils

一些工具方法, 主要为了解决标准库中棘手/经常复用但没有的情况

<a>![分割线](https://www.helloimg.com/images/2022/07/01/ZM0SoX.png)</a>

# clean

```shell
go get -u all
go mod tidy
```