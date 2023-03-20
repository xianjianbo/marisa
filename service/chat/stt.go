package chat

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func Ogg2Wav(oggFilePath string, wavFilePath string) (err error) {
	ffmpegBin := filepath.Join(GetProjectRootDir(), "bin", "ffmpeg")
	ffmpegCmd := exec.Command(ffmpegBin, "-i", oggFilePath, wavFilePath)
	if err = ffmpegCmd.Run(); err != nil {
		err = errors.Wrap(err, "ffmpegCmd")
		return
	}

	return
}

func DownloadVoiceFile(voiceURL string, outputFilePath string) (err error) {
	resp, err := http.Get(voiceURL)
	if err != nil {
		err = errors.Wrap(err, "get file url")
	}
	defer resp.Body.Close()

	out, err := os.Create(outputFilePath)
	if err != nil {
		err = errors.Wrap(err, "create file")
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		err = errors.Wrap(err, "copy file")
	}

	return
}

func (c *ChatService) RecognizeVoice(voiceURL string) (text string, err error) {
	oggFlePath := filepath.Join(GetVoiceDataDir(), fmt.Sprintf("/user/%s.ogg", uuid.New().String()))
	wavFlePath := filepath.Join(GetVoiceDataDir(), fmt.Sprintf("/user/%s.wav", uuid.New().String()))

	err = DownloadVoiceFile(voiceURL, oggFlePath)
	if err != nil {
		err = errors.Wrap(err, "DownloadVoiceFile")
		return
	}

	err = Ogg2Wav(oggFlePath, wavFlePath)
	if err != nil {
		err = errors.Wrap(err, "Ogg2Wav")
		return
	}

	text, err = STT(wavFlePath)
	if err != nil {
		err = errors.Wrap(err, "STT")
		return
	}

	return
}

func STT(wavFilePath string) (text string, err error) {
	speechKey := os.Getenv("SPEECH_KEY")
	speechRegion := os.Getenv("SPEECH_REGION")

	audioConfig, err := audio.NewAudioConfigFromWavFileInput(wavFilePath)
	if err != nil {
		err = errors.Wrap(err, "NewAudioConfigFromWavFileInput")
		return
	}
	defer audioConfig.Close()

	speechConfig, err := speech.NewSpeechConfigFromSubscription(speechKey, speechRegion)
	if err != nil {
		err = errors.Wrap(err, "NewSpeechConfigFromSubscription")
		return
	}
	defer speechConfig.Close()

	speechRecognizer, err := speech.NewSpeechRecognizerFromConfig(speechConfig, audioConfig)
	if err != nil {
		err = errors.Wrap(err, "NewSpeechRecognizerFromConfig")
		return
	}
	defer speechRecognizer.Close()

	finished := make(chan struct{}, 1)
	speechRecognizer.Recognized(func(event speech.SpeechRecognitionEventArgs) {
		defer event.Close()
		text += event.Result.Text
	})
	speechRecognizer.Canceled(func(event speech.SpeechRecognitionCanceledEventArgs) {
		defer event.Close()
		finished <- struct{}{}
	})

	speechRecognizer.StartContinuousRecognitionAsync()
	defer speechRecognizer.StopContinuousRecognitionAsync()

	select {
	case <-finished:
	case <-time.After(20 * time.Second):
		err = errors.New("ContinuousRecognitionAsync time out")
		return
	}

	return
}
