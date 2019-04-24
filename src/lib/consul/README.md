## consul使用说明

包含服务提供方和服务消费方两种角色

### 服务提供方

```golang
// 参数说明
// 1. consul 服务IP
// 2. consul端口
// 3. 服务IP
// 4. 服务端口
// 5. 服务名
// 6. Check路径，可以配置成完整http路径，（可选，不填时将以IP/PORT做为服务检查方式）
consul.ConsulRegistService("127.0.0.1", 8500, "127.0.0.1", 10000, "httpServer", "/check")

// 兼容之前的调用方式
consul.ConsulRegistGrpc(cfg conf.Conf)

// 新调用方式，5个参数，consul.ConsulRegisteService的别名，参数含义一样
consul.ConsulRegistGrpc5("127.0.0.1", 8500, "127.0.0.1", 8000, "grpc_test")

```
简单调用方法

```go
consul.Init(conf.Conf) // 参数为conf.Conf类型实例
consul.RegisterService(service_ip, service_port, service_name)
consul.RegisterGrpc(service_ip, service_port, service_name)
consul.RegisterGrpcDefault() // 无参数,使用配置文件中的默认配置
```


### 服务消费方

```golang
// 开启Consul自动发现服务
// consul_addr consul服务的地址，如 127.0.0.1:8500
// service_name 限定服务名
// interval 发现服务间隔
consul.DoDiscover(consul_addr, service_name, 5)

// 开启Watch模式，与DoDiscover互斥
// consul_addr consul服务的IP端口
// service_name 想监控的service_name列表,如为空则检控所有service
consul.DoWatch(consul_addr, service_name...)

// gRPC
// 参数与返回与grpc.Dial相同
consul.DoGrpcRequest(service_name, opts)

// HTTP
// 根据service_name生成正确url
// service_name 服务名
// path 路径，可以为完整路径和绝对路径
// 返回： *url.URL, error
consul.GenerateURL(service_name, path)

// 获取http.Request
// service_name 服务名
// path 完整路径或绝对路径
// body io.Reader
consul.NewPostRequest(service_name, path, body)
consul.NewGetRequest(service_name, path, body)

// 以下接口返回*http.Response, error
// GET请求
// service_name 服务名
// path 完整路径或绝对路径
consul.HttpGet(service_name, path)

// POST请求
// service_name 服务名
// path 完整路径或绝对路径
// contentType 同http.Post
// body io.Reader 同http.Post
consul.HttpGet(service_name, path)

// POST请求
// service_name 服务名
// path 完整路径或绝对路径
// data url.Values 同http.PostForm
consol.HttpPostForm(service_name, path, data)
```

简单调用方法

```go
consul.Init(conf) // conf.Conf 类型实例
consul.Watch(service_name...) // 需监控的服务名参数列表
consul.Discover(service_name)  // 需监控的服务名
```
