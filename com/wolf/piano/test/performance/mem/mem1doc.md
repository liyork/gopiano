场景：如何分析Golang程序的内存使用情况？

参看[memtest1.go]

看出来，没有退出的snippet_mem进程有约830m的内存被占用
直观上来说，这个程序在test()函数执行完后，切片contaner的内存应该被释放，不应该占用830M那么大。
使用GODEBUG来分析程序的内存使用情况
添加环境变量GODEBUG='gctrace=1'来跟踪打印垃圾回收器信息
GODEBUG='gctrace=1' ./snippet_mem
设置gctrace=1会使得垃圾回收器在每次回收时汇总所回收内存的大小以及耗时， 并将这些内容汇总成单行内容打印到标准错误输出中。

分析一下本次GC的结果。
$ GODEBUG='gctrace=1' ./snippet_mem

2020/03/02 11:22:37 Start.
2020/03/02 11:22:37  ===> loop begin.
gc 1 @0.002s 5%: 0.14+0.45+0.002 ms clock, 0.29+0/0.042/0.33+0.005 ms cpu, 4->4->0 MB, 5 MB goal, 2 P
gc 2 @0.003s 4%: 0.13+3.7+0.019 ms clock, 0.27+0/0.037/2.8+0.038 ms cpu, 4->4->2 MB, 5 MB goal, 2 P
gc 3 @0.008s 3%: 0.002+1.1+0.001 ms clock, 0.005+0/0.083/1.0+0.003 ms cpu, 6->6->2 MB, 7 MB goal, 2 P
gc 4 @0.010s 3%: 0.003+0.99+0.002 ms clock, 0.006+0/0.041/0.82+0.004 ms cpu, 5->5->2 MB, 6 MB goal, 2 P
gc 5 @0.011s 4%: 0.079+0.80+0.003 ms clock, 0.15+0/0.046/0.51+0.006 ms cpu, 6->6->3 MB, 7 MB goal, 2 P
gc 6 @0.013s 4%: 0.15+3.7+0.002 ms clock, 0.31+0/0.061/3.3+0.005 ms cpu, 8->8->8 MB, 9 MB goal, 2 P
gc 7 @0.019s 3%: 0.004+2.5+0.005 ms clock, 0.008+0/0.051/2.1+0.010 ms cpu, 20->20->6 MB, 21 MB goal, 2 P
gc 8 @0.023s 5%: 0.014+3.7+0.002 ms clock, 0.029+0.040/1.2/0+0.005 ms cpu, 15->15->8 MB, 16 MB goal, 2 P
gc 9 @0.031s 4%: 0.003+1.6+0.001 ms clock, 0.007+0.094/0/0+0.003 ms cpu, 19->19->10 MB, 20 MB goal, 2 P
gc 10 @0.034s 3%: 0.006+5.2+0.004 ms clock, 0.013+0/0.045/5.0+0.008 ms cpu, 24->24->13 MB, 25 MB goal, 2 P
gc 11 @0.040s 3%: 0.12+2.6+0.002 ms clock, 0.24+0/0.043/2.5+0.004 ms cpu, 30->30->16 MB, 31 MB goal, 2 P
gc 12 @0.043s 3%: 0.11+4.4+0.002 ms clock, 0.23+0/0.044/4.1+0.005 ms cpu, 38->38->21 MB, 39 MB goal, 2 P
gc 13 @0.049s 3%: 0.008+10+0.040 ms clock, 0.017+0/0.045/10+0.080 ms cpu, 47->47->47 MB, 48 MB goal, 2 P
gc 14 @0.070s 2%: 0.004+12+0.002 ms clock, 0.008+0/0.062/12+0.005 ms cpu, 122->122->41 MB, 123 MB goal, 2 P
gc 15 @0.084s 2%: 0.11+11+0.038 ms clock, 0.22+0/0.064/3.9+0.076 ms cpu, 93->93->93 MB, 94 MB goal, 2 P
gc 16 @0.122s 1%: 0.005+25+0.010 ms clock, 0.011+0/0.12/24+0.021 ms cpu, 238->238->80 MB, 239 MB goal, 2 P
gc 17 @0.149s 1%: 0.004+36+0.003 ms clock, 0.009+0/0.051/36+0.006 ms cpu, 181->181->101 MB, 182 MB goal, 2 P
gc 18 @0.187s 1%: 0.12+19+0.004 ms clock, 0.25+0/0.049/19+0.008 ms cpu, 227->227->126 MB, 228 MB goal, 2 P
gc 19 @0.207s 1%: 0.096+27+0.004 ms clock, 0.19+0/0.077/0.73+0.009 ms cpu, 284->284->284 MB, 285 MB goal, 2 P
gc 20 @0.287s 0%: 0.005+944+0.040 ms clock, 0.011+0/0.048/1.3+0.081 ms cpu, 728->728->444 MB, 729 MB goal, 2 P
2020/03/02 11:22:38  ===> loop end.
2020/03/02 11:22:38 force gc.
gc 21 @1.236s 0%: 0.004+0.099+0.001 ms clock, 0.008+0/0.018/0.071+0.003 ms cpu, 444->444->0 MB, 888 MB goal, 2 P (forced)
2020/03/02 11:22:38 Done.
GC forced
gc 22 @122.455s 0%: 0.010+0.15+0.003 ms clock, 0.021+0/0.025/0.093+0.007 ms cpu, 0->0->0 MB, 4 MB goal, 2 P
GC forced
gc 23 @242.543s 0%: 0.007+0.075+0.002 ms clock, 0.014+0/0.022/0.085+0.004 ms cpu, 0->0->0 MB, 4 MB goal, 2 P
GC forced
gc 24 @362.545s 0%: 0.018+0.19+0.006 ms clock, 0.037+0/0.055/0.15+0.013 ms cpu, 0->0->0 MB, 4 MB goal, 2 P
GC forced
gc 25 @482.548s 0%: 0.012+0.25+0.005 ms clock, 0.025+0/0.025/0.11+0.010 ms cpu, 0->0->0 MB, 4 MB goal, 2 P
GC forced
gc 26 @602.551s 0%: 0.009+0.10+0.003 ms clock, 0.018+0/0.021/0.075+0.006 ms cpu, 0->0->0 MB, 4 MB goal, 2 P
GC forced
gc 27 @722.554s 0%: 0.012+0.30+0.005 ms clock, 0.025+0/0.15/0.22+0.011 ms cpu, 0->0->0 MB, 4 MB goal, 2 P
GC forced
gc 28 @842.556s 0%: 0.027+0.18+0.003 ms clock, 0.054+0/0.11/0.14+0.006 ms cpu, 0->0->0 MB, 4 MB goal, 2 P
...

