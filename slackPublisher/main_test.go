package main

import (
	"bytes"
	"embed"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slacktest"
)

//go:embed testdata/*
var testdata embed.FS

func TestNoRequiredEnviroments(t *testing.T) {
	type env struct {
		key   string
		value string
	}

	tests := []struct {
		name        string
		expected    string
		enviroments []env
	}{
		{
			name:        "NoSlackToken",
			expected:    "No token.",
			enviroments: []env{{key: "SLACK_CHANNEL_ID", value: "channel"}},
		},
		{
			name:        "NoSlackChannelId",
			expected:    "No channel.",
			enviroments: []env{{key: "SLACK_TOKEN", value: "token"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			for _, e := range tt.enviroments {
				t.Setenv(e.key, e.value)
			}

			var buf bytes.Buffer
			log.SetOutput(&buf)
			defaultFlags := log.Flags()
			log.SetFlags(0)
			defer func() {
				log.SetOutput(os.Stderr)
				log.SetFlags(defaultFlags)
				buf.Reset()
			}()

			main()

			actual := strings.TrimRight(buf.String(), "\n")
			if tt.expected != actual {
				t.Errorf("name: %s\nexpected: %v\nactual: %v", tt.name, tt.expected, actual)
			}
		})
	}
}

func TestSendFromMessage(t *testing.T) {
	type env struct {
		key   string
		value string
	}

	type expected struct {
		postMessage int
		uploadFile  int
		log         string
	}

	tests := []struct {
		name     string
		env      []env
		api      bool
		expected expected
	}{
		{
			name:     "SendNoMessage",
			env:      []env{},
			api:      true,
			expected: expected{postMessage: 0, uploadFile: 0, log: ""},
		},
		{
			name:     "SendAsMessageIsSuccess",
			env:      []env{{key: "MESSAGE", value: strings.Repeat("a", 2999)}},
			api:      true,
			expected: expected{postMessage: 1, uploadFile: 0, log: "Successful: post message"},
		},
		{
			name:     "sendAdFileIsSuccess",
			env:      []env{{key: "MESSAGE", value: strings.Repeat("a", 3000)}},
			api:      true,
			expected: expected{postMessage: 0, uploadFile: 1, log: "Successful: upload file"},
		},
		{
			name:     "SendAsMessageIsFail",
			env:      []env{{key: "MESSAGE", value: strings.Repeat("a", 2999)}},
			api:      false,
			expected: expected{postMessage: 1, uploadFile: 0, log: "can not post message too_many_attachments"},
		},
		{
			name:     "sendAdFileIsFail",
			env:      []env{{key: "MESSAGE", value: strings.Repeat("a", 3000)}},
			api:      false,
			expected: expected{postMessage: 0, uploadFile: 1, log: "can not upload file EOF"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			var buf bytes.Buffer
			log.SetOutput(&buf)
			defaultFlags := log.Flags()
			log.SetFlags(0)
			defer func() {
				log.SetOutput(os.Stderr)
				log.SetFlags(defaultFlags)
				buf.Reset()
			}()

			var countPostMessage = 0
			var countUploadFile = 0
			ts := slacktest.NewTestServer(func(c slacktest.Customize) {
				var path string
				if tt.api {
					path += "ok.json"
				} else {
					path += "error.json"
				}
				c.Handle("/chat.postMessage", func(w http.ResponseWriter, _ *http.Request) {
					countPostMessage += 1
					res, _ := testdata.ReadFile("testdata/chatPostMessage/" + path)
					w.Write(res)
				})
				c.Handle("/files.getUploadURLExternal", func(w http.ResponseWriter, _ *http.Request) {
					countUploadFile += 1
					res, _ := testdata.ReadFile("testdata/filesGetUploadURLExternal/" + path)
					w.Write(res)
				})
				c.Handle("/files.completeUploadExternal", func(w http.ResponseWriter, _ *http.Request) {
					res, _ := testdata.ReadFile("testdata/filesCompleteUploadExternal/" + path)
					w.Write(res)
				})
				c.Handle("/dummy.uploadUrl", func(w http.ResponseWriter, _ *http.Request) {
					res, _ := testdata.ReadFile("testdata/dummyUploadUrl/" + path)
					w.Write(res)
				})
			})
			ts.Start()
			client := slack.New("testToken", slack.OptionAPIURL(ts.GetAPIURL()))
			for _, e := range tt.env {
				t.Setenv(e.key, e.value)
			}
			Client{Client: client}.sendFromMessage("channel")

			if tt.expected.postMessage != countPostMessage {
				t.Errorf("name: %s\ncount post message\nexpected: %d\nactual: %d", tt.name, tt.expected.postMessage, countPostMessage)
			}
			if tt.expected.uploadFile != countUploadFile {
				t.Errorf("name: %s\ncount upload file\nexpected: %d\nactual: %d", tt.name, tt.expected.uploadFile, countUploadFile)
			}

			actual := strings.TrimRight(buf.String(), "\n")
			if tt.expected.log != actual {
				t.Errorf("name: %s\nexpected: %s\nactual: %s", tt.name, tt.expected.log, actual)
			}
		})
	}
}

func TestSendFromFile(t *testing.T) {

	type expected struct {
		postMessage int
		uploadFile  int
		log         string
	}

	tests := []struct {
		name     string
		message  string
		api      bool
		expected expected
	}{
		{
			name:     "SendNoFile",
			message:  "",
			api:      true,
			expected: expected{postMessage: 0, uploadFile: 0, log: ""},
		},
		{
			name:     "SendAsMessageIsSuccess",
			message:  strings.Repeat("a", 2999),
			api:      true,
			expected: expected{postMessage: 1, uploadFile: 0, log: "Successful: post message"},
		},
		{
			name:     "sendAdFileIsSuccess",
			message:  strings.Repeat("a", 3000),
			api:      true,
			expected: expected{postMessage: 0, uploadFile: 1, log: "Successful: upload file"},
		},
		{
			name:     "SendAsMessageIsFail",
			message:  strings.Repeat("a", 2999),
			api:      false,
			expected: expected{postMessage: 1, uploadFile: 0, log: "can not post message too_many_attachments"},
		},
		{
			name:     "sendAdFileIsFail",
			message:  strings.Repeat("a", 3000),
			api:      false,
			expected: expected{postMessage: 0, uploadFile: 1, log: "can not upload file EOF"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Helper()

			var buf bytes.Buffer
			log.SetOutput(&buf)
			defaultFlags := log.Flags()
			log.SetFlags(0)
			defer func() {
				log.SetOutput(os.Stderr)
				log.SetFlags(defaultFlags)
				buf.Reset()
			}()

			var countPostMessage = 0
			var countUploadFile = 0
			ts := slacktest.NewTestServer(func(c slacktest.Customize) {
				var path string
				if tt.api {
					path += "ok.json"
				} else {
					path += "error.json"
				}
				c.Handle("/chat.postMessage", func(w http.ResponseWriter, _ *http.Request) {
					countPostMessage += 1
					res, _ := testdata.ReadFile("testdata/chatPostMessage/" + path)
					w.Write(res)
				})
				c.Handle("/files.getUploadURLExternal", func(w http.ResponseWriter, _ *http.Request) {
					countUploadFile += 1
					res, _ := testdata.ReadFile("testdata/filesGetUploadURLExternal/" + path)
					w.Write(res)
				})
				c.Handle("/files.completeUploadExternal", func(w http.ResponseWriter, _ *http.Request) {
					res, _ := testdata.ReadFile("testdata/filesCompleteUploadExternal/" + path)
					w.Write(res)
				})
				c.Handle("/dummy.uploadUrl", func(w http.ResponseWriter, _ *http.Request) {
					res, _ := testdata.ReadFile("testdata/dummyUploadUrl/" + path)
					w.Write(res)
				})
			})
			ts.Start()
			client := slack.New("testToken", slack.OptionAPIURL(ts.GetAPIURL()))
			if len(tt.message) > 0 {
				f, _ := os.CreateTemp("", "message.txt")
				f.WriteString(tt.message)
				t.Setenv("FILE_PATH", f.Name())
				defer func() {
					os.Remove("message.txt")
				}()
			}

			Client{Client: client}.sendFromFile("channel")

			if tt.expected.postMessage != countPostMessage {
				t.Errorf("name: %s\ncount post message\nexpected: %d\nactual: %d", tt.name, tt.expected.postMessage, countPostMessage)
			}
			if tt.expected.uploadFile != countUploadFile {
				t.Errorf("name: %s\ncount upload file\nexpected: %d\nactual: %d", tt.name, tt.expected.uploadFile, countUploadFile)
			}

			actual := strings.TrimRight(buf.String(), "\n")
			if tt.expected.log != actual {
				t.Errorf("name: %s\nexpected: %s\nactual: %s", tt.name, tt.expected.log, actual)
			}
		})
	}
}
