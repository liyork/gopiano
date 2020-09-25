$ go run demo2.go 
2020/03/02 18:21:17  ===> [Start].
2020/03/02 18:21:17  ===> Alloc:71280(bytes) HeapIdle:66633728(bytes) HeapReleased:66600960(bytes)
2020/03/02 18:21:17  ===> loop begin.
2020/03/02 18:21:18  ===> Alloc:132535744(bytes) HeapIdle:336756736(bytes) HeapReleased:155721728(bytes)
2020/03/02 18:21:38  ===> loop end.
2020/03/02 18:21:38  ===> Alloc:598300600(bytes) HeapIdle:609181696(bytes) HeapReleased:434323456(bytes)
2020/03/02 18:21:38  ===> [force gc].
2020/03/02 18:21:38  ===> [Done].
2020/03/02 18:21:38  ===> Alloc:55840(bytes) HeapIdle:1207427072(bytes) HeapReleased:434266112(bytes)
2020/03/02 18:21:38  ===> Alloc:56656(bytes) HeapIdle:1207394304(bytes) HeapReleased:434266112(bytes)
2020/03/02 18:21:48  ===> Alloc:56912(bytes) HeapIdle:1207394304(bytes) HeapReleased:1206493184(bytes)
2020/03/02 18:21:58  ===> Alloc:57488(bytes) HeapIdle:1207394304(bytes) HeapReleased:1206493184(bytes)
2020/03/02 18:22:08  ===> Alloc:57616(bytes) HeapIdle:1207394304(bytes) HeapReleased:1206493184(bytes)
c2020/03/02 18:22:18  ===> Alloc:57744(bytes) HeapIdle:1207394304(bytes) HeapReleased:1206493184(by
可以看到，打印`[Done].`之后那条trace信息，Alloc已经下降，即内存已被垃圾回收器回收。
在`2020/03/02 18:21:38`和`2020/03/02 18:21:48`的两条trace信息中，HeapReleased开始上升，即垃圾回收器把内存归还给系统。
