gRPC Name Resolver（名称解析）是 gRPC 核心功能之一，目前大部分 gRPC Name Resolver 都采用 ETCD 来实现，provider注册到etcd然后由resolver拉取

通过自定义nsResolver 将服务名解析成对应实例
init 阶段会通过 Register() 方法将 nsResolverBuilder 注册到 grpc 内部的一个全局 map 中
采用 Builder 模式在包初始化时创建并注册构造 nsResover 的 nsResolverBuilder 实例

ns.Dial() 方法使用 callerService 和 calleeService 构造服务 URI，并使用此 URI 作为参数调用 grpc.DialContext() 方法，来构造 *grpc.ClientConn 实例。
grpc.DialContext() 方法接收三个参数：ctx、target、opts，
target：就是根据我们自定义的协议名、callerService、CalleeService 生成的 URI，
比如本例中 target 参数值为 ns://my-caller-service:@my-callee-service，其中 ns 为协议名。grpc 可通过协议名查表来获取对应的 resolverBuilder。opts：是一个变长参数，表示拨号配置选项

func DialContext(ctx context.Context, target string, opts ...DialOption) (conn *ClientConn, err error) {
  // 构造 ClientConn 实例
  cc := &ClientConn{
    target:            target,
    csMgr:             &connectivityStateManager{},
    conns:             make(map[*addrConn]struct{}),
    dopts:             defaultDialOptions(),
    blockingpicker:    newPickerWrapper(),
    czData:            new(channelzData),
    firstResolveEvent: grpcsync.NewEvent(),
  }
  cc.retryThrottler.Store((*retryThrottler)(nil))
  cc.ctx, cc.cancel = context.WithCancel(context.Background())
 
 
  for _, opt := range opts {
    opt.apply(&cc.dopts)
  }
  ...
 
 
  // 如果用户指定了 timeout 超时配置，那么初始化一个带超时的 ctx
  // 如果 defer 阶段已超时，则抛出 j 错误
  if cc.dopts.timeout > 0 {
    var cancel context.CancelFunc
    ctx, cancel = context.WithTimeout(ctx, cc.dopts.timeout)
    defer cancel()
  }
  defer func() {
    select {
    case <-ctx.Done():
      conn, err = nil, ctx.Err()
    default:
    }
  }()
  ...
 
 
  // Name Resolver 核心逻辑，初始化 resolverBuilder，代码中首先会判断下用户是否指定 resolverBuilder
  // - 如果有指定 resolverBuilder，则直接使用此 resolverBuilder。
  // - 如果用户没有指定 resolverBuilder，那么 grpc 做如下操作：
  //    - 通过 parseTarget 方法解析用户传入的 target，本例中即 `ns://my-caller-service:@my-callee-service`，获取 Scheme（协议名）、authority（包含 callerService、calleeService）。
  //    - 查询指定协议对应的 resolverBuilder。
  if cc.dopts.resolverBuilder == nil {
    // 解析用户传入的 target
    cc.parsedTarget = parseTarget(cc.target)
    // 通过协议名查表获取对应的 resolverBuilder
    cc.dopts.resolverBuilder = resolver.Get(cc.parsedTarget.Scheme)
    // 如果表中没查到对应的 resolverBuilder，则使用默认协议查询对应的 resolverBuilder
    // 默认协议为 `passthrough`，它会从用户解析的 target 中直接读取 endpoint 地址
    if cc.dopts.resolverBuilder == nil {
      cc.parsedTarget = resolver.Target{
        Scheme:   resolver.GetDefaultScheme(),
        Endpoint: target,
      }
      cc.dopts.resolverBuilder = resolver.Get(cc.parsedTarget.Scheme)
    }
  } else {
    cc.parsedTarget = resolver.Target{Endpoint: target}
  }
  ...
 
 
  // 使用上面初始化的 resolverBuilder 构建 resolver
  // 初始化 resolverWrapper
  rWrapper, err := newCCResolverWrapper(cc)
  if err != nil {
    return nil, fmt.Errorf("failed to build resolver: %v", err)
  }
  cc.mu.Lock()
  cc.resolverWrapper = rWrapper
  cc.mu.Unlock()
 
 
  // 如果客户端配置了 WithBlock option，则会轮询 ClientConn 状态，如果 ClientConn 就绪，则返回 ClientConn。
  // 如果直到 ctx 超时或被 Cancel ClientConn 依然未就绪，则抛出 ctx.Err() 错误。
  if cc.dopts.block {
    for {
      s := cc.GetState()
      // 1. 如果 ClientConn 状态为 Ready 则返回此 ClientConn
      // 2. 如果 ClientConn 状态并非 Ready，且用户配置了 FailOnNonTempDialError，当前 ClientConn 状态为 TransientFailure，且 lbPicker 尝试和服务端实例建立连接时产生错误。根据错误性质做如下处理：
      //    2.1. 如果此错误是非临时性的错误，则抛出此错误
      //    2.2. 如果此错误是临时性的错误，则继续轮询 ClientConn 状态，直至 ctx 超时或被外部 Cancel
      if s == connectivity.Ready {
        break
      } else if cc.dopts.copts.FailOnNonTempDialError && s == connectivity.TransientFailure {
        if err = cc.blockingpicker.connectionError(); err != nil {
          terr, ok := err.(interface {
            Temporary() bool
          })
          if ok && !terr.Temporary() {
            return nil, err
          }
        }
      }
      if !cc.WaitForStateChange(ctx, s) {
        return nil, ctx.Err()
      }
    }
  }
 
 
  return cc, nil
}


