$ go run demo2.go 

2019/04/06 17:37:52 start.
2019/04/06 17:37:52 Alloc:49328(bytes) HeapIdle:66494464(bytes) HeapReleased:0(bytes)
2019/04/06 17:37:52 > loop.
2019/04/06 17:37:52 Alloc:238510080(bytes) HeapIdle:364863488(bytes) HeapReleased:334856192(bytes)
2019/04/06 17:37:52 < loop.
2019/04/06 17:37:52 Alloc:207053496(bytes) HeapIdle:664731648(bytes) HeapReleased:396263424(bytes)
2019/04/06 17:37:52 force gc.
2019/04/06 17:37:52 done.
2019/04/06 17:37:52 Alloc:49864(bytes) HeapIdle:871768064(bytes) HeapReleased:396255232(bytes)
2019/04/06 17:37:52 Alloc:51056(bytes) HeapIdle:871727104(bytes) HeapReleased:396222464(bytes)
// ... 省略部分日志
2019/04/06 17:42:32 Alloc:52304(bytes) HeapIdle:871718912(bytes) HeapReleased:396214272(bytes)
2019/04/06 17:42:42 Alloc:52416(bytes) HeapIdle:871718912(bytes) HeapReleased:396214272(bytes)
2019/04/06 17:42:52 Alloc:52528(bytes) HeapIdle:871718912(bytes) HeapReleased:603217920(bytes)
2019/04/06 17:43:02 Alloc:52640(bytes) HeapIdle:871718912(bytes) HeapReleased:871653376(bytes)
2019/04/06 17:43:12 Alloc:52752(bytes) HeapIdle:871718912(bytes) HeapReleased:871653376(bytes)
2019/04/06 17:43:22 Alloc:52864(bytes) HeapIdle:871718912(bytes) HeapReleased:871653376(bytes)
可以看到，打印done.之后那条trace信息，Alloc已经下降，即内存已被垃圾回收器回收。
在2019/04/06 17:42:52和2019/04/06 17:43:02的两条trace信息中，HeapReleased开始上升，即垃圾回收器把内存归还给系统。
距离打印done.时有5分钟时间间隔。

总结
golang的垃圾回收器在回收了应用层的内存后，有可能并不会立即将回收的内存归还给操作系统。
观察应用层代码使用的内存大小，可以观察Alloc字段。
观察程序从系统申请的内存以及归还给系统的情况，可以观察HeapIdle和HeapReleased字段。
