package speech

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/xianjianbo/marisa/library/config"
)

type STTResponse struct {
	RecognitionStatus string `json:"RecognitionStatus"`
	Offset            int    `json:"Offset"`
	Duration          int    `json:"Duration"`
	DisplayText       string `json:"DisplayText"`
}

func SpeechToText(audio []byte) (text string, err error) {
	url := fmt.Sprintf("https://%s.stt.speech.microsoft.com/speech/recognition/conversation/cognitiveservices/v1?language=en-US", config.SpeechRegion)
	clientResp, err := resty.New().R().
		SetHeader("Content-Type", "audio/ogg; codecs=opus").
		SetHeader("Ocp-Apim-Subscription-Key", config.SpeechKey).
		SetBody(audio).
		Post(url)
	if err != nil {
		err = errors.Wrap(err, "STT connection")
		return
	}

	sttResponse := STTResponse{}
	if err = json.Unmarshal(clientResp.Body(), &sttResponse); err != nil {
		err = errors.Wrap(err, "STTResponse error, http code "+strconv.Itoa(clientResp.StatusCode()))
		return
	}
	if sttResponse.RecognitionStatus != "Success" {
		err = errors.New("STT failed")
		return
	}

	text = sttResponse.DisplayText
	return
}
