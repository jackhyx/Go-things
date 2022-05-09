package main

import (
	"fmt"
	"net/http"
	"strings"
)

// 我们把得分计算从 handler 移到函数 GetPlayerScore 中，这就是使用接口重构的正确方法

// PlayerStore stores score information about players.
type PlayerStore interface {
	GetPlayerScore(name string) int
}

// PlayerServer is a HTTP interface for player information.
type PlayerServer struct {
	store PlayerStore
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	// 调用 store.GetPlayerStore 来获得得分
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}