func newCCResolverWrapper(cc *ClientConn) (*ccResolverWrapper, error) {
  ...
  ccr := &ccResolverWrapper{
    cc:     cc,
    addrCh: make(chan []resolver.Address, 1),
    scCh:   make(chan string, 1),
  }
 
 
  var err error
  // rb.Build() 调用指定 resolveBuilder 的 Build 方法，本例中会执行我们定义的 nsResolverBuilder.Builder() 方法
  ccr.resolver, err = rb.Build(cc.parsedTarget, ccr, resolver.BuildOption{DisableServiceConfig: cc.dopts.disableServiceConfig})
  if err != nil {
    return nil, err
  }
  return ccr, nil
}


对自定义 grpc Resolver 做个总结，整个工作流大概如下所示：
客户端启动时，引入自定义的 resolver 包（比如本例中我们自定义的 ns 包）
引入 ns 包，在 init() 阶段，构造自定义的 resolveBuilder，并将其注册到 grpc 内部的 resolveBuilder 表中（其实是一个全局 map，key 为协议名，比如 ns；value 为构造的 resolveBuilder，比如 nsResolverBuilder）。
客户端启动时通过自定义 Dail() 方法构造 grpc.ClientConn 单例
grpc.DialContext() 方法内部解析 URI，分析协议类型，并从 resolveBuilder 表中查找协议对应的 resolverBuilder。比如本例中我们定义的 URI 协议类型为 ns，对应的 resolverBuilder 为 nsResolverBuilder
找到指定的 resolveBuilder 后，调用 resolveBuilder 的 Build() 方法，构建自定义 resolver，同时开启协程，通过此 resolver 更新被调服务实例列表。
Dial() 方法接收主调服务名和被调服务名，并根据自定义的协议名，基于这两个参数构造服务的 URI
Dial() 方法内部使用构造的 URI，调用 grpc.DialContext() 方法对指定服务进行拨号
grpc 底层 LB 库对每个实例均创建一个 subConnection，最终根据相应的 LB 策略，选择合适的 subConnection 处理某次 RPC 请求。




一个grpc客户端首先通过resolver获取多个IP地址,这些IP地址可能是服务端的地址也可能是一个load balancer的地址,随着地址返回的还有一个service config,其中会指明load balance的策略(例如round_robin或者grpclb)
客户端初始化load balance策略.如果resolver返回的地址只要有一个是load balance的地址,客户端就会使用grpclb策略.否则的话根据service config的配置决定.如果service config没有指定load balance粗略,则默认选用pick-first
创建到每一个server address的subchannel
grpclb策略下,client连接load balancer,并向其请求服务器地址
grpc server会向load balancer汇报其负载情况
load balancer向client返回一个server list,client向每一个server都建立一个subchannel
对每一个rpc请求,load balancer决定哪一个subchannel应该被使用.
整理其流程可知有如下两种方式获取服务端地址:
resolver返回多个server地址,然后client根据round-robin或者pick-first或者random的策略选择一个server去连接.
resolver返回load-balancer的地址,load-balancer会去做server的负载检查,探活策略,然后根据负载均衡策略返回一个地址给client.



关键步骤有两步:
init()函数中进行注册
Build函数中生成一个resolver并且调用该resolver的start函数,start函数中会调用r.cc.UpdateState函数
在grpc.Dial()中会将resolver builder赋值给如下变量

cc.dopts.resolverBuilder
然后执行:
rWrapper, err := newCCResolverWrapper(cc)
...
cc.mu.Lock()
cc.resolverWrapper = rWrapper
cc.mu.Unlock()
在newCCResolverWrapper中会调用builder的build函数并赋值给rWrapper的resolver字段.rWrapper是一个ccResolverWrapper结构,所以UpdateState实际调用的是ccResolverWrapper的方法.通过追踪代码,发现UpdateState最终调用的是ClientConn结构体的updateResolverState,


resolver主要解决从域名获取地址这一步,这一步有可能直接获取到服务端地址也可能获取到balancer的地址
