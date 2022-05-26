```
package main
// 获取gin
import "github.com/gin-gonic/gin"

// 主函数
func main() {
	// 取r是router的缩写
    r := gin.Default()
	// 这里非常简单，很像deno、node的路由吧
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    // 监听端口8080
    r.Run(":8080")
}

```

#### 搭建脚手架 server.go
```go
package server

import (
	"github.com/gin-gonic/gin"
)
// 这里是定义一个接口，解决上述弊端的规范性
type IController interface {
	// 这个传参就是脚手架主程
	Router(server *Server)
}

// 定义一个脚手架
type Server struct {
	*gin.Engine
    // 路由分组一会儿会用到
	g *gin.RouterGroup
}

// 初始化函数
func Init() *Server {
	// 作为Server的构造器
	s := &Server{Engine: gin.New()}
    // 返回作为链式调用
	return s
}

// 监听函数，更好的做法是这里的端口应该放到配置文件
func (this *Server) Listen() {
	this.Run(":8080")
}

// 这里是路由的关键代码，这里会挂载路由
func (this *Server) Route(controllers ...IController) *Server {
	// 遍历所有的控制层，这里使用接口，就是为了将Router实例化
	for _, c := range controllers {
		c.Router(this)
	}
	return this
}

func (this *Server) GroupRouter(group string, controllers ...IController) *Server {
	this.g = this.Group(group)
	for _, c := range controllers {
		c.Router(this)
	}
	return this
}

```
#### 控制层 controller.go
```go
package controller
import (
	"github.com/gin-gonic/gin"
	"feihu/server"
)

// 这里的gin引擎直接移到脚手架server里
type UserController struct {
}

// 这里是构造函数
func NewUserController() *UserController {
	return &UserController{}
}

// 这里是业务方法
func (this *UserController) GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"data": "hello world",
		})
	}
}

// 这里依然是处理路由的地儿，而由于我们定义了接口规范，就必须实现Router方法
func (this *UserController) Router (server *server.Server) {
	server.Handle("GET", "/", this.GetUser())
}

```

主函数
```go
package main

import (
	. "feihu/controller"
	"feihu/server"
)

func main () {
	server.
		Init().
		Route(
			NewUserController(),
		).
        // 这里就是路由分组啦
		GroupRouter("v1",
			NewOrderController(),
		).
		Listen()
}

```