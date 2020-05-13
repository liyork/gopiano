#!/bin/sh
echo Building sparkdevo/href-counter:build
# 构建编译应用程序的镜像
docker build --no-cache -t sparkdevo/href-counter:build . -f Dockerfile.build
# 创建应用容器
docker create --name extract sparkdevo/href-counter:build
# 拷贝容器中编译好的应用到本地,即app
docker cp extract:/go/src/github.com/sparkdevo/href-counter/app ./app
docker rm -f extract

echo Building sparkdevo/href-counter:latest
# 构建运行应用程序的镜像
docker build --no-cache -t sparkdevo/href-counter:latest .