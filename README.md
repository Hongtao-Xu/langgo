# 1.goreleaser的使用

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

