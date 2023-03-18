package record

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 聊天记录表
type Record struct {
	ID         int64     `gorm:"column:id" db:"id" json:"id" form:"id"`                                                  //主键
	CreateTime time.Time `gorm:"column:create_time;default:null" db:"create_time" json:"create_time" form:"create_time"` //创建记录时间
	UpdateTime time.Time `gorm:"column:update_time;default:null" db:"update_time" json:"update_time" form:"update_time"` //更新记录时间
	SessionID  string    `gorm:"column:session_id" db:"session_id" json:"session_id" form:"session_id"`                  //会话id
	UserID     string    `gorm:"column:user_id" db:"user_id" json:"user_id" form:"user_id"`                              //用户id
	Role       string    `gorm:"column:role" db:"role" json:"role" form:"role"`                                          //角色: system, user, assistant
	Content    string    `gorm:"column:content" db:"content" json:"content" form:"content"`                              //消息
}

type RecordModel interface {
	InsertRecords(ctx *gin.Context, tx *gorm.DB, records []*Record) (err error)
	GetRecordsBySessionID(ctx *gin.Context, sessionID string, size int) (records []*Record, err error)
}