分析
先看在`test()`函数执行完后
立即打印的`gc 21`那行的信息。`444->444->0 MB, 888 MB goal`表示垃圾回收器已经把444M的内存标记为非活跃的内存。
再看下一个记录gc 22。0->0->0 MB, 4 MB goal表示垃圾回收器中的全局堆内存大小由888M下降为4M。

结论
1、在test()函数执行完后，demo程序中的切片容器所申请的堆空间都被垃圾回收器回收了。
2、如果此时在top指令查询内存的时候，如果依然是800+MB，说明垃圾回收器回收了应用层的内存后，（可能）并不会立即将内存归还给系统。


demo程序之后会每隔一段时间打印一些gc信息，汇总如下：
GC forced
gc 26 @120.562s 0%: 0.008+0.18+0.005 ms clock, 0.016+0/0.051/0.10+0.010 ms cpu, 0->0->0 MB, 8 MB goal, 2 P
scvg0: inuse: 0, idle: 959, sys: 959, released: 447, consumed: 512 (MB)
GC forced
gc 27 @240.562s 0%: 0.005+0.19+0.005 ms clock, 0.010+0/0.063/0.13+0.010 ms cpu, 0->0->0 MB, 4 MB goal, 2 P
GC forced
scvg1: 512 MB released
scvg1: inuse: 0, idle: 959, sys: 959, released: 959, consumed: 0 (MB)
gc 28 @360.564s 0%: 0.007+0.099+0.004 ms clock, 0.014+0/0.036/0.13+0.008 ms cpu, 0->0->0 MB, 4 MB goal, 2 P
GC forced
gc 29 @480.565s 0%: 0.006+0.30+0.005 ms clock, 0.013+0/0.048/0.12+0.010 ms cpu, 0->0->0 MB, 4 MB goal, 2 P
scvg2: 0 MB released
scvg2: inuse: 0, idle: 959, sys: 959, released: 959, consumed: 0 (MB)
GC forced
gc 30 @600.566s 0%: 0.004+0.11+0.005 ms clock, 0.009+0/0.045/0.15+0.010 ms cpu, 0->0->0 MB, 4 MB goal, 2 P
scvg3: inuse: 0, idle: 959, sys: 959, released: 959, consumed: 0 (MB)
GC forced
gc 31 @720.566s 0%: 0.004+0.081+0.004 ms clock, 0.009+0/0.024/0.10+0.008 ms cpu, 0->0->0 MB, 4 MB goal, 2 P
GC forced
gc 32 @840.567s 0%: 0.006+0.12+0.005 ms clock, 0.012+0/0.039/0.17+0.010 ms cpu, 0->0->0 MB, 4 MB goal, 2 P
scvg4: inuse: 0, idle: 959, sys: 959, released: 959, consumed: 0 (MB)

接下来看scvg相关的信息。该信息在demo程序每运行一段时间后打印一次。
scvg0时consumed(从系统申请的大小)为512M。此时内存还没有归还给系统。
scvg1时consumed为0，并且scvg1的released(归还给系统的大小)=(scvg0 released + scvg0 consumed)。此时内存已归还给系统。
我们通过top命令查看，内存占用下降为38M。
之后打印的scvg信息不再有变化。
结论：垃圾回收器在一段时间后，（可能）会将回收的内存归还给系统。