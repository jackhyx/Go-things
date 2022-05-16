package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// 无缓冲 channel 的常见用途
//Go 语言倡导： Do not communicate by sharing memory; instead, share memory by communicating. 不要通过共享内存来通信，而是通过通信来共享内存
//多 goroutine 通信：信号
//基于无 buffer channel，可以实现一对一和一对多的信号传递。

type signal struct{}

//接收一个函数，在子 routine 里执行，然后返回一个 channel，用于主 routine 等待
func spawn(f func()) <-chan signal {
	c := make(chan signal)
	go func() {
		fmt.Println("exec f in child_routine")
		f()
		fmt.Println("f exec finished, write to channel")
		c <- signal{}
	}()
	return c
}

//测试使用无 buffer channel 实现信号
func main() {
	//模拟主 routine 等待子 routine

	worker := func() {
		fmt.Println("do some work")
		time.Sleep(1 * time.Second)
	}

	fmt.Println("start a worker...")
	c := spawn(worker)

	fmt.Println("spawn finished, read channel...")
	<-c //读取，阻塞等待

	fmt.Println("worker finished")
}

// 一对多
// 关闭一个无 buffer channel 会让所有阻塞在这个 channel 上的 read 操作返回，基于此我们可以实现 1 对 n 的“广播”机制。
var waitGroup sync.WaitGroup

type signal struct{}

func spawnGroup(f func(ind int), count int, groupSignal chan struct{}) <-chan signal {
	c := make(chan signal) //用于让主 routine 阻塞的 channel
	waitGroup.Add(count)   //等待总数

	//创建 n 个 goroutine
	for i := 0; i < count; i++ {
		go func(index int) {
			<-groupSignal //读取阻塞，等待通知执行

			//fmt.Println("exec f in child_routine, index: ", i);
			//⚠️注意上面注释的代码，这里不能直接访问 for 循环的 i，因为这个是复用的，会导致访问的值不是目标值

			fmt.Println("exec f in child_routine, index: ", index)
			f(index)
			fmt.Println(index, " exec finished, write to channel")

			waitGroup.Done()
		}(i + 1)
	}

	//创建通知主 routine 结束的 routine，不能阻塞当前函数
	go func() {
		//需要同步等待所有子 routine 执行完
		waitGroup.Wait()
		c <- signal{} //写入数据
	}()
	return c
}

func main() {
	worker := func(i int) {
		fmt.Println("do some work, index ", i)
		time.Sleep(3 * time.Second)
	}

	groupSignal := make(chan struct{})
	c := spawnGroup(worker, 5, groupSignal)

	fmt.Println("main routine: close channel")
	close(groupSignal) //通知刚创建的所有 routine

	fmt.Println("main routine: read channel...")
	<-c //阻塞在这里

	fmt.Println("main routine: all worker finished")
}

// 多 goroutine 同步：通过阻塞，替代锁

type NewCounter struct {
	c chan int
	i int
}

func CreateNewCounter() *NewCounter {
	counter := &NewCounter{
		c: make(chan int),
		i: 0,
	}

	go func() {
		for {
			counter.i++
			counter.c <- counter.i //每次加一，阻塞在这里
		}
	}()

	return counter
}

func (c *NewCounter) Increase() int {
	return <-c.c //读取到的值，是上一次加一
}

//多协程并发增加计数，通过 channel 写入阻塞，读取时加一
func main() {
	fmt.Println("\ntestCounterWithChannel ->>>")

	group := sync.WaitGroup{}
	counter := CreateNewCounter()

	for i := 0; i < 10; i++ {
		group.Add(1)

		go func(i int) {
			count := counter.Increase()
			fmt.Printf("Goroutine%d, count %d \n", i, count)
		}(i)
	}

	group.Wait()

}

// 带缓冲 channel 的常见用途
//  消息队列
//  channel 的特性符合对消息队列的要求：

// 跨 goroutine 访问安全
// FIFO
// 可设置容量
// 异步收发

// Go 支持 channel 的初衷是将它作为 Goroutine 间的通信手段，它并不是专门用于消息队列场景的。
// 如果你的项目需要专业消息队列的功能特性，比如支持优先级、支持权重、支持离线持久化等，那么 channel 就不合适了，可以使用第三方的专业的消息队列实现。

//计数信号量
// 由于带 buffer channel 的特性（容量满时写入会阻塞），可以用它的容量表示同时最大并发数量。

var active = make(chan struct{}, 3) //"信号量"，最多 3 个
var jobs = make(chan int, 10)

//使用带缓存的 channel，容量就是信号量的大小
func testSemaphoreWithBufferChannel() {

	//先写入数据，用作表示任务
	go func() {
		for i := 0; i < 9; i++ {
			jobs <- i + 1
		}
		close(jobs)
	}()

	var wg sync.WaitGroup

	for j := range jobs {
		wg.Add(1)

		//执行任务
		go func(i int) {
			//通知开始执行，当容量用完时，阻塞
			active <- struct{}{}

			//fmt.Println("exec job ", i)
			log.Printf("exec job: %d, length of active: %d \n", i, len(active))
			time.Sleep(2 * time.Second)

			//执行完，通知结束
			<-active
			wg.Done()

		}(j)
	}

	wg.Wait()
}

// 上面的代码中，我们用 channel jobs 表示要执行的任务（这里为 8 个），然后用 channel active 表示信号量（最多三个）。
// 然后在 8 个 goroutine 里执行任务，每个任务耗时 2s。在每次执行任务前，先写入 channel 表示获取信号量；执行完后读取，表示释放信号量。
// 由于信号量最多三个，所以同一时刻最多能有 3 个任务得以执行。
