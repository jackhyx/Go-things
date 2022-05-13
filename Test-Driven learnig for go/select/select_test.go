package _select

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

/*
在标准库中有一个 net/http/httptest 包，它可以让你轻易建立一个 HTTP 模拟服务器（mock HTTP server）。
我们改为使用模拟测试，这样我们就可以控制可靠的服务器来测试了
*/

/*

httptest.NewServer 接受一个我们传入的 匿名函数 http.HandlerFunc。
http.HandlerFunc 是一个看起来类似这样的类型：type HandlerFunc func(ResponseWriter, *Request)。
这些只是说它是一个需要接受一个 ResponseWriter 和 Request 参数的函数，这对于 HTTP 服务器来说并不奇怪。
结果呢，这里并没有什么彩蛋，这也是如何在 Go 语言写一个 真实的 HTTP 服务器的方法。唯一的区别就是我们把它封装成一个易于测试的 httptest.NewServer，它会找一个可监听的端口，然后测试完你就可以关闭它了。
我们让两个服务器中慢的那一个短暂地 time.Sleep 一段时间，当我们请求时让它比另一个慢一些。然后两个服务器都会通过 w.WriteHeader(http.StatusOK) 返回一个 OK 给调用者
func TestRacer(t *testing.T) {

	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(20 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))

	fastServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	slowURL := slowServer.URL
	fastURL := fastServer.URL

	want := fastURL
	got := Racer(slowURL, fastURL)

	if got != want {
		t.Errorf("got '%s', want '%s'", got, want)
	}

	slowServer.Close()
	fastServer.Close()
}
*/

func TestRacer(t *testing.T) {

	t.Run("compares speeds of servers, returning the url of the fastest one", func(t *testing.T) {
		slowServer := makeDelayedServer(20 * time.Millisecond)
		fastServer := makeDelayedServer(0 * time.Millisecond)

		defer slowServer.Close()
		defer fastServer.Close()

		slowURL := slowServer.URL
		fastURL := fastServer.URL

		want := fastURL
		got, err := Racer(slowURL, fastURL)

		if err != nil {
			t.Fatalf("did not expect an error but got one %v", err)
		}

		if got != want {
			t.Errorf("got '%s', want '%s'", got, want)
		}
	})

	t.Run("returns an error if a server doesn't respond within 10s", func(t *testing.T) {
		server := makeDelayedServer(25 * time.Millisecond)

		defer server.Close()

		_, err := ConfigurableRacer(server.URL, server.URL, 20*time.Millisecond)

		if err == nil {
			t.Error("expected an error but didn't get one")
		}
	})
}
func makeDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
}

/*

对每个 URL：
我们用 time.Now() 来记录请求 URL 前的时间。
然后用 http.Get 来请求 URL 的内容。这个函数返回一个 http.Response 和一个 error，但目前我们不关心它们的值。
time.Since 获取开始时间并返回一个 time.Duration 时间差。

*/
/*
func Racer(a, b string) (winner string) {
	aDuration := measureResponseTime(a)
	bDuration := measureResponseTime(b)

	if aDuration < bDuration {
		return a
	}

	return b
}

func measureResponseTime(url string) time.Duration {
	start := time.Now()
	http.Get(url)
	return time.Since(start)
}
*/

// 进程同步
// Go 在并发方面很在行，为什么我们要一个接一个地测试哪个网站更快呢？我们应该能够同时测试两个。
// 我们并不关心请求的 准确响应时间，我们只是需要知道哪个更快返回而已
var tenSecondTimeout = 10 * time.Second

func Racer(a, b string) (winner string, error error) {
	return ConfigurableRacer(a, b, tenSecondTimeout)
}

func ConfigurableRacer(a, b string, timeout time.Duration) (winner string, error error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
	}
}

//使用 select 时，time.After 是一个很好用的函数。当你监听的 channel 永远不会返回一个值时你可以潜在地编写永远阻塞的代码，尽管在我们的案例中它没有发生。time.After 会在你定义的时间过后发送一个信号给 channel 并返回一个 chan 类型（就像 ping 那样）。
//对我们来说这完美了；如果 a 或 b 谁胜出就返回谁，但如果测试达到 10 秒，那么 time.After 会发送一个信号并返回一个 error。

func ping(url string) chan bool {
	ch := make(chan bool)
	go func() {
		http.Get(url)
		ch <- true
	}()
	return ch
}
