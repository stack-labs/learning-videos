# Micro API

>
> - 准备工做
>   - Go: `v1.13`
>	- micro工具:https://github.com/micro/micro#install
>	- tutorials代码:git clone git@github.com:micro-in-cn/tutorials.git $GOPATH/src/github.com/micro-in-cn/tutorials
>	- learning-videos代码:git clone git@github.com:micro-in-cn/learning-videos.git $GOPATH/src/github.com/micro-in-cn/learning-videos
>	- [etcd](https://etcd.io/)
>	- [consul](https://www.consul.io/)

## 微服务架构-网关

<img src="/docs/Micro API/img/micro-arch.png" width="75%">

## micro网关

#### 演示: 启动网关

**API**

> 前提已经安装micro工具，演示版本`1.18.0`

```bash
micro api
```

```bash
curl http://localhost:8080                                                                                 
{"version": "1.18.0"}
```

**Web**

> 虽然定义的不是网关，但也可以作为网关使用

```bash
micro web
```

访问服务:http://localhost:8082

### Options

**global options**
```bash
micro -h
--registry value                Registry for discovery. etcd, mdns [$MICRO_REGISTRY]
--registry_address value        Comma-separated list of registry addresses [$MICRO_REGISTRY_ADDRESS]
--server_name value             Name of the server. go.micro.srv.example [$MICRO_SERVER_NAME]
--transport value               Transport mechanism used; http [$MICRO_TRANSPORT]
```

**command options**
```bash
micro api -h
--address value    Set the api address e.g 0.0.0.0:8080 [$MICRO_API_ADDRESS]
--handler value    Specify the request handler to be used for mapping HTTP requests to services; {api, event, http, rpc} [$MICRO_API_HANDLER]
--namespace value  Set the namespace used by the API e.g. com.example.api [$MICRO_API_NAMESPACE]
--resolver value   Set the hostname resolver used by the API {host, path, grpc} [$MICRO_API_RESOLVER]
--enable_rpc       Enable call the backend directly via /rpc [$MICRO_API_ENABLE_RPC]
```

> 有关`handler`、`resolver`、`rpc`的介绍参考官方文档：[API Gateway](https://micro.mu/docs/api.html#handlers)

```bash
micro web -h
--address value    Set the web UI address e.g 0.0.0.0:8082 [$MICRO_WEB_ADDRESS]
--namespace value  Set the namespace used by the Web proxy e.g. com.example.web [$MICRO_WEB_NAMESPACE]
```

#### 演示: options

**global options --registry**
```bash
micro --registry=etcd api
micro --registry=etcd web
```

**command options --address & --namespace**

> 教程演示用例[micro-in-cn/tutorials/examples/basic-practices/micro-api](https://github.com/micro-in-cn/tutorials/tree/master/examples/basic-practices/micro-api)

```bash
micro --registry=etcd api --address=:9080 --namespace=com.hbchen.api

# micro-in-cn/tutorials/examples/basic-practices/micro-api/meta
# 错误示范，不指定server_name，
go run meta.go --registry=etcd
# 正确示范
go run meta.go --registry=etcd --server_name=com.hbchen.api.example
```

```bash
curl -XGET "http://localhost:9080/example?name=john"
curl -XPOST -H 'Content-Type: application/json' -d '{"name": "john"}' "http://localhost:9080/example"
```

### 服务发现

> 自定义`namespace`适合启动不同类型的`API`

<img src="/docs/Micro API/img/micro-ds.png" width="75%">

### 路由

**Handler**

| - | 类型 | 说明
----|----|----
1 | rpc | 通过RPC向go-micro应用转送请求，只接收GET和POST请求，GET转发`RawQuery`，POST转发`Body`
2 | api | 与rpc差不多，但是会把完整的http头封装向下传送，不限制请求方法
3 | http或proxy | 以反向代理的方式使用**API**，相当于把普通的web应用部署在**API**之后，让外界像调api接口一样调用web服务
4 | web | 与http差不多，但是支持websocket
5 | event | 代理event事件服务类型的请求
6 | meta* | 默认值，元数据，通过在代码中的`Endpoint`配置选择使用上述中的某一个处理器，默认RPC

- `rpc`或`api`模式同样可以使用`Endpoint`定义路由。

<img src="/docs/Micro API/img/micro-router.png" width="75%">

- router过程
	- endpoint
		- 自定义路由
	- resolver
		- 路径规则

**Resolver**

`rpc`需要服务名称`go.micro.api.greeter` + 方法名`Greeter.Hello`

请求路径    |    后台服务    |    接口方法
----    |    ----    |    ----
/foo/bar    |    go.micro.api.foo    |    Foo.Bar
/foo/bar/baz    |    go.micro.api.foo    |    Bar.Baz
/foo/bar/baz/cat    |    go.micro.api.foo.bar    |    Baz.Cat
/v1/foo/bar    |    go.micro.api.v1.foo    |    Foo.Bar
/v1/foo/bar/baz    |    go.micro.api.v1.foo    |    Bar.Baz
/v2/foo/bar    |    go.micro.api.v2.foo    |    Foo.Bar
/v2/foo/bar/baz    |    go.micro.api.v2.foo    |    Bar.Baz


`proxy`只需要服务名称，用于服务发现，将http请求转发到对应的服务

请求路径    |    服务    |    后台服务路径
---    |    ---    |    ---
/foo    |   go.micro.api.foo	|   /foo
/foo/bar	|   go.micro.api.foo	|   /foo/bar
/greeter    |    go.micro.api.greeter    |    /greeter
/greeter/:name    |    go.micro.api.greeter    |    /greeter/:name


#### 演示: API Handler

<details>
  <summary> 默认网关 </summary>
  
```bash
micro --registry=etcd api

# micro-in-cn/tutorials/examples/basic-practices/micro-api/meta
go run meta.go --registry=etcd
```

```bash
curl -XGET "http://localhost:8080/example?name=john"
curl -XPOST -H 'Content-Type: application/json' -d '{"name": "john"}' "http://localhost:8080/example"
```

</details>

**--handler=api**
```bash
micro --registry=etcd api --handler=api

# micro-in-cn/tutorials/examples/basic-practices/micro-api/api
go run api.go --registry=etcd
```

```bash
curl -XGET "http://localhost:8080/example/call?name=john"
curl -XPOST -H 'Content-Type: application/json' -d '{data:123}' http://localhost:8080/example/foo/bar
```


**--handler=proxy**
```bash
micro --registry=etcd api --handler=proxy

# micro-in-cn/tutorials/examples/basic-practices/micro-api/proxy
go run proxy.go --registry=etcd
```

```bash
curl -XGET "http://localhost:8080/example/call?name=john"
curl -H 'Content-Type: application/json' -d '{"name": "john"}' http://localhost:8080/example/foo/bar
```

## 自定义网关

- import
	- 适合简单定制，如go-micro组件、增加插件，参考[micro-in-cn/starter-kit/gateway](https://github.com/micro-in-cn/starter-kit/tree/master/gateway)
- fork
	- 需要修改网关源码
	
> 不管需求大小都建议在项目中自己编译`micro`工具，确保开发、生产等环境一致

> 以下示例在[Micro API/example](/docs/Micro%20API/example)：`main_01.go`、`main_02.go`
### go-micro组件

- registry
	- consul
	- kubernetes
- transport
	- tcp
	- grpc

<details>
  <summary> 自定义组件 </summary>
  
```go
package main

import (
	"github.com/micro/micro/cmd"

	// go-micro plugins
	_ "github.com/micro/go-plugins/registry/consul"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	_ "github.com/micro/go-plugins/transport/grpc"
	_ "github.com/micro/go-plugins/transport/tcp"
)

func main() {
	cmd.Init()
}
```

</details>

#### 演示: Registry&Transport

```bash
 go build -o bin/micro_01 main_01.go
```

```go
./bin/micro_01 --registry=consul --transport=tcp api
```

```bash
#https://github.com/micro-in-cn/tutorials/tree/master/examples/basic-practices/micro-api/meta
go run meta.go --registry=consul --transport=tcp
```

```bash
curl -XGET "http://localhost:8080/example?name=john"
curl -XPOST -H 'Content-Type: application/json' -d '{"name": "john"}' "http://localhost:8080/example"
```

### plugin

plugin是使用网关的关键，类似各种web框架的中间件，通过`HTTP`请求上下文的前置、后置处理实现拦截、装饰等各种场景的需求，如：

- 跨域
- 认证鉴权
- 监控
- 限流
- 链路追踪
- 日志
- 流量染色
_ ……

**plugin接口**
```go
type Plugin interface {
	// Global Flags
	Flags() []cli.Flag
	// Sub-commands
	Commands() []cli.Command
	// Handle is the middleware handler for HTTP requests. We pass in
	// the existing handler so it can be wrapped to create a call chain.
	Handler() Handler
	// Init called when command line args are parsed.
	// The initialised cli.Context is passed in.
	Init(*cli.Context) error
	// Name of the plugin
	String() string
}
```

<details>
  <summary> Wrap ResponseWriter </summary>
  
- StatusCode
- In/Out Bytes
- ……
  
```go
type WrapWriter struct {
	StatusCode  int
	wroteHeader bool

	http.ResponseWriter
}

func (ww *WrapWriter) WriteHeader(statusCode int) {
	ww.wroteHeader = true
	ww.StatusCode = statusCode
	ww.ResponseWriter.WriteHeader(statusCode)
}
```
  
</details>

#### 演示: Metrics

```bash
go build -o bin/micro_02 main_02.go
```

```go
./bin/micro_02 --registry=consul --transport=tcp api
```

访问服务:http://localhost:8080/metrics

做些访问数据，再看`metrics`结果
```bash
curl -XGET "http://localhost:8080/example?name=john"
curl -XGET "http://localhost:8080/example?name=john"
curl -XGET "http://localhost:8080/example?name=john"
……

# 或
hey -z 10s -c 1 "http://localhost:8080/example?name=john"
```
