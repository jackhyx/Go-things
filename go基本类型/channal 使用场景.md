
#### 什么场景下用channel合适呢？

* 通过全局变量加锁同步来实现通讯，并不利于多个协程对全局变量的读写操作。
* 加锁虽然可以解决goroutine对全局变量的抢占资源问题，但是影响性能，违背了原则。
#### 总结：为了解决上述的问题，我们可以引入channel，使用channel进行协程goroutine间的通信。


#### for range 从通道中取值，通道关闭时for range 退出
```
// channel练习 go  for range从chan中取值
   ch1 := make(chan int)
   ch2 := make(chan int)

   // 开启goroutine 把0-100写入到ch1通道中
   go func() {
      for i := 0; i < 100; i++ {
         ch1 <- i
      }
      close(ch1)
   }()

// 开启goroutine 从ch1中取值，值的平方赋值给 ch2
   go func() {
      for {
         i,ok := <-ch1 //通道取值后 再取值 ok = false
         if ok {
            ch2 <- i*i
         }else {
            break
         }
      }
      close(ch2)
   }()

// 主goroutine 从ch2中取值 打印输出
// for x := chan 有值取值，通道关闭时跳出goroutine
   for i :=range ch2{
      fmt.Println(i)
   }

```
#### channel升级，单通道，只读通道和只写通道
```
func counter(in chan<- int) {
   defer close(in)
   for i := 0; i < 100; i++ {
      in <- i
   }
}

func square(in chan<- int, out <-chan int) {
   defer close(in)
   for i := range out {
      in <- i * i
   }
}

func output(out <-chan int)  {
   for i:=range out{
      fmt.Println(i)
   }
}

// 改写成单向通道
func main() {
   ch1 := make(chan int)
   ch2 := make(chan int)
   go counter(ch1)
   go square(ch2, ch1)
   output(ch2)
}
```
#### goroutine work pool，可以防止goroutine暴涨或者泄露
```
//使用work pool 防止goroutine的泄露和暴涨
func worker(id int, jobs <-chan int, results chan<- int) {
   for j := range jobs {
      fmt.Printf("worker:%d start job:%d\n", id, j)
      time.Sleep(time.Second)
      fmt.Printf("worker:%d end job:%d\n", id, j)
      results <- j * 2
   }
}
      
func main() {
   jobs := make(chan int, 100)
   results := make(chan int, 100)
   // 开启3个goroutine
   for w := 1; w <= 3; w++ {
      go worker(w, jobs, results)
   }
   // 5个任务
   for j := 1; j <= 5; j++ {
      jobs <- j
   }
   close(jobs)
   // 输出结果
   for a := 1; a <= 5; a++ {
      <-results
   }
}
```
#### goroutine使用select case多路复用，满足我们同时从多个通道接收值的需求

```
//使用select语句能提高代码的可读性。
//可处理一个或多个channel的发送/接收操作。
//如果多个case同时满足，select会随机选择一个。
//对于没有case的select{}会一直等待，可用于阻塞main函数。
ch := make(chan int, 1)
go func() {
   for i := 0; i < 10; i++ {
      select {
      case x := <-ch:
         fmt.Println(x)
      case ch <- i:
      }
   }
}()

```

#### goroutine加锁 排它锁 读写锁

```
var x int64
var wg sync.WaitGroup
//添加互斥锁
var lock sync.Mutex

func main() {
   wg.Add(2)
   go add()
   go add()
   wg.Wait()
   fmt.Println(x)
}

func add() {
   for i := 0; i < 5000; i++ {
      lock.Lock() //加锁
      x = x + 1
      lock.Unlock() //解锁
   }
   wg.Done()
}
```