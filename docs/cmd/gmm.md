#  gmm

Golang package Mirror Manager

- [x] 结果排序
- [x] 协程加速

***

```shell
go install github.com/Weidows/Golang/cmd/gmm@latest
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