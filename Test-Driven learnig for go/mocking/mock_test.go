package main

import (
	"bytes"
	"reflect"
	"testing"
)

/*
func Countdown(out *bytes.Buffer) {
	fmt.Fprint(out, "3")
} */ // 我们正在使用 fmt.Fprint 传入一个 io.Writer（例如 *bytes.Buffer）并发送一个 string。

func TestCountdown(t *testing.T) {

	t.Run("prints 3 to Go!", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		Countdown(buffer, &CountdownOperationsSpy{})

		got := buffer.String()
		want := `3
2
1
Go!`

		if got != want {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})

	t.Run("sleep after every print", func(t *testing.T) {
		spySleepPrinter := &CountdownOperationsSpy{}
		Countdown(spySleepPrinter, spySleepPrinter)

		want := []string{
			sleep,
			write,
			sleep,
			write,
			sleep,
			write,
			sleep,
			write,
		}

		if !reflect.DeepEqual(want, spySleepPrinter.Calls) {
			t.Errorf("wanted calls %v got %v", want, spySleepPrinter.Calls)
		}
	})
}

// 我们现在在 Sleeper 上有两个测试监视器;
// 所以我们现在可以重构我们的测试，一个测试被打印的内容，另一个是确保我们在打印时间 sleep
//  更新测试以注入对我们监视器的依赖，并断言 sleep 被调用了 4 次。
