package daemon

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xianjianbo/marisa/library/resource"
	chatservice "github.com/xianjianbo/marisa/service/chat"
)

func BotUpdatesHandler() {
	chatService := chatservice.NewChatService()

	uconfig := tgbotapi.NewUpdate(0)
	uconfig.Timeout = 60
	updates := resource.TelegramBotClient.GetUpdatesChan(uconfig)

	for update := range updates {
		if update.Message != nil {

			var ask, reply string
			var err error
			var chatOutput chatservice.ChatOutput
			switch {
			case update.Message.Text != "":
				ask = update.Message.Text
			case update.Message.Voice != nil:
				voiceURL, err := resource.TelegramBotClient.GetFileDirectURL(update.Message.Voice.FileID)
				if err != nil {
					reply = MarisaReply2Text("Soory...I can't recognise this message currently. Please try again later.")
					goto ToReply
				}

				ask, err = chatService.RecognizeVoice(voiceURL)
				if err != nil {
					reply = MarisaReply2Text("Soory...I can't recognise this message currently. Please try again later.")
					goto ToReply
				}

			default:
				// unsupported message type
				reply = MarisaReply2Text("Soory...I can't recognise this message currently. You can try text or voice message.")
				goto ToReply
			}

			log.Printf("%s: %s", GetUserName(*update.Message.From), ask)
			chatOutput, err = chatService.Chat(context.Background(), chatservice.ChatInput{
				SessionID: strconv.FormatInt(update.Message.Chat.ID, 10),
				UserID:    strconv.FormatInt(update.Message.From.ID, 10),
				UserName:  GetUserName(*update.Message.From),
				Ask:       ask,
			})
			if err != nil {
				fmt.Println("Got a chat error: ", err)
				reply = MarisaReply2Text("Oops, something went wrong.")
			} else {
				reply = MarisaReply2Text(chatOutput.Reply)
				if update.Message.Voice != nil {
					reply = MarisaReply2Voice(ask, chatOutput.Reply)
				}
			}
			log.Printf("%s: %s", "Marisa", chatOutput.Reply)

		ToReply:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			resource.TelegramBotClient.Send(msg)

			voice := tgbotapi.NewVoice(update.Message.Chat.ID, tgbotapi.FileReader{
				Reader: bytes.NewReader(chatOutput.OggVoice),
			})
			resource.TelegramBotClient.Send(voice)

		}
	}
}

func GetUserName(user tgbotapi.User) string {
	if user.UserName != "" {
		return user.UserName
	} else {
		names := make([]string, 0)
		if user.FirstName != "" {
			names = append(names, user.FirstName)
		}
		if user.LastName != "" {
			names = append(names, user.LastName)
		}
		return strings.Join(names, " ")
	}
}

func MarisaReply2Text(replyText string) string {
	return `Marisa: ` + replyText
}

func MarisaReply2Voice(askText, replyText string) string {
	return fmt.Sprintf("You: %s\n\nMarisa: %s", askText, replyText)
}
