package tgfileutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mymmrac/telego"
	"io"
	"net/http"
	"os"
	"strings"
)

type MsgFile struct {
	Type      string
	MessageID int
	Date      int
	From      int64
	Data      string
}

func (m *MsgFile) GetData(b *telego.Bot) (string, error) {
	if m.Type != "text" {
		var getFileParams telego.GetFileParams
		getFileParams.FileID = m.Data
		file, err := b.GetFile(&getFileParams)
		if err == nil {
			return file.FilePath, nil
		} else {
			return "", fmt.Errorf("Error getting FileID from Telegram API")
		}
	} else {
		return m.Data, nil
	}
}

func GetMsgFile(message *telego.Message) (MsgFile, error) {
	var output MsgFile

	output.Date = int(message.Date)
	output.MessageID = message.MessageID
	output.From = message.From.ID

	if message.Audio != nil {
		output.Data = message.Audio.FileID
		output.Type = "audio"
		return output, nil
	} else if message.Document != nil {
		output.Data = message.Document.FileID
		output.Type = "doc"
		return output, nil
	} else if message.Video != nil {
		output.Type = "video"
		output.Data = message.Video.FileID
		return output, nil
	} else if message.Voice != nil {
		output.Type = "voice"
		output.Data = message.Voice.FileID
		return output, nil
	} else if len(message.Photo) != 0 {
		output.Type = "photo"
		output.Data = message.Photo[0].FileID
		return output, nil
	} else if len(message.Text) != 0 {
		output.Type = "text"
		output.Data = message.Text
		return output, nil
	} else {
		return output, fmt.Errorf("Unkown message type")
	}
}

// TODO: Make the downloaded image as Full size
// TODO: Make test with all file types
func DownloadFile(msgFile *MsgFile, token string, botBaseUrl string, filePath string) (int64, error) {
	if msgFile.Type != "text" {
		segments := strings.Split(filePath, "/")
		fileName := segments[len(segments)-1]
		fullURLFile := fmt.Sprintf("%+v/file/bot%+v/%+v", botBaseUrl, token, filePath)

		file, err := os.Create(fileName)
		if err != nil {
			return 0, fmt.Errorf("Error on file creation with: %+v", err)
		}

		client := http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
		}
		// Put content on file
		resp, err := client.Get(fullURLFile)
		if err != nil {
			return 0, fmt.Errorf("Error getting file with: %+v", err)
		}
		defer resp.Body.Close()

		size, err := io.Copy(file, resp.Body)
		defer file.Close()
		if err != nil {
			return 0, fmt.Errorf("Error downloading file with: %+v", err)
		}

		return size, nil
	} else {
		return 0, fmt.Errorf("The object is not a multimedia file")
	}
}

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}
