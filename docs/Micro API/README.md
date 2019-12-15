# Micro API

## 微服务架构-网关

## micro网关

### 网关服务发现

- commands
	- `micro api`
	- `micro web`

- namespace
	- `go.micro.api.*`
	- `go.micro.web.*`

### API

```bash
micro api
```

```bash
curl localhost:8080                                                                                 
{"version": "1.18.0"}
```

> 教程演示用例[micro-in-cn/tutorials/examples/basic-practices/micro-api](https://github.com/micro-in-cn/tutorials/tree/master/examples/basic-practices/micro-api)

```bash
cd meta
go run meta.go
```

```bash
curl -XGET "http://localhost:8080/example?name=john"
curl -XPOST -H 'Content-Type: application/json' -d '{"name": "john"}' "http://localhost:8080/example"
```

```bash
micro api -h
--address value    Set the api address e.g 0.0.0.0:8080 [$MICRO_API_ADDRESS]
--handler value    Specify the request handler to be used for mapping HTTP requests to services; {api, event, http, rpc} [$MICRO_API_HANDLER]
--namespace value  Set the namespace used by the API e.g. com.example.api [$MICRO_API_NAMESPACE]
--resolver value   Set the hostname resolver used by the API {host, path, grpc} [$MICRO_API_RESOLVER]
--enable_rpc       Enable call the backend directly via /rpc [$MICRO_API_ENABLE_RPC]
```

**address & namespace**
```bash
micro api --address=:9080 --namespace=com.hbchen.api
go run meta.go --server_name=com.hbchen.api.example
```

- router过程
	- endpoint
		- 自定义路由
	- resolver
		- request -> endpoint Name & Method（）
	- registry
		- services

- handler
	- meta
	- rpc
	- api
	- proxy / http
	- web

### Web

> 虽然定义的不是网关，但可以作为网关来用

```bash
micro web
```

## 自定义网关

****
- import
	- 适合简单定制，如go-micro组件、增加插件，参考[micro-in-cn/starter-kit/gateway](https://github.com/micro-in-cn/starter-kit/tree/master/gateway)
- fork
	- 需要修改网关源码

### go-micro组件

- registry
	- kubernetes
	- consul
- transport
	- tcp

**Global Options**
```bash
--registry 
Registry for discovery. etcd, mdns
--registry_address
Comma-separated list of registry addresses
--server_name 
Name of the server. go.micro.srv.example
--transport 
Transport mechanism used; http
--enable_stats
Enable stats
```

### plugin

- Flags
- Commands
- Handler
- Init

- 跨域
- 认证鉴权
- 监控
- 限流
- 链路追踪
- 日志
- 流量染色

#### Metrics

- StatusCode
- In/Out Bytes
