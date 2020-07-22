package unaryrpc

import "google.golang.org/grpc/resolver"

// “A name resolver takes a service name and returns a list of IPs of the backends.”
type exampleResolverBuilder struct{}

func (*exampleResolverBuilder) Build(target resolver.Target,
	cc resolver.ClientConn,
	opts resolver.BuildOptions) (resolver.Resolver, error) {

	r := &exampleResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			"exampleServiceName": []string{"addr1", "addr2"},
		},
	}
	r.start()
	return r, nil
}
func (*exampleResolverBuilder) Scheme() string { return "exampleScheme" }

type exampleResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *exampleResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*exampleResolver) Close()                                  {}

func init() {
	resolver.Register(&exampleResolverBuilder{})
}
