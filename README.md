# Kin

这是一个小型的 Go Web 框架，模仿了 Gin 的设计和功能。实现了如下功能：

1. **Web 框架入口设计**：基于 net/http 标准库实现了 web 框架的入口，提供了自定义 Engine 接口，支持自定义路由和中间件。
2. **上下文 Context 设计**：设计了上下文 Context，用于封装请求和响应，简化消息头设置，并实现了常用的访问和响应方法。
3. **路由管理**：采用 Trie 树存储和查询路由，实现了路由注册和查询功能，支持动态参数和通配符路由。
4. **中间件支持**：实现了中间件机制，允许在请求处理流程前后插入额外的处理逻辑，支持中间件的顺序控制和错误恢复。
5. **模板渲染**：支持模板渲染功能，将请求的地址映射到实际文件存储地址，通过 net/http 库实现静态文件服务。

## 框架大致原型

项目结构如下：

```shell
D:\DEVELOP\GO\GOWORKPLACE\KIN
│  .gitignore
│  go.mod
│  main.go				# 测试文件
│  makefile
│  README.md
├─kin
│      context.go		# 上下文设计 进行请求和响应的封装以及实现常用的访问和响应的方法
│      kin.go			# 框架入口
│      logger.go		# 记录日志的中间件
│      recovery.go		# 错误恢复的中间件
│      router.go		# 将从 kin.go 中抽离的 router 方法实现
│      router_test.go	# 单元测试
│      trie.go			# 通过 trie 树存储和查询路由
│
└─static				# 存放本地文件
        file1.jpg
        file2.md
        template.html
```

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
	// ...
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

## 中间件

中间件类似路由处理函数(HandlFunc)，区别在于中间件返回的是一个闭包。中间件保存在 `Context` 中，因为中间件不仅作用在处理流程前，也可以作用在处理流程后，即在用户定义的 Handler 处理完毕后，还可以执行剩下的操作。
中间件通过 `Next()` 方法递归的触发，由索引来标示顺序。每次调用 `Next()` ，控制权就交给下一个中间件。

## 模板

框架需要做的是将请求的地址映射到文件实际的存储地址，接着找到文件后，如何返回这一步，net/http 库已经实现了。

eg: 我们将静态文件放在 `/assets/` 下， 服务上文件存储在 `./static/`，接着将 `./static/` 映射到 `/assets/`，访问 `localhost/assets/file` 时就会解析为 `./static/file` （file 为 static 路径下文件的相对路径）

在 `Engine` 中加上以下两个字段 `*template.Template` 和 `template.FuncMap` 对象，前者存储全局加载的模板，后者存储自定义的渲染函数。

```go
type Engine struct {
	// ...
	// serve as html render
	htmlTemplates *template.Template	// store all html templates
	funcMap template.FuncMap	// render func
}
```

## 错误恢复

由于我们在处理错误时都是之间 `panic(err)`，这样可能由于错误的请求使得服务器宕机，为避免这种情况，我们使用 `recover()` 来恢复错误。
每当错误发生时 `panic(err)` 之前会处理 `defer` 的任务，因此我们可以在 `defer` 中使用 `recover()` 来进行错误恢复。
