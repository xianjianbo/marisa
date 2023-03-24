package bootstrap

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xianjianbo/marisa/daemon"
	"github.com/xianjianbo/marisa/library/config"
	"github.com/xianjianbo/marisa/library/resource"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init() {
	initConfig()
	initMysql()
	initTelegramBot()
	StartBotUpdatesHandler()
}

func initConfig() {
	config.SpeechRegion = os.Getenv("MS_SPEECH_REGION")
	if config.SpeechRegion == "" {
		panic("env MS_SPEECH_REGION is null")
	}
	config.SpeechKey = os.Getenv("MS_SPEECH_KEY")
	if config.SpeechKey == "" {
		panic("env MS_SPEECH_KEY is null")
	}
	config.OpenAIApiKey = os.Getenv("OPENAI_API_KEY")
	if config.OpenAIApiKey == "" {
		panic("env OPENAI_API_KEY is null")
	}
}

func initMysql() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	mysqlUserName := os.Getenv("MYSQL_USERNAME")
	mysqPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlDBNameMarisa := os.Getenv("MYSQL_DB_NAME_MARISA")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlUserName, mysqPassword, mysqlHost, mysqlDBNameMarisa)
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
	go daemon.BotUpdatesHandler()
}
