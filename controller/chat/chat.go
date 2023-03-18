package chat

import (
	"net/http"

	"github.com/gin-gonic/gin"
	chatservice "github.com/xianjianbo/marisa/service/chat"
)

func Chat(ctx *gin.Context) {
	var input chatservice.ChatInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err_msg": err.Error(),
		})
		return
	}

	resp, err := chatservice.NewChatService().Chat(ctx, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err_msg": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
	return
}
