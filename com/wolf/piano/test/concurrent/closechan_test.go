package concurrent

import (
	"log"
	"math/rand"
	_ "net/http/pprof"
	"strconv"
	"sync"
	"testing"
	"time"
)

// 问题：
//在不能更改channel状态的情况下，没有简单普遍的方式来检查channel是否已经关闭了
//若关闭已经关闭的channel会导致panic，所以在closer(关闭者)不知道channel是否已经关闭的情况下去关闭channel是很危险的
//发送值到已经关闭的channel会导致panic，所以如果sender(发送者)在不知道channel是否已经关闭的情况下去向channel发送值是很危险的

//没有一个内置函数可以检查一个channel是否已经关闭。如果你能确定不会向channel发送任何值，那么也确实需要一个简单的方法来检查channel是否已经关闭：
//没有一种适用的方式来检查channel是否已经关闭了。但是，就算有一个简单的 closed(chan T) bool函数来检查channel是否已经关闭，它的用处还是很有限的，就像内置的len函数用来检查缓冲channel中元素数量一样。原因就在于，已经检查过的channel的状态有可能在调用了类似的方法返回之后就修改了，因此返回来的值已经不能够反映刚才检查的channel的当前状态了。
//尽管在调用closed(ch)返回true的情况下停止向channel发送值是可以的，但是如果调用closed(ch)返回false，那么关闭channel或者继续向channel发送值就不安全了（会panic）。
//
//The Channel Closing Principle
//一个适用的原则是不要从接收端关闭channel，也不要在多个并发发送端中关闭channel
//如果sender(发送者)只是唯一的sender或者是channel最后一个活跃的sender，那么你应该在sender的goroutine关闭channel，从而通知receiver(s)(接收者们)已经没有值可以读了。维持这条原则将保证永远不会发生向一个已经关闭的channel发送值或者关闭一个已经关闭的channel。
//
//可以使用SafeSend和SafeClose就是对要关闭时添加recover
//我们应该要理解为什么Go不支持内置SafeSend和SafeClose函数，
// 原因就在于并不推荐从接收端或者多个并发发送端关闭channel。Golang甚至禁止关闭只接收（receive-only）的channel。

// todo 有个问题参见:concurrent/main/testconcurrent2chan.go

// ====== 针对各种场景，下面介绍不用使用panic/recover和sync包，纯粹是利用channel的解决方案。

// 场景1:M个receivers，一个sender，sender通过关闭data channel说“不再发送”
//最简单的场景了，就只是当sender不想再发送的时候让sender关闭data 来关闭channel：
func TestOneSenderMultiReceiver(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	// ...
	const MaxRandomNumber = 100000
	const NumReceivers = 100

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	// ...
	dataCh := make(chan int, 100)

	// the sender
	go func() {
		for {
			if value := rand.Intn(MaxRandomNumber); value == 0 {
				// the only sender can close the channel safely.
				close(dataCh)
				return
			} else {
				dataCh <- value
			}
		}
	}()

	// receivers
	for i := 0; i < NumReceivers; i++ {
		go func() {
			defer wgReceivers.Done()

			// receive values until dataCh is closed and
			// the value buffer queue of dataCh is empty.
			for value := range dataCh {
				log.Println(value)
			}
		}()
	}

	wgReceivers.Wait()
}

