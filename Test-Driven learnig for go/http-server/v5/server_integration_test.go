package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
集成测试对较大型的测试很有用，但你必须牢记：
集成测试更难编写
测试失败时，可能很难知道原因（通常它是集成测试组件中的错误），因此可能更难修复
有时运行较慢（因为它们通常与“真实”组件一起使用，比如数据库）
因此，建议你研究一下 金字塔测试。

我们正在尝试集成两个组件：InMemoryPlayerStore 和 PlayerServer。
然后我们发起 3 个请求，为玩家记录 3 次获胜。我们并不太关心测试中的返回状态码，因为和集成得好不好无关。
我们真正关心的是下一个响应（所以我们用变量存储 response），因为我们要尝试并获得 player 的得分
*/
func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := NewInMemoryPlayerStore()
	server := PlayerServer{store}
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))
	assertStatus(t, response.Code, http.StatusOK)

	assertResponseBody(t, response.Body.String(), "3")
}
