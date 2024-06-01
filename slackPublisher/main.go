package main

import (
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack"
)

type Client struct {
	*slack.Client
}

func main() {
	token, isTokenPresent := os.LookupEnv("SLACK_TOKEN")
	if !isTokenPresent {
		log.Fatalln("No token.")
		return
	}
	channel, isChannelPresent := os.LookupEnv("SLACK_CHANNEL_ID")
	if !isChannelPresent {
		log.Fatalln("No channel.")
		return
	}
	client := Client{slack.New(token)}
	message, isMessagePresent := os.LookupEnv("MESSAGE")
	boader := 3000
	if isMessagePresent {
		l := len(message)
		if l < boader {
			_, _, e := client.PostMessage(channel, slack.MsgOptionText(message, false))
			if e != nil {
				log.Fatalln("can not upload file", e)
			}
			log.Println("Successful: post message")
		} else {
			f, err := os.CreateTemp("", "message.txt")
			if err != nil {
				log.Fatalln("can not create temp file", err)
				return
			}
			f.WriteString(message)
			_, e := client.UploadFileV2(slack.UploadFileV2Parameters{File: f.Name(), Channel: channel, InitialComment: f.Name(), Title: f.Name(), Filename: f.Name(), FileSize: l})
			if e != nil {
				log.Fatalln("can not upload file", e)
			}
			log.Println("Successful: upload file")
		}
	}
	filePath, isFilePresent := os.LookupEnv("FILE_PATH")
	if isFilePresent {
		text, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalln("can not read file", err)
			return
		}
		l := len(text)
		if l < boader {
			_, _, e := client.PostMessage(channel, slack.MsgOptionText(string(text), false))
			if e != nil {
				log.Fatalln("can not post message", e)
			}
			log.Println("Successful: post message")
		} else {
			arrays := strings.Split(filePath, "/")
			name := arrays[len(arrays)-1]
			_, e := client.UploadFileV2(slack.UploadFileV2Parameters{File: filePath, Channel: channel, InitialComment: name, Title: name, Filename: name, FileSize: l})
			if e != nil {
				log.Fatalln("can not upload file", e)
			}
			log.Println("Successful: upload file")
		}
	}
}
