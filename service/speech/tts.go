package speech

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/xianjianbo/marisa/library/config"
)

type TTSPlayload struct {
	XMLName xml.Name `xml:"speak"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Lang    string   `xml:"lang,attr"`
	Voice   struct {
		Text   string `xml:",chardata"`
		Lang   string `xml:"lang,attr"`
		Gender string `xml:"gender,attr"`
		Name   string `xml:"name,attr"`
	} `xml:"voice"`
}

func TextToSpeech(text string) (audio []byte, err error) {
	url := fmt.Sprintf("https://%s.tts.speech.microsoft.com/cognitiveservices/v1", config.SpeechRegion)
	payload := `<speak version='1.0' xml:lang='en-US'><voice xml:lang='en-US' xml:gender='Female' name='en-CA-ClaraNeural'>` + text + `</voice></speak>`

	clientResp, err := resty.New().R().
		SetHeader("Ocp-Apim-Subscription-Key", config.SpeechKey).
		SetHeader("Content-Type", "application/ssml+xml").
		SetHeader("X-Microsoft-OutputFormat", "ogg-48khz-16bit-mono-opus").
		SetBody(payload).
		Post(url)

	if err != nil {
		err = errors.Wrap(err, "TTS connection")
		return
	}
	if clientResp.StatusCode() != 200 {
		err = errors.New("TTS failed: http code " + strconv.Itoa(clientResp.StatusCode()))
	}

	audio = clientResp.Body()
	return
}
