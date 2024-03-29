package chat

import (
	"context"

	"github.com/google/uuid"
	"github.com/xianjianbo/marisa/model/record"
)

type ChatInput struct {
	Ask       string `json:"ask"`
	UserID    string `json:"user_id"`
	UserName  string `json:"user_name"`
	SessionID string `json:"session_id"`
}

type ChatOutput struct {
	Reply     string `json:"reply"`
	SessionID string `json:"session_id"`
}

func (c *ChatService) Chat(ctx context.Context, input ChatInput) (output ChatOutput, err error) {
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
		historys, err = c.RecordModel.GetRecordsBySessionID(ctx, input.SessionID, 6)
		if err != nil {
			return
		}

		for i := len(historys) - 1; i >= 0; i-- {
			history := historys[i]
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
			UserName:  input.UserName,
			Role:      RoleUser,
			Content:   input.Ask,
		},
		{
			SessionID: input.SessionID,
			UserID:    input.UserID,
			UserName:  input.UserName,
			Role:      RoleAssistant,
			Content:   output.Reply,
		},
	}); err != nil {
		return
	}

	output.SessionID = input.SessionID

	return
}
