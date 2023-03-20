package gorm

import (
	"context"

	"github.com/xianjianbo/marisa/model/record"
	"gorm.io/gorm"
)

var _ record.RecordModel = &RecordModel{}

type RecordModel struct {
	DB *gorm.DB
}

func NewRecordModel(db *gorm.DB) record.RecordModel {
	return &RecordModel{
		DB: db,
	}
}

const (
	RecordTableName string = "marisa_chat_record"
)

func (d *RecordModel) InsertRecords(ctx context.Context, tx *gorm.DB, records []*record.Record) (err error) {
	if tx == nil {
		tx = d.DB
	}
	err = tx.WithContext(ctx).Table(RecordTableName).Create(records).Error
	return
}

func (d *RecordModel) GetRecordsBySessionID(ctx context.Context, sessionID string, size int) (records []*record.Record, err error) {
	err = d.DB.WithContext(ctx).Table(RecordTableName).Where(map[string]interface{}{
		"session_id": sessionID,
	}).Order("id").Limit(size).Scan(&records).Error
	return
}
