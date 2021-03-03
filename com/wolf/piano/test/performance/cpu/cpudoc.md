3. 场景：如何分析Golang程序的CPU性能情况？
性能分析注意事项
a.性能分析必须在一个
可重复的、稳定的环境中来进行。
b.机器必须闲置
不要在共享硬件上进行性能分析;
不要在性能分析期间，在同一个机器上去浏览网页
c.注意省电模式和过热保护，如果突然进入这些模式，会导致分析数据严重不准确
d.不要使用虚拟机、共享的云主机，太多干扰因素，分析数据会很不一致
e.不要在 macOS 10.11 及以前的版本运行性能分析，有 bug，之后的版本修复了
f.关闭电源管理、过热管理
g.绝不要升级，以保证测试的一致性，以及具有可比性
一定要在多个环境中，执行多次，以取得可参考的、具有相对一致性的测试结果。



几种方式来查看进程的cpu性能情况.
A. Web界面查看
http://127.0.0.1:10000/debug/pprof/

通过pprof查看包括(阻塞信息、cpu信息、内存堆信息、锁信息、goroutine信息等等), 我们这里关心的cpu的性能的profile信息.
有关profile下面的英文解释大致如下:
“CPU配置文件。您可以在秒GET参数中指定持续时间。获取概要文件后，请使用go tool pprof命令调查概要文件。”
要是想得到cpu性能,就是要获取到当前进程的profile文件


B. 使用pprof工具查看
pprof 的格式如下
go tool pprof [binary] [profile]
binary: 必须指向生成这个性能分析数据的那个二进制可执行文件；
profile: 必须是该二进制可执行文件所生成的性能分析数据文件。
binary 和 profile 必须严格匹配。

go tool pprof ./demo4 profile
(pprof) top

几列数据,需要说明一下
flat：当前函数占用CPU的耗时
flat：:当前函数占用CPU的耗时百分比
sun%：函数占用CPU的耗时累计百分比
cum：当前函数加上调用当前函数的函数占用CPU的总耗时
cum%：当前函数加上调用当前函数的函数占用CPU的总耗时百分比
最后一列：函数名称
通过结果我们可以看出, 该程序的大部分cpu性能消耗在 main.getSoneBytes()方法中,其中math/rand取随机数消耗比较大.


或者
通过go tool pprof得到profile文件
go tool pprof http://localhost:10000/debug/pprof/profile?seconds=60
(pprof) top


D.可视化查看
go tool pprof ./demo4 profile
(pprof) web
graphviz工具,需要安装一下

Mac安装
brew install graphviz  
将graphviz安装目录下的bin文件夹添加到Path环境变量中。 在终端输入dot -version查看是否安装成功。
得到一个svg的可视化文件在/tmp路径下
这样我们就能比较清晰的看到函数之间的调用关系,方块越大的表示cpu的占用越大.


目前最后使用火焰图方式：
go get -u github.com/google/pprof 
pprof -http 127.0.0.1:9090 http://targetip/debug/pprof/profile\?seconds\=120
