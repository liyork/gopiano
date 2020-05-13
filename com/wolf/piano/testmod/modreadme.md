不再需要gopath以及src等目录，不再需要设定GOPATH变量,使用默认的gopath即可，即go env中的$HOME/go
cd gopiano
go mod init github.com/liyork/gopiano
go run/build com/wolf/piano/testmod/start.go  --自动相关依赖