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

const boader = 3000

func main() {
	token, isTokenPresent := os.LookupEnv("SLACK_TOKEN")
	if !isTokenPresent {
		log.Println("No token.")
		return
	}
	channel, isChannelPresent := os.LookupEnv("SLACK_CHANNEL_ID")
	if !isChannelPresent {
		log.Println("No channel.")
		return
	}
	client := Client{slack.New(token)}
	client.sendFromMessage(channel)
	client.sendFromFile(channel)
}

func (client Client) sendFromMessage(channel string) {
	message, isMessagePresent := os.LookupEnv("MESSAGE")
	if isMessagePresent {
		l := len(message)
		if l < boader {
			_, _, e := client.PostMessage(channel, slack.MsgOptionText(message, false))
			if e != nil {
				log.Println("can not post message", e)
				return
			}
			log.Println("Successful: post message")
		} else {
			f, err := os.CreateTemp("", "message.txt")
			if err != nil {
				log.Println("can not create temp file", err)
				return
			}
			f.WriteString(message)
			_, e := client.UploadFileV2(slack.UploadFileV2Parameters{File: f.Name(), Channel: channel, InitialComment: f.Name(), Title: f.Name(), Filename: f.Name(), FileSize: l})
			if e != nil {
				log.Println("can not upload file", e)
				return
			}
			log.Println("Successful: upload file")
		}
	}
}

func (client Client) sendFromFile(channel string) {
	filePath, isFilePresent := os.LookupEnv("FILE_PATH")
	if isFilePresent {
		text, err := os.ReadFile(filePath)
		if err != nil {
			log.Println("can not read file", err)
			return
		}
		l := len(text)
		if l < boader {
			_, _, e := client.PostMessage(channel, slack.MsgOptionText(string(text), false))
			if e != nil {
				log.Println("can not post message", e)
				return
			}
			log.Println("Successful: post message")
		} else {
			arrays := strings.Split(filePath, "/")
			name := arrays[len(arrays)-1]
			_, e := client.UploadFileV2(slack.UploadFileV2Parameters{File: filePath, Channel: channel, InitialComment: name, Title: name, Filename: name, FileSize: l})
			if e != nil {
				log.Println("can not upload file", e)
				return
			}
			log.Println("Successful: upload file")
		}
	}
}
