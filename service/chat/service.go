package chat

import (
	"github.com/xianjianbo/marisa/library/resource"
	"github.com/xianjianbo/marisa/model/record"
	recordmodel "github.com/xianjianbo/marisa/model/record/gorm"
)

type ChatService struct {
	RecordModel record.RecordModel
}

func NewChatService() *ChatService {
	return &ChatService{
		RecordModel: recordmodel.NewRecordModel(resource.MysqlClientGorm),
	}
}
