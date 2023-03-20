package bootstrap

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xianjianbo/marisa/daemon"
	"github.com/xianjianbo/marisa/library/resource"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init() {
	initMysql()
	initTelegramBot()
	StartBotUpdatesHandler()
}

func initMysql() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:root123456@tcp(127.0.0.1:3306)/marisa?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	resource.MysqlClientGorm = db
	return
}

func initTelegramBot() {
	botClient, err := tgbotapi.NewBotAPI("5967868931:AAGieVOnVt4--6MLqkylT6fzCVROXi7zqYY")
	if err != nil {
		panic(err)
	}
	resource.TelegramBotClient = botClient
	return
}

func StartBotUpdatesHandler() {
	daemon.BotUpdatesHandler()
}
