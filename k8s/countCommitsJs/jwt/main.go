package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

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
	if appId == "" || installId == "" {
		return
	}
	claims := jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * 10).Unix(),
		"iss": appId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	keyRow, _ := os.ReadFile("private-key.pem")
	key, _ := jwt.ParseRSAPrivateKeyFromPEM(keyRow)
	tokenString, err := token.SignedString(key)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(jsonErr)
	}
	fmt.Println(response.Token)
}
