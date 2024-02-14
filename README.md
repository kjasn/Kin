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

## 基于 net/http 标准库实现 web 框架的入口

通过 `http.ListenAndServe()` 启动 web 服务时，第一个参数是 web 服务地址，第二个是一个 Handler 类型的参数。Handler 是一个接口类型，实现了 `ServeHTTP(http.ResponseWriter, *http.Request)` 方法，由此，我们可以自定义一个实现了该方法的接口，将这个接口实例作为第二个参数 （第二个参数为 nil 时表示使用标准库的接口实例）。

自定义一个简单的 Engine 接口：

```go
type Engine struct {
   router map[string]HandlerFunc
}
```

接着实现 `New、GET、POST、Run` 等方法。Run 方法封装 `http.ListenAndServe()`，其中 `ServeHTTP` 方法会在每次服务器接收到请求时被调用，而存储在 router 的中间件会在 `ServeHTTP` 中被调用。

## 设计上下文 Context

1. 用户在每次请求和响应时都要设置消息头(Header)，消息头中包含状态码，消息类型等，实在太麻烦，所以需要进行封装。

Context 中必须要 `*http.Request`和`http.ResponseWriter`用来发送请求和根据请求构造响应。再加上 状态码(StatusCode)、请求路由(Path) 和 请求方法 (Method)，接着实现 `PostForm、Query、SetHeader、String、JSON、HTML`等方法

```go
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string  // eg: GET POST...
	Params map[string]string
	// response info
	StatusCode int
}
```

从 kin.go 中抽离出路由相关的实现放到 router.go 中，同时将 handler 的参数改为 Context 类型。

## 使用 Trie 树存储和查询路由

定义 树结点 结构体如下：

```go
type node struct {
	pattern  string  // complete router path to match
	part     string  // segment of router path at current node
	children []*node // child nodes
	isWild   bool    // contain parameter(:id) or wildcard (*)
}
```

pattern: 完整的请求路由，只在路由段最后一个结点才会设置 pattern，否则为空 eg: `/test/:id/a` ，只有在 `a` 结点才设置 pattern 为 `/test/:id/a` ，

由此可用来判断是否匹配成功： `/test/12`，匹配结束，最后一个结点 `12`的 pattern 为空，即路由表不存在该路由。而 `/test/12/a`，`a`的 pattern 非空，匹配成功。

part: 当前结点的路由段，eg: `/a/b/c`中 `"" 、b、c` 都是
isWild: 用来标记是否为 动态参数 或 通配路由(\*)
path: 实际请求的路由，eg: `/test/123/a` (对应 pattern 的示例路由)
parts: 由 pattern 或 path 按 `/`划分而来。 eg: `/test/:id/a => [test :id a]、/test/123/a => [test 123 a]`

路由的注册和查询由 `insert` 和 `search` 完成，二者都递归查询路由表，但 `insert`查询到一个匹配的结点就立刻返回，`search`则会查询所有匹配的结点，返回一个这个结点数组，然后遍历这些结点继续递归的查询下一层路由，直到查询到完全匹配的路由。

路由注册的顺序很重要，更具体的路由需要在通用的路由（如：参数化的路由）之前注册，举个例子：

```go
router.GET("/:lang", func(ctx *kin.Context) {
   ctx.String(http.StatusOK, "this is a dynamic route")
})

router.GET("/cpp", func(ctx *kin.Context) {
   ctx.String(http.StatusOK, "this is cpp url")
})
```

以上注册了两个路由，先注册的是更通用的路由 `/:lang` ，接着再注册 `/cpp`时，getRoute 会匹配到
