package concurrency

import (
	"reflect"
	"testing"
	"time"
)

type WebsiteChecker func(string) bool
type result struct {
	string
	bool
}

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)
	resultChannel := make(chan result)
	// 因为开启 goroutine 的唯一方法就是将 go 放在函数调用前面，所以当我们想要启动 goroutine 时，我们经常使用 匿名函数。一个匿名函数文字看起来和正常函数声明一样，但没有名字;
	// 匿名函数有许多有用的特性;首先，它们可以在声明的同时执行 —— 这就是匿名函数末尾的 () 实现的。
	// 其次，它们维护对其所定义的词汇作用域的访问权 —— 在声明匿名函数时所有可用的变量也可在函数体内使用
	for _, url := range urls {
		go func(u string) {
			resultChannel <- result{u, wc(u)}
		}(url)
	}
	for i := 0; i < len(urls); i++ {
		result := <-resultChannel
		results[result.string] = result.bool
	}
	return results
}

// 这里的问题是变量 url 被重复用于 for 循环的每次迭代 —— 每次都会从 urls 获取新值。但是我们的每个 goroutine 都是 url 变量的引用 —— 它们没有自己的独立副本。所以他们 都 会写入在迭代结束时的 url —— 最后一个 url。这就是为什么我们得到的结果是最后一个 url。
// 通过给每个匿名函数一个参数 url(u)，然后用 url 作为参数调用匿名函数，我们确保 u 的值固定为循环迭代的 url 值，重新启动 goroutine。u 是 url 值的副本，因此无法更改
func mockWebsiteChecker(url string) bool {
	if url == "waat://furhurterwe.geds" {
		return false
	}
	return true
}

func TestCheckWebsites(t *testing.T) {
	websites := []string{
		"http://google.com",
		"http://blog.gypsydave5.com",
		"waat://furhurterwe.geds",
	}

	actualResults := CheckWebsites(mockWebsiteChecker, websites)

	want := len(websites)
	got := len(actualResults)
	if want != got {
		t.Fatalf("Wanted %v, got %v", want, got)
	}

	expectedResults := map[string]bool{
		"http://google.com":          true,
		"http://blog.gypsydave5.com": true,
		"waat://furhurterwe.geds":    false,
	}

	if !reflect.DeepEqual(expectedResults, actualResults) {
		t.Fatalf("Wanted %v, got %v", expectedResults, actualResults)
	}
}

// 基准测试使用一百个网址的 slice 对 CheckWebsites 进行测试，并使用 WebsiteChecker 的伪造实现。slowStubWebsiteChecker 故意放慢速度。它使用 time.Sleep 明确等待 20 毫秒，然后返回 true。

func slowStubWebsiteChecker(_ string) bool {
	time.Sleep(20 * time.Millisecond)
	return true
}

func BenchmarkCheckWebsites(b *testing.B) {
	urls := make([]string, 100)
	for i := 0; i < len(urls); i++ {
		urls[i] = "a url"
	}

	for i := 0; i < b.N; i++ {
		CheckWebsites(slowStubWebsiteChecker, urls)
	}
}

/*
通常在 Go 中，当调用函数 doSomething() 时，我们等待它返回（即使它没有值返回，我们仍然等待它完成）。我们说这个操作是 阻塞 的 —— 它让我们等待它完成。Go 中不会阻塞的操作将在称为 goroutine 的单独 进程 中运行。将程序想象成从上到下读 Go 的 代码，当函数被调用执行读取操作时，进入每个函数「内部」。当一个单独的进程开始时，就像开启另一个 reader（阅读程序）在函数内部执行读取操作，原来的 reader 继续向下读取 Go 代码。
要告诉 Go 开始一个新的 goroutine，我们把一个函数调用变成 go 声明，通过把关键字 go 放在它前面：go doSomething()
*/
