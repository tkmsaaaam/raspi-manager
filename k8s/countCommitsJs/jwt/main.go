package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Request struct {
	Repositories []string `json:"repositories"`
}

type Response struct {
	Token               string
	ExpiresAt           time.Time
	Permissions         Permissions
	RepositorySelection string
}

type Permissions struct {
	Metadata string
}

func main() {
	appId := os.Getenv("APP_ID")
	installId := os.Getenv("INSTALL_ID")
	if appId == "" {
		log.Println("APP_ID is required")
		return
	}
	if installId == "" {
		log.Println("INSTALL_ID is required")
		return
	}

	var d = time.Minute * 10
	if durationStr := os.Getenv("DURATION"); durationStr != "" {
		duration, e := time.ParseDuration(durationStr)
		if e == nil {
			d = duration
		}
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(d).Unix(),
		"iss": appId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	var keyFile = "private-key.pem"
	if path := os.Getenv("KEY_FILE_PATH"); path != "" {
		keyFile = path
	}

	keyRow, _ := os.ReadFile(keyFile)
	key, _ := jwt.ParseRSAPrivateKeyFromPEM(keyRow)
	tokenString, err := token.SignedString(key)
	if err != nil {
		log.Println(err)
		return
	}

	url := "https://api.github.com/app/installations/" + installId + "/access_tokens"
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	req.Header.Set("Accept", "application/vnd.github+json")
	client := new(http.Client)
	resp, err := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	response := &Response{}
	jsonErr := json.Unmarshal(body, response)
	if err != nil {
		log.Println(jsonErr)
		return
	}
	if output := os.Getenv("OUTPUT_TO_FILE"); output == "true" {
		f, err := os.Create("result.json")
		if err == nil {
			body, e := json.Marshal(response)
			if e == nil {
				f.Write(body)
			}
		}
	}
	fmt.Println(response.Token)
}
