package main

import (
	"fmt"
	"net/http"
	"strings"
)

// PlayerStore stores score information about players.
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
}

/*
我们更改了 PlayerServer 的第二个属性，删除了命名属性 router http.ServeMux，并用 http.Handler 替换了它；这被称为 嵌入。
Go 没有提供典型的，类型驱动的子类化概念，但它具有通过在结构或接口中嵌入类型来“借用”一部分实现的能力。
高效 Go - 嵌入
这意味着我们的 PlayerServer 现在已经有了 http.Handler 所有的方法，也就是 ServeHTTP。
为了“填充” http.Handler，我们将它分配给我们在 NewPlayerServer 中创建的 router。我们可以这样做是因为 http.ServeMux 具有 ServeHTTP 方法。
这允许我们删除我们的 ServeHTTP 方法，因为我们已经通过嵌入类型公开了它。
嵌入是一个非常有意思的语法特性。你可以用它将接口组成新的接口。

你必须小心使用嵌入类型，因为你将公开所有嵌入类型的公共方法和字段。在我们的例子中它是可以的，因为我们只是嵌入了 http.Handler 这个 接口。
如果我们懒一点，嵌入了 http.ServeMux（混合类型），它仍然可以工作 但 PlayerServer 的用户就可以给我们的服务器添加新路由了，因为 Handle(path, handler) 会公开。
嵌入类型时，真正要考虑的是对你公开的 API 有什么影响。
滥用嵌入最终会污染你的 API，并暴露你的类型的内部信息，这是个常见的错误。
现在我们重新构建了我们的应用，我们可以轻易地添加新的路由，并让 /league 端点有了一个新的开始。现我我们需要让它返回一些有用的信息。

*/

// PlayerServer is a HTTP interface for player information.
type PlayerServer struct {
	store        PlayerStore
	http.Handler // router *http.ServeMux
}

// 把一个路由作为一个请求来处理并调用它挺奇怪的（并且效率低下）。我们想要的理想情况是有一个 NewPlayerServer 这样的函数，它可以取得依赖并进行一次创建路由的设置。每个请求都可以使用该路由的一个实例。
// NewPlayerServer creates a PlayerServer with routing configured.

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)

	p.store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	p.Handler = router

	return p
}

/*  func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}
PlayerServer 现在需要储存一个路由。
我们已经把创建 ServeHTTP 路由的动作移到 NewPlayerServer，这样只需要完成一次，而不是每次请求都要做。
在所有测试和程序代码中，用到 PlayerServer{&store} 的地方你需要更新为 NewPlayerServer(&store)。
*/

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
