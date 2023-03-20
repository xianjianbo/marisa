package chat

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
	"github.com/google/uuid"
)

func (c *ChatService) TTS(text string) (int64, []byte, error) {

	fileName := uuid.NewString()
	wavFilePath := filepath.Join(GetVoiceDataDir(), fmt.Sprintf("%s.wav", fileName))
	oggFilePath := filepath.Join(GetVoiceDataDir(), fmt.Sprintf("%s.ogg", fileName))

	defer func() {
		// TODO
		// if errRemove := os.Remove(wavFilePath); errRemove != nil {
		// 	fmt.Println("remove fail:" + wavFilePath)
		// }
		// if errRemove := os.Remove(oggFilePath); errRemove != nil {
		// 	fmt.Println("remove fail:" + wavFilePath)
		// }
	}()

	if err := Text2WavAudio(text, wavFilePath); err != nil {
		err = errors.Wrap(err, "Text2WavAudio")
		return 0, nil, err
	}

	d, err := Wav2Ogg(wavFilePath, oggFilePath)
	if err != nil {
		err = errors.Wrap(err, "Wav2Ogg")
		return 0, nil, err
	}

	oggData, err := ioutil.ReadFile(oggFilePath)
	if err != nil {
		return 0, nil, err
	}

	return d, oggData, nil
}

// TODO
func GetProjectRootDir() string {
	projectRootDir, _ := os.Getwd()
	return projectRootDir
}

func GetVoiceDataDir() string {
	voiceDataDir := filepath.Join(GetProjectRootDir(), "data", "voice")

	// Create if not exsit
	_, err := os.Stat(voiceDataDir)
	if os.IsNotExist(err) {
		_ = os.MkdirAll(voiceDataDir, os.ModePerm)
	}
	return voiceDataDir
}

func Wav2Ogg(wavFilePath string, oggFilePath string) (duration int64, err error) {
	ffmpegBin := filepath.Join(GetProjectRootDir(), "bin", "ffmpeg")
	ffmpegCmd := exec.Command(ffmpegBin, "-i", wavFilePath, "-acodec", "libopus", oggFilePath)
	if err = ffmpegCmd.Run(); err != nil {
		err = errors.Wrap(err, "ffmpegCmd")
		return
	}

	durationCmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("%s -i %s  2>&1 | grep 'Duration' | cut -d ' ' -f 4 | sed s/,//", ffmpegBin, wavFilePath))
	outputDuration, err := durationCmd.CombinedOutput()
	if err != nil {
		err = errors.Wrap(err, "durationCmd")
		return
	}

	if duration, err = DurationFormatHour2Second(string(outputDuration)); err != nil {
		return
	}

	return
}

// DurationHour2Second 将ffmpeg返回的小时格式duration转换为妙
// 如 "00:01:09.00" 转换为 69 秒
func DurationFormatHour2Second(hour string) (second int64, err error) {
	hour = strings.TrimSpace(hour)
	if !strings.Contains(hour, ":") {
		err = errors.New("ErrInvalidDurationString")
		return
	}
	timeArr := strings.Split(hour, ":")
	if len(timeArr) != 3 {
		err = errors.New("ErrInvalidDurationString")
		return
	}
	h, err := strconv.ParseFloat(timeArr[0], 64)
	if err != nil {
		return
	}
	m, err := strconv.ParseFloat(timeArr[1], 64)
	if err != nil {
		return
	}
	s, err := strconv.ParseFloat(timeArr[2], 64)
	if err != nil {
		return
	}

	return int64(h*3600 + m*60 + s), nil
}

func Text2WavAudio(text string, filePath string) error {
	// This example requires environment variables named "SPEECH_KEY" and "SPEECH_REGION"
	speechKey := os.Getenv("SPEECH_KEY")
	speechRegion := os.Getenv("SPEECH_REGION")

	audioConfig, err := audio.NewAudioConfigFromWavFileOutput(filePath)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return err
	}
	defer audioConfig.Close()
	speechConfig, err := speech.NewSpeechConfigFromSubscription(speechKey, speechRegion)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return err
	}
	defer speechConfig.Close()

	speechConfig.SetSpeechSynthesisVoiceName("en-CA-ClaraNeural")

	speechSynthesizer, err := speech.NewSpeechSynthesizerFromConfig(speechConfig, audioConfig)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return err
	}
	defer speechSynthesizer.Close()

	task := speechSynthesizer.SpeakTextAsync(text)
	var outcome speech.SpeechSynthesisOutcome
	select {
	case outcome = <-task:
	case <-time.After(60 * time.Second):
		fmt.Println("Timed out")
		return err
	}
	defer outcome.Close()
	if outcome.Error != nil {
		fmt.Println("Got an error: ", outcome.Error)
		return err
	}

	if outcome.Result.Reason != common.SynthesizingAudioCompleted {
		cancellation, _ := speech.NewCancellationDetailsFromSpeechSynthesisResult(outcome.Result)
		fmt.Printf("CANCELED: Reason=%d.\n", cancellation.Reason)

		if cancellation.Reason == common.Error {
			fmt.Printf("CANCELED: ErrorCode=%d\nCANCELED: ErrorDetails=[%s]\nCANCELED: Did you set the speech resource key and region values?\n",
				cancellation.ErrorCode,
				cancellation.ErrorDetails)
		}
	}

	return nil

}
