# Kin

This is a simple Go web framework that mimics the design and functionality of Gin.
Studying from [@极客兔兔](https://geektutu.com/post/gee.html)

这是一个简单的 Go Web 框架，模仿了 Gin 的设计和功能。跟着 @极客兔兔 的博客学。

## 框架大致原型

1. main.go : 测试用
2. router_test.go : 单元测试

3. kin.go

作为框架入口, 抽离出 router api 放到 router.go 中实现

4. context.go

- 封装 http.ResponseWriter 和 http.Request 以及相关的方
- 实现常用的访问 Query 和 PostForm 参数方法
- 实现 String/JSON/HTML 响应方法

5. router.go
   将从 kin.go 中抽离的 router 方法实现

6. trie.go

通过 trie 树存储和查询路由
