package main

import (
	"log"
	"net/http"
	"os"
)

const dbFileName = "game.db.json"

// 我们创建了一个文件作为数据库。
// 第 2 个参数 os.OpenFile 允许你定义打开文件的权限，在我们的例子中，O_RDWR 意味着我们想要读写权限，os.O_CREATE 是指如果文件不存在，则创建该文件。
// 第 3 个参数表示设置文件的权限，在我们的示例中，所有用户都可以读写文件
func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store := &FileSystemPlayerStore{db}
	server := NewPlayerServer(store)

	log.Fatal(http.ListenAndServe(":5000", server))
}
