package daemon

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/forPelevin/gomoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xianjianbo/marisa/library/resource"
	chatservice "github.com/xianjianbo/marisa/service/chat"
	speechservice "github.com/xianjianbo/marisa/service/speech"
)

func BotUpdatesHandler() {
	uconfig := tgbotapi.NewUpdate(0)
	uconfig.Timeout = 60
	updates := resource.TelegramBotClient.GetUpdatesChan(uconfig)

	for update := range updates {
		if update.Message != nil {
			reply, voice := HandleMessage(*update.Message)

			textMsg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			resource.TelegramBotClient.Send(textMsg)

			if len(voice) == 0 {
				continue
			}

			voiceMsg := tgbotapi.NewVoice(update.Message.Chat.ID, tgbotapi.FileReader{Reader: bytes.NewReader(voice)})
			resource.TelegramBotClient.Send(voiceMsg)
		}
	}
}

func HandleMessage(message tgbotapi.Message) (reply string, voice []byte) {
	var (
		ask string
		err error
	)
	chatService := chatservice.NewChatService()

	switch {
	case message.Text != "":
		ask = message.Text
	case message.Voice != nil:
		var voiceURL string
		if voiceURL, err = resource.TelegramBotClient.GetFileDirectURL(message.Voice.FileID); err != nil {
			log.Println("GetFileDirectURL error: ", err)
			reply = MarisaReply2Text("Soory...I can't recognise this message currently. Please try again later.")
			return
		}

		resp, err := http.Get(voiceURL)
		if err != nil {
			log.Println("http.Get error: ", err)
			reply = MarisaReply2Text("Soory...I can't recognise this message currently. Please try again later.")
			return
		}
		defer resp.Body.Close()

		voiceBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("io.ReadAllL error: ", err)
			reply = MarisaReply2Text("Soory...I can't recognise this message currently. Please try again later.")
			return
		}

		if ask, err = speechservice.SpeechToText(voiceBytes); err != nil {
			log.Println("RecognizeVoice error: ", err)
			reply = MarisaReply2Text("Soory...I can't recognise this message currently. Please try again later.")
			return
		}
	default:
		// unsupported message type
		reply = MarisaReply2Text("Soory...I can't recognise this message currently. You can try text or voice message.")
		return
	}

	log.Println("### UserID: ", strconv.Itoa(int(message.From.ID)))
	log.Println("### "+GetTGUserName(*message.From)+": ", ask)

	var chatOutput chatservice.ChatOutput
	if chatOutput, err = chatService.Chat(context.Background(), chatservice.ChatInput{
		SessionID: strconv.FormatInt(message.Chat.ID, 10),
		UserID:    strconv.FormatInt(message.From.ID, 10),
		UserName:  GetTGUserName(*message.From),
		Ask:       ask,
	}); err != nil {
		log.Println("got a chat error: ", err)
		reply = MarisaReply2Text("Oops, something went wrong.")
		return
	}

	reply = MarisaReply2Text(chatOutput.Reply)
	if message.Voice != nil {
		reply = MarisaReply2Voice(ask, chatOutput.Reply)
	}

	log.Println("### Marisa: ", chatOutput.Reply)

	if voice, err = speechservice.TextToSpeech(gomoji.RemoveEmojis(chatOutput.Reply)); err != nil {
		log.Println("TextToSpeech error: ", err)
		return
	}

	return
}

func GetTGUserName(user tgbotapi.User) string {
	if user.UserName != "" {
		return user.UserName
	}

	names := make([]string, 0)
	if user.FirstName != "" {
		names = append(names, user.FirstName)
	}
	if user.LastName != "" {
		names = append(names, user.LastName)
	}
	return strings.Join(names, " ")

}

func MarisaReply2Text(replyText string) string {
	return fmt.Sprintf("Marisa: %s", replyText)
}

func MarisaReply2Voice(askText, replyText string) string {
	return fmt.Sprintf("You: %s\n\nMarisa: %s", askText, replyText)
}
