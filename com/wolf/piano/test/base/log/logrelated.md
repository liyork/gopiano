go profiling
zap 很好地解决了日志组件的低性能的问题
日志属于 io 密集型的组件，这类组件如何做到高性能低成本，这也将直接影响到服务成本。
Opentracing
Goroutine 是一个轻量级线程管理器，它用于完成一个 “简单的” 任务。因此它不应该负责日志。它可能导致并发问题，因为在每个 goroutine 中使用 log.New() 会重复接口，所有日志器会并发尝试访问同一个 io.Writer。
为了限制对性能的影响以及避免并发调用 io.Writer，库通常使用一个特定的 goroutine 用于日志输出。
使用异步库
logmatic.io