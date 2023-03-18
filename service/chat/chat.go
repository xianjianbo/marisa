package chat

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xianjianbo/marisa/model/record"
)

type ChatInput struct {
	Ask       string `json:"ask"`
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
}

type ChatOutput struct {
	Reply     string `json:"reply"`
	SessionID string `json:"session_id"`
}

func (c *ChatService) Chat(ctx *gin.Context, input ChatInput) (output ChatOutput, err error) {
	messages := []Message{
		{
			Role:    RoleSystem,
			Content: newPromptSpokenEnglishTeacherAndImprover(),
		},
	}

	if input.SessionID == "" {
		input.SessionID = uuid.NewString()
	} else {
		historys := make([]*record.Record, 0)
		historys, err = c.RecordModel.GetRecordsBySessionID(ctx, input.SessionID, 4)
		if err != nil {
			return
		}

		for _, history := range historys {
			messages = append(messages, Message{
				Role:    history.Role,
				Content: history.Content,
			})
		}
	}

	messages = append(messages, Message{
		Role:    RoleUser,
		Content: input.Ask,
	})

	if output.Reply, err = GPT(ctx, messages); err != nil {
		return
	}

	if err = c.RecordModel.InsertRecords(ctx, nil, []*record.Record{
		{
			SessionID: input.SessionID,
			UserID:    input.UserID,
			Role:      RoleUser,
			Content:   input.Ask,
		},
		{
			SessionID: input.SessionID,
			UserID:    input.UserID,
			Role:      RoleAssistant,
			Content:   output.Reply,
		},
	}); err != nil {
		return
	}

	output.SessionID = input.SessionID

	return
}