//场景2:一个receiver，N个sender，receiver通过关闭一个额外的signal channel说“请停止发送”
//不能让receiver关闭data channel，因为这么做将会打破channel closing principle。
// 但是我们可以让receiver关闭一个额外的signal channel来通知sender停止发送值：
// 对于额外的signal channel来说，它的sender是data channel的receiver。
// 这个额外的signal channel被它唯一的sender关闭，遵守了channel closing principle。
// 最后对于dataCh并没有读写了，没有关闭可能是go的建议
func TestMultiSenderOneReceiver(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	// ...
	const MaxRandomNumber = 100000
	const NumSenders = 1000

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(1)

	// ...
	dataCh := make(chan int, 100)
	stopCh := make(chan struct{})
	// stopCh is an additional signal channel.
	// Its sender is the receiver of channel dataCh.
	// Its reveivers are the senders of channel dataCh.

	// senders
	for i := 0; i < NumSenders; i++ {
		go func() {
			for {
				value := rand.Intn(MaxRandomNumber)

				select {
				// 收到receiver的stop则不再发送
				case <-stopCh:
					return
				case dataCh <- value:
				}
			}
		}()
	}

	// the receiver
	go func() {
		defer wgReceivers.Done()

		for value := range dataCh {
			if value == MaxRandomNumber-1 {
				// the receiver of the dataCh channel is
				// also the sender of the stopCh cahnnel.
				// It is safe to close the stop channel here.
				close(stopCh)
				return
			}

			log.Println(value)
		}
	}()

	// ...
	wgReceivers.Wait()
}

// 场景3:M个receiver，N个sender，它们当中任意一个通过通知一个moderator（仲裁者）关闭额外的signal channel来说“让我们结束游戏吧”
//我们不能让任意的receivers和senders关闭data channel，也不能让任何一个receivers通过关闭一个额外的signal channel来通知
// 所有的senders和receivers退出游戏。这么做的话会打破channel closing principle。
// 但是，我们可以引入一个moderator来关闭一个额外的signal channel。这个例子的一个技巧是怎么通知moderator去关闭额外的signal channel：
// 注意channel toStop的缓冲大小是1.这是为了避免当mederator goroutine 准备好之前第一个通知就已经发送了，导致丢失。
func TestMultiSenderMultiReceiver(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(0)

	// ...
	const MaxRandomNumber = 100000
	const NumReceivers = 10
	const NumSenders = 1000

	wgReceivers := sync.WaitGroup{}
	wgReceivers.Add(NumReceivers)

	// ...
	dataCh := make(chan int, 100)
	stopCh := make(chan struct{}) //只有读和close动作
	// stopCh is an additional signal channel.
	// Its sender is the moderator goroutine shown below.
	// Its reveivers are all senders and receivers of dataCh.
	toStop := make(chan string, 1)
	// the channel toStop is used to notify the moderator
	// to close the additional signal channel (stopCh).
	// Its senders are any senders and receivers of dataCh.
	// Its reveiver is the moderator goroutine shown below.

	var stoppedBy string

	// moderator
	go func() {
		stoppedBy = <-toStop // part of the trick used to notify the moderator
		// to close the additional signal channel.
		close(stopCh)
	}()

	// senders
	for i := 0; i < NumSenders; i++ {
		go func(id string) {
			for {
				value := rand.Intn(MaxRandomNumber)
				if value == 0 {
					// here, a trick is used to notify the moderator
					// to close the additional signal channel.
					select {
					case toStop <- "sender#" + id:
					default:
					}
					return
				}

				// the first select here is to try to exit the
				// goroutine as early as possible.
				select {
				case <-stopCh:
					return
				default:
				}

				select {
				case <-stopCh:
					return
				case dataCh <- value:
				}
			}
		}(strconv.Itoa(i))
	}

	// receivers
	for i := 0; i < NumReceivers; i++ {
		go func(id string) {
			defer wgReceivers.Done()

			for {
				// same as senders, the first select here is to
				// try to exit the goroutine as early as possible.
				select {
				case <-stopCh:
					return
				default:
				}

				select {
				case <-stopCh:
					return
				case value := <-dataCh:
					if value == MaxRandomNumber-1 {
						// the same trick is used to notify the moderator
						// to close the additional signal channel.
						select {
						case toStop <- "receiver#" + id:
						default:
						}
						return
					}

					log.Println(value)
				}
			}
		}(strconv.Itoa(i))
	}

	// ...
	wgReceivers.Wait()
	log.Println("stopped by", stoppedBy)
}

// 这里没有一种场景要求你去打破channel closing principle。如果你遇到了这种场景，请思考一下你的设计并重写你的代码。
// 用Go编程就像在创作艺术。
