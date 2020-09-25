package ns

import (
	"context"
	"strings"
	"time"

	_ "google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

const (
	// 同步周期
	syncNSInterval = 1 * time.Second
)

// 实现resolver.Resolver接口
type nsResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
	ctx    context.Context
	cancel context.CancelFunc
}

// 定时更新服务的实例的可用地址
func (r *nsResolver) watcher() {
	r.updateCC()
	ticker := time.NewTicker(syncNSInterval)
	for {
		select {
		case <-r.ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			r.updateCC()
		}
	}
}

// updateCC 更新 resolver.Resolver.ClientConn 配置
func (r *nsResolver) updateCC() {
	instances, err := r.getInstances(r.target)
	if err != nil || len(instances.calleeIns) == 0 {
		println("[mis] error retrieving instances from Mis", r.target, err)
		return
	}
	//convert
	var newAddrs []resolver.Address
	for _, k := range instances.calleeIns {
		newAddrs = append(newAddrs, resolver.Address{Addr: k})
	}

	// 更新实例列表
	// grpc 底层 LB 组件对每个服务端实例创建一个 subConnection。并根据设定的 LB 策略，选择合适的 subConnection 处理某次 RPC 请求。
	r.cc.UpdateState(resolver.State{Addresses: newAddrs})
}

// 实现了 resolver.Resolver.ResolveNow 方法
func (*nsResolver) ResolveNow(o resolver.ResolveNowOptions) {}

// 实现了 resolver.Resolver.Close 方法
func (r *nsResolver) Close() {
	r.cancel()
}

type instances struct {
	callerService string
	calleeService string
	calleeIns     []string
}

// getInstances 获取指定服务所有可用的实例列表，通过从注册中心获取
func (r *nsResolver) getInstances(target resolver.Target) (s *instances, e error) {
	auths := strings.Split(target.Authority, "@")
	return &instances{
		callerService: auths[0],
		calleeService: target.Endpoint,
		calleeIns:     []string{},
	}, nil
}
