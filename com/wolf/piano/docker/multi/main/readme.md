在应用了容器技术的软件开发过程中，控制容器镜像的大小可是一件费时费力的事情。如果我们构建的镜像既是编译软件的环境，又是软件最终的运行环境，这是很难控制镜像大小的。所以常见的配置模式为：分别为软件的编译环境和运行环境提供不同的容器镜像。比如为编译环境提供一个 Dockerfile.build，用它构建的镜像包含了编译软件需要的所有内容，比如代码、SDK、工具等等。同时为软件的运行环境提供另外一个单独的 Dockerfile，它从 Dockerfile.build 中获得编译好的软件，用它构建的镜像只包含运行软件所必须的内容。这种情况被称为构造者模式(builder pattern)，本文将介绍如何通过 Dockerfile 中的 multi-stage 来解决构造者模式带来的问题。

Dockerfile.build 
构建编译应用程序的镜像

Dockerfile
把构建好的应用程序部署到正式镜像

build.sh
把整个构建过程整合

执行：
./build.sh

查看构建后的镜像大小:
docker images|grep href
可以观察到，用于编译应用程序的容器镜像(sparkdevo/href-counter:build)大小接近 700M，
而用于生产环境的容器镜像(sparkdevo/href-counter:latest)只有 10.3 M，这样的大小在网络间传输的效率是很高的。

检查构建的容器是否可以正常的工作：
docker run -e url=https://www.cnblogs.com/ sparkdevo/href-counter:latest

采用上面的构建过程，我们需要维护两个 Dockerfile 文件和一个脚本文件 build.sh
docker 针对这种情况提供的解决方案：multi-stage。
multi-stage 允许我们在 Dockerfile 中完成类似前面 build.sh 脚本中的功能，每个 stage 可以理解为构建一个容器镜像，后面的 
stage 可以引用前面 stage 中创建的镜像。所以我们可以使用下面单个的 Dockerfile 文件实现前面的需求：
如Dockerfile.multi
同时存在多个 FROM 指令，每个 FROM 指令代表一个 stage 的开始部分
可以把一个 stage 的产物拷贝到另一个 stage 中。本例中的第一个 stage 完成了应用程序的构建，内容和前面的 Dockerfile.build 是一样的。第二个 stage 中的 COPY 指令通过 --from=0 引用了第一个 stage ，并把应用程序拷贝到了当前 stage 中
docker build --no-cache -t sparkdevo/href-counter:multi . -f Dockerfile.multi
运行：
docker run -e url=https://www.cnblogs.com/ sparkdevo/href-counter:multi
新生成的镜像有没有特别之处呢：
除了 sparkdevo/href-counter:multi 镜像，还生成了一个匿名的镜像。因此，所谓的 multi-stage 不过时多个 Dockerfile 的语法糖罢了
可以为 stage 命名的，然后就可以通过名称来引用 stage 了
把第一个 stage 使用 as 语法命名为 builder，然后在后面的 stage 中通过名称 builder 进行引用 --from=builder
只有 17.05 以及之后的版本才开始支持