package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
为了测试服务器，我们需要通过 Request 来发送请求，并期望监听到 handler 向 ResponseWriter 写入了什么。

我们用 http.NewRequest 来创建一个请求。第一个参数是请求方法，第二个是请求路径。
nil 是请求实体，不过在这个场景中不用发送请求实体。

net/http/httptest 自带一个名为 ResponseRecorder 的监听器，所以我们可以用这个。
它有很多有用的方法可以检查应答被写入了什么;
*/
func TestGETPlayers(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	PlayerServer(response, request)

	t.Run("returns Pepper's score", func(t *testing.T) {
		got := response.Body.String()
		want := "20"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

}
