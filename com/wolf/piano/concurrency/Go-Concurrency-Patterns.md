Go Concurrency Patterns

Concurrency features in Go
Questions:
Why is concurrency supported?
What is concurrency, anyway?
Where does the idea come from?
What is it good for?
How do I use it?

Concurrency is not parallelism, although it enables parallelism.

What is a goroutine? It's an independently executing function, launched by a go statement.
It has its own call stack, which grows and shrinks as required.
It's very cheap. It's practical to have thousands, even hundreds of thousands of goroutines.
It's not a thread.
There might be only one thread in a program with thousands of goroutines.
Instead, goroutines are multiplexed dynamically onto threads as needed to keep all the goroutines running.


A channel in Go provides a connection between two goroutines, allowing them to communicate.

When the main function executes <–c, it will wait for a value to be sent.
Similarly, when the boring function executes c <– value, it waits for a receiver to be ready.
A sender and receiver must both be ready to play their part in the communication. Otherwise we wait until they are.
Thus channels both communicate and synchronize.


Don't communicate by sharing memory, share memory by communicating.

Patterns
Channels are first-class values, just like strings or integers.


The select statement provides another way to handle multiple channels.
It's like a switch, but each case is a communication:
- All channels are evaluated.
- Selection blocks until one communication can proceed, which then does.
- If multiple can proceed, select chooses pseudo-randomly.
- A default clause, if present, executes immediately if no channel is ready.



Don't overdo it
They're fun to play with, but don't overuse these ideas.
Goroutines and channels are big ideas. They're tools for program construction.
But sometimes all you need is a reference counter.
Go has "sync" and "sync/atomic" packages that provide mutexes, condition variables, etc. They provide tools for smaller problems.
Often, these things will work together to solve a bigger problem.
Always use the right tool for the job.


In programming, concurrency is the composition of independently executing processes, 
while parallelism is the simultaneous execution of (possibly related) computations. 
Concurrency is about dealing with lots of things at once. Parallelism is about doing lots of things at once.

