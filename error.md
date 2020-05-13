1. go run 时错误
vendor/k8s.io/client-go/rest/request.go:598:31: not enough arguments in call to watch.NewStreamWatcher
	have (*versioned.Decoder)
	want (watch.Decoder, watch.Reporter)
	
由于有版本问题。
方案：
require (
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/gohouse/converter v0.0.3 // indirect
	github.com/gohouse/gorose v1.0.5
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6 // indirect
	golang.org/x/sys v0.0.0-20191228213918-04cbcbbfeed8 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.2.7 // indirect
	k8s.io/api v0.17.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/utils v0.0.0-20191218082557-f07c713de883 // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20191004102349-159aefb8556b
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20191004105649-b14e3c49469a
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191004074956-c5d2f014d689
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.3.0
)


2. 