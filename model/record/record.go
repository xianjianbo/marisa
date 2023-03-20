package record

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// 聊天记录表
type Record struct {
	ID         int64     `gorm:"column:id" db:"id" json:"id" form:"id"`
	CreateTime time.Time `gorm:"column:create_time;default:null" db:"create_time" json:"create_time" form:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;default:null" db:"update_time" json:"update_time" form:"update_time"`
	SessionID  string    `gorm:"column:session_id" db:"session_id" json:"session_id" form:"session_id"`
	UserID     string    `gorm:"column:user_id" db:"user_id" json:"user_id" form:"user_id"`
	UserName   string    `gorm:"column:user_name" db:"user_name" json:"user_name" form:"user_name"`
	Role       string    `gorm:"column:role" db:"role" json:"role" form:"role"`
	Content    string    `gorm:"column:content" db:"content" json:"content" form:"content"`
}

type RecordModel interface {
	InsertRecords(ctx context.Context, tx *gorm.DB, records []*Record) (err error)
	GetRecordsBySessionID(ctx context.Context, sessionID string, size int) (records []*Record, err error)
}
