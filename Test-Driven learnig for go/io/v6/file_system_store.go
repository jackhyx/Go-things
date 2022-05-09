package main

import (
	"encoding/json"
	"io"
)

// FileSystemPlayerStore stores players in the filesystem.
// 我们可以创建一个构造函数，该构造函数可以为我们执行一些初始化操作，并将 league 作为值存储在我们的 FileSystemStore 中，以便在读取中使用。
type FileSystemPlayerStore struct {
	database io.ReadWriteSeeker
	league   League
}

// NewFileSystemPlayerStore creates a FileSystemPlayerStore.
func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore {
	database.Seek(0, 0)
	league, _ := NewLeague(database)

	return &FileSystemPlayerStore{
		database: database,
		league:   league,
	}
}

// 这样，我们只需从磁盘读取一次。我们现在可以替换以前的所有从磁盘上获得 league 的调用，并且只使用 f.league。
// GetLeague returns the scores of all the players.

func (f *FileSystemPlayerStore) GetLeague() League {
	return f.league
}

// GetPlayerScore retrieves a player's score.
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {

	player := f.league.Find(name)

	if player != nil {
		return player.Wins
	}

	return 0
}

// RecordWin will store a win for a player, incrementing wins if already known.
func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}

	f.database.Seek(0, 0)
	json.NewEncoder(f.database).Encode(f.league)
}
