package resource

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

var MysqlClientGorm *gorm.DB

var TelegramBotClient *tgbotapi.BotAPI
