# Langgo Framework

> Langgo是一款go语言开发应用的框架。


## 目录

- [安装](#安装)
- [快速开始](#快速开始)
- [grpc](#grpc)
- [mysql](#mysql)
- [goreleaser使用](#goreleaser使用)


## 安装

基于go 1.19开发

1. 安装langgo
```
go get -u github.com/langwan/langgo
```

2. 导入
```
import "github.com/langwan/langgo"
```


## 快速开始

```go
package main

import (
	"github.com/langwan/langgo"
	"github.com/langwan/langgo/components/hello"
	"github.com/langwan/langgo/core/log"
)

func main() {
	langgo.Run(&hello.Instance{Message: "hello component"})
	log.Logger("component", "hello").Info().Msg(hello.GetInstance().Message)
}
```

## grpc

grpc支持单机模式和etcd服务发现两种模式，可以参考examples/grpc/single和examples/grpc/etcd两个例子。

## mysql

参考 `examples/mysql`

mysql配置支持多个mysql账号，例如：

```yaml
mysql:
  main:
    dsn: main:123456@tcp(localhost:3306)/simple?charset=utf8mb4&parseTime=True&loc=Local
    conn_max_lifetime: 1h
    max_idle_conns: 1
    max_open_conns: 10
  order:
    dsn: order:123456@tcp(localhost:3306)/simple?charset=utf8mb4&parseTime=True&loc=Local
    conn_max_lifetime: 1h
    max_idle_conns: 1
    max_open_conns: 10

```

这样可以支持项目会拥有多个mysql数据库

```go
langgo.Run(&mysql.Instance{})
var one int
mysql.Main().Raw("SELECT 1").Scan(&one)
fmt.Println(one)
```

`mysql.Main()` 表示获取配置中`main`下的mysql配置，如果想获取`order`，需要使用 `mysql.Get("order")`



## goreleaser使用

```go
//初始化项目
goreleaser init

//发布在本地的版本
goreleaser release --snapshot --rm-dist //删除本地打包文件夹并打包

//发布在github的版本
//1.设置token
export GITHUB_TOKEN="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"

//2.添加tag
git tag -a v0.1.0 -m "First Release"
git push origin v0.1.0

//3.发布
goreleaser release
```

