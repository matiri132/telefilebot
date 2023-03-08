package main

import (
	"fmt"
	fu "github.com/matiri132/telefilebot/pkg/fileutils"
	"github.com/mymmrac/telego"
	"os"
)

func main() {
	botToken := os.Getenv("TOKEN")
	baseBotUrl := os.Getenv("BOT_BASE_URL")
	// Note: Please keep in mind that default logger may expose sensitive information,
	// use in development only
	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Get updates channel
	// (more on configuration in examples/updates_long_polling/main.go)
	updates, _ := bot.UpdatesViaLongPolling(nil)
	// Stop reviving updates from update channel
	defer bot.StopLongPolling()

	// Loop through all updates when they came
	for update := range updates {
		fmt.Printf("Update: %+v\n\n", update.Message)
		msg_file, err := fu.GetMsgFile(update.Message)
		msg_data, err1 := msg_file.GetData(bot)
		if err == nil && err1 == nil {
			if msg_file.Type != "text" {
				if size, err := fu.DownloadFile(&msg_file, botToken, baseBotUrl, msg_data); err != nil {
					fmt.Printf("Error in DownloadFile with: %+v\n", err)
				} else {
					fmt.Printf("Downloaded a file %s with size %d\n", msg_data, size)
				}
			} else {
				fmt.Printf("Text: %+v\n", msg_data)
			}
		} else {
			if err == nil {
				fmt.Printf("Error getting message metadata: %+v\n", err)
			} else if err1 == nil {
				fmt.Printf("Error getting multimedia file path: %+v\n", err)
			}
		}
	}
}
