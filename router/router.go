package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xianjianbo/marisa/controller/chat"
)

func InitRouter() {
	r := gin.Default()
	groupV1 := r.Group("v1")
	groupV1.POST("/chat", chat.Chat)
	r.Run() // listen and serve on 0.0.0.0:8080
}
