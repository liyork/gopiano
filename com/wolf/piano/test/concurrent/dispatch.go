package concurrent

// 设计一个分发器，当有数据进入分发器后，需要将数据分发到多个处理器处理，每个处理器可以想象为一个协程，处理器在没有数据的时候要阻塞
// 第一版本：
// chan1 --> Dispatcher --> cha1n2 --> Processor
//				        --> chan3 --> Processor
// 这种设计的最大的问题在于多个Processor处理时长不同会造成木桶效应，多个Processor会被一个Processor拖累(即Dispatcher协成会被最长Processor卡住)。
// 可以给chan加缓冲，试问缓冲设计多大合适呢？如果真存在一个处理非常慢的Processor多大的缓冲都无济于事。

// 第二版本：
// chan1 --> Dispatcher --> chan2 --> buffer --> chan3 --> Processor
//                                           --> chan3 --> Processor
// 用一个专用的协程实现从chan2读取数据放到缓冲中(不会阻塞dispatcher协成)，然后再从缓冲中读取数据放到chan3中，chan2和chan3都是无缓冲。
// 这样分发器不会由于任何一个Processor慢被拖累，同时缓冲Buffer可以设计成弹性的Buffer，不会被设定成一个固定的值。
