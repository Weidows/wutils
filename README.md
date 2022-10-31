<!--
 * @?: *********************************************************************
 * @Author: Weidows
 * @Date: 2022-08-30 14:51:11
 * @LastEditors: Weidows
 * @LastEditTime: 2022-09-01 18:57:17
 * @FilePath: \Gmm\README.md
 * @Description:
 * @!: *********************************************************************
-->

> Some `demos` and `utils` in learning & developing golang.

# Demo-gmm

- <details>

    <summary> Golang package Mirror Manager </summary>

  ***

    ```console
    ╰─ gmm test
    proxys
            352ms   aliyun
            278ms   proxy-cn
            642ms   proxy-io
            269ms   baidu
            1002ms  tencent
            406ms   huawei
            837ms   default
    sumdbs
            2073ms  default
            789ms   google
            1957ms  sumdb-io
    ```
    
    ```console
    ╰─ gmm proxy huawei                                                                                                pwsh   95  12:17:56  
    Proxy use huawei https://repo.huaweicloud.com/repository/goproxy
    
    ╰─ gmm sumdb default                                                                                              pwsh   95  12:17:17 
    Sumdb use default https://sum.golang.org
    ```

  </details>

```shell
go install github.com/Weidows/Golang/tree/master/demo/gmm@latest
```

<a>![分割线](https://www.helloimg.com/images/2022/07/01/ZM0SoX.png)</a>

# Utils

一些工具方法, 主要为了解决标准库中棘手/经常复用但没有的情况