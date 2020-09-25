module github.com/liyork/gopiano

go 1.13

require (
	github.com/astaxie/beego v1.12.1
	github.com/go-redis/redis v6.15.7+incompatible
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/gohouse/converter v0.0.3 // indirect
	github.com/gohouse/gorose v1.0.5
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.4.2
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/onsi/ginkgo v1.13.0 // indirect
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644 // indirect
	github.com/spf13/pflag v1.0.5
	golang.org/x/net v0.0.0-20200520004742-59133d7f0dd7
	golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.25.0
	gopkg.in/inf.v0 v0.9.1 // indirect
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
