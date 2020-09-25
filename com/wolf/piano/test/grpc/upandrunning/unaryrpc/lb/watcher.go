package main

import (
	"google.golang.org/grpc/naming"
	"net"
)

// watcher is the implementaion of grpc.naming.Watcher
type watcher struct {
	re            *resolver // re: Etcd Resolver
	isInitialized bool
}

// Close do nothing
func (w *watcher) Close() {
}

// Next to return the updates
func (w *watcher) Next() ([]*naming.Update, error) {
	// prefix is the etcd prefix/value to watch

	// check if is initialized
	addrs := extractAddrs()
	//if not empty, return the updates or watcher new dir
	if l := len(addrs); l != 0 {
		updates := make([]*naming.Update, l)
		for i := range addrs {
			updates[i] = &naming.Update{Op: naming.Add, Addr: addrs[i]}
		}
		return updates, nil
	}
	return nil, nil
}

var i = 1

func extractAddrs() []string {
	i++
	if i%2 == 0 {
		return nil
	}
	net.ResolveIPAddr("ip", "www.baidu.com")
	//fmt.Println(addr, addr.IP, err)
	addrs := []string{"10.0.13.55:23444", "10.0.13.55:23445"}
	return addrs
}
