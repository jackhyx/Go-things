


#### 无缓冲channel
* 无缓冲通道上的发送操作将会阻塞，直到另一个goroutine在对应的通道上执行接收操作，这时值传送完成，两个goroutine都可以继续执行。
* 相反，如果接收操作先执行，接收方goroutine将阻塞，直到另一个goroutine在同一个通道上发送一个值
* 使用无缓冲通道进行的通信导致发送和接收goroutine同步化。因此，无缓冲通道也称为同步通道。
* 当一个值在无缓冲通道上传递时，接收值后发送方goroutine才被再次唤醒
#### 单向channel
  从上边的程序，我们可以看出来，我们通过createWorker函数创建出来的通道，是用来发送数据的，所以可以将createWorker的返回值写成这样
 ``` 
 func createWorker() chan<- int {
  ......
  }
  ```
这样使用者一眼就可以看出来，这个函数返回的channel是一个发送数据的单向channel
####  只写
  chan <- int

### 只读
<- chan int

* 如果使用一个无缓冲通道，有3个goroutine向通道中发送数据，两个比较慢的goroutine将被卡住，因为在它们发送响应结果到通道的时候没有goroutine来接收。
* 这个情况叫作goroutine泄漏，它属于一个bug。不像回收变量，泄露的goroutine不会自动回收，所以确保goroutine在不再需要的时候可以自动结束
无缓冲和缓冲通道的选择、缓冲通道容量大小的选择，都会对程序的正确性产生影响。
* 无缓冲通道提供强同步保障，因为每一次发送都需要和一次对应的接收同步；对于缓冲通道，这些操作则是解耦的。
* 如果我们知道要发送的值数量的上限，通常会创建一个容量是使用上限的缓冲通道，在接收第一个值前就完成所有的发送。
* 在内存无法提供缓冲容量的情况下，可能导致程序死锁

* channel一旦close了，接收方还是能从channel中接收到数据，收到的是channel中元素类型的零值，因为我们创建的是一个chan int类，所以它的零值就是0，因此我们看到打印出来很多的0 
* 我们就需要在接收方从通道中获取数据的时候进行判断，下边对worker函数进行修改，具体如下:
```
 func worker(id int, c chan int)  {
	for  {
		n, ok := <-c //n为获取到的具体的数，ok就是，是否还有值(如果close了，就没值了)
		if !ok {
			break
		}
		fmt.Printf("worker %d, received %d\n", id, n)
	}
}

除了用上边那种判断ok的方式，还可以通过range来遍历通道，等通道中没数据了，就不会再接收了，还是对worker进行修改，具体如下：
func worker(id int, c chan int)  {
	for n := range c {
		fmt.Printf("worker %d, received %d\n", id, n)
	}
}
  
  //注意，上边都是建立在channel被close的情况，如果没有close，其它goroutine一直发，他就会一直收，直到main执行结束
```

#### 无缓冲 channel 的常见用途
* Go 语言倡导： Do not communicate by sharing memory; instead, share memory by communicating. 
* 不要通过共享内存来通信，而是通过通信来共享内存
* 多 goroutine 通信：信号 
* 基于无 buffer channel，可以实现一对一和一对多的信号传递。
```
package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)
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
```
* 一对多
// 关闭一个无 buffer channel 会让所有阻塞在这个 channel 上的 read 操作返回，基于此我们可以实现 1 对 n 的“广播”机制
```
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
```
#### 多 goroutine 同步：通过阻塞，替代锁
```
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
```
#### 多协程并发增加计数，通过 channel 写入阻塞，读取时加一
```
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
```
#### 带缓冲 channel 的常见用途
* 消息队列 
* channel 的特性符合对消息队列的要求： 
* 跨 goroutine 访问安全 
* FIFO 
* 可设置容量 
* 异步收发

* Go 支持 channel 的初衷是将它作为 Goroutine 间的通信手段，它并不是专门用于消息队列场景的。 
* 如果你的项目需要专业消息队列的功能特性，比如支持优先级、支持权重、支持离线持久化等，那么 channel 就不合适了，可以使用第三方的专业的消息队列实现。

#### 计数信号量
// 由于带 buffer channel 的特性（容量满时写入会阻塞），可以用它的容量表示同时最大并发数量。

var active = make(chan struct{}, 3) //"信号量"，最多 3 个
var jobs = make(chan int, 10)

* 使用带缓存的 channel，容量就是信号量的大小
```
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
```
* 上面的代码中，我们用 channel jobs 表示要执行的任务（这里为 8 个），然后用 channel active 表示信号量（最多三个）。
* 然后在 8 个 goroutine 里执行任务，每个任务耗时 2s。在每次执行任务前，先写入 channel 表示获取信号量；执行完后读取，表示释放信号量。 
* 由于信号量最多三个，所以同一时刻最多能有 3 个任务得以执行。
