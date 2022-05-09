package main

import (
	"strings"
	"testing"
)

func TestFileSystemStore(t *testing.T) {
	// 我们使用 strings.NewReader 会返回一个 Reader，这是我们的 FileSystemStore 函数中用来读取数据的。
	// 在 main 中我们将打开一个文件，它也是一个 Reader。
	t.Run("league from a reader", func(t *testing.T) {
		database := strings.NewReader(`[
			{"Name": "Cleo", "Wins": 10},
			{"Name": "Chris", "Wins": 33}]`)

		store := FileSystemPlayerStore{database}

		got := store.GetLeague()

		want := []Player{
			{"Cleo", 10},
			{"Chris", 33},
		}

		assertLeague(t, got, want)

		// read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})
}
