
/*
   当counter增加时，直接return
   当counter减少时, 判断条件：counter > 0 || waiter == 0

   true时，直接return
   false（等待线程都完成且有等待者）时，statep复位为0，通过semap信号量唤醒所有等待者

 */
func (wg *WaitGroup) Add(delta int) {
	//从数组中拿到stetep（counter+waiter的组合）和semap信号量的内存地址
	statep, semap := wg.state()
	//stetep原子加操作，高位32bit是counter,实际counter+1
	state := atomic.AddUint64(statep, uint64(delta)<<32)
	//state的高位32bit，表示couter的计数值
	v := int32(state >> 32)
	//state的低位32bit，表示waiter的等待者数量
	w := uint32(state)
	// couter不能小于0
	if v < 0 {
		panic("sync: negative WaitGroup counter")
	}
	// 需要避免错误操作：Add和Wait并发操作，否则会panic
	if w != 0 && delta > 0 && v == int32(delta) {
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}
	// 如果还有等待线程未完成或者并没有等待者，直接return
	if v > 0 || w == 0 {
		return
	}
	// 需要避免错误操作：Add和Wait并发操作，否则会panic
	if *statep != state {
		panic("sync: WaitGroup misuse: Add called concurrently with Wait")
	}
	// 将statep复位为0（ counter和waiter都置为0）
	*statep = 0
	// 有多少个等待者就往semap循环发信号量（其实就是semap+1），Wait等待有一个调用	// runtime_Semacquire(semap)就是在等待这个信号量
	for ; w != 0; w-- {
		runtime_Semrelease(semap, false, 0)
	}
}
// 主线程循环对waiter原子操作+1直到成功后，然后阻塞等待semap信号量而被唤醒，最后return
func (wg *WaitGroup) Wait() {
// 从数组中拿到stetep（counter+waiter的组合）和semap信号量的内存地址
	statep, semap := wg.state()
	for {
//从内存总线中加载最新的statep值
		state := atomic.LoadUint64(statep)
//state的高位32bit，表示couter的计数值
		v := int32(state >> 32)
//state的低位32bit，表示waiter的等待者数量
		w := uint32(state)
//如果couter为0，表示当前已经没有在运行的等待线程了
		if v == 0 {
			return
		}
// CAS操作statep+1，低位属于waiter,即waiter+1
		if atomic.CompareAndSwapUint64(statep, state, state+1) {
// CAS操作成功后，阻塞等待semap信号为非零，竞争到会将semap-1，并唤醒线程
			runtime_Semacquire(semap)
			if *statep != 0 {
				panic("sync: WaitGroup is reused before previous Wait has returned")
			}
			return
		}
// CAS操作失败了，重新进入循环
	}
}
