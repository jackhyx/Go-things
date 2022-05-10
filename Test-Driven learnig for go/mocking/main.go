package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

/*
func Countdown(out io.Writer) {
	fmt.Fprint(out, "3")
} // 使用通用接口代替。
*/

type Sleeper interface {
	Sleep()
}
type SpySleeper struct {
	Calls int
}

func (s *SpySleeper) Sleep() {
	s.Calls++
}

type ConfigurableSleeper struct {
	duration time.Duration
}

func (o *ConfigurableSleeper) Sleep() {
	time.Sleep(o.duration)
}

type CountdownOperationsSpy struct {
	Calls []string
}

func (s *CountdownOperationsSpy) Sleep() {
	s.Calls = append(s.Calls, sleep)
}

func (s *CountdownOperationsSpy) Write(p []byte) (n int, err error) {
	s.Calls = append(s.Calls, write)
	return
}

const write = "write"
const sleep = "sleep"

// 监视器（spies）是一种 mock，它可以记录依赖关系是怎样被使用的。
// 它们可以记录被传入来的参数，多少次等等。
// 在我们的例子中，我们跟踪记录了 Sleep() 被调用了多少次，这样我们就可以在测试中检查它。

const finalWord = "Go!"
const countdownStart = 3

func Countdown(out io.Writer, sleeper Sleeper) {
	for i := countdownStart; i > 0; i-- {
		sleeper.Sleep()
		fmt.Fprintln(out, i)
	}

	sleeper.Sleep()
	fmt.Fprint(out, finalWord)
}

func main() {
	sleeper := &ConfigurableSleeper{1 * time.Second}
	Countdown(os.Stdout, sleeper)
}

/*
func Countdown(out io.Writer, sleeper Sleeper) {
    for i := countdownStart; i > 0; i-- {
        sleeper.Sleep()
    }

    for i := countdownStart; i > 0; i-- {
        fmt.Fprintln(out, i)
    }

    sleeper.Sleep()
    fmt.Fprint(out, finalWord)
}
如果你运行测试，它们仍然应该通过，即使实现是错误的。
让我们再用一种新的测试来检查操作的顺序是否正确。
*/
/*
Mocking
测试可以通过，软件按预期的工作。但是我们有一些问题：
我们的测试花费了 4 秒的时间运行
每一个关于软件开发的前沿思考性文章，都强调快速反馈循环的重要性。
缓慢的测试会破坏开发人员的生产力。
想象一下，如果需求变得更复杂，将会有更多的测试。对于每一次新的 Countdown 测试，我们是否会对被添加到测试运行中 4 秒钟感到满意呢？
我们还没有测试这个函数的一个重要属性。
我们有个 Sleeping 的注入，需要抽离出来然后我们才可以在测试中控制它。
如果我们能够 mock time.Sleep，我们可以用 依赖注入 的方式去来代替「真正的」time.Sleep，然后我们可以使用断言 监视调用
*/
