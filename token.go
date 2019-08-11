package main

import "encoding/json"
import "fmt"
import "io/ioutil"
import "log"
import "net/http"
import "net/url"
import "os"
import "strings"
import "time"

var myClient = &http.Client{Timeout: 10 * time.Second}

type Token struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
	ExpiresIn    int      `json:"expires_in"`
	ExpireIn     int      `json:"expire_in"`
}

type TokenFile struct {
	Token     *Token    `json:"token"`
	Refreshed time.Time `json:"refreshed"`
}

func readCacheToken(tokenFilePath string) *TokenFile {
	jsonFile, err := os.Open(tokenFilePath)
	if err != nil {
		log.Panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var tokenFile TokenFile
	err = json.Unmarshal(byteValue, &tokenFile)
	if err != nil {
		log.Panic(err)
	}

	return &tokenFile
}

func writeCacheToken(tokenFilePath string, r *Token) {
	file, err := json.MarshalIndent(TokenFile{
		Token:     r,
		Refreshed: time.Now(),
	}, "", " ")
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile(tokenFilePath, file, 0600)
	if err != nil {
		log.Panic(err)
	}
}

func doRefresh(refreshToken string) *Token {
	reqBody := fmt.Sprintf(
		"grant_type=refresh_token&refresh_token=%s&client_id=%s&client_secret=%s", url.QueryEscape(refreshToken), clientId, clientSecret)

	url := "https://app.velux-active.com/oauth2/token"
	req, err := http.NewRequest("POST", url, strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if Debug {
		log.Printf("token refresh: %+v", req)
	}

	resp, err := myClient.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	r := new(Token)
	err = json.Unmarshal(body, &r)
	if err != nil {
		log.Panic(err)
	}

	if r.AccessToken == "" {
		log.Panicf("invalid response: %s", body)
	}

	return r
}

func refreshToken(tokenFilePath string) *Token {
	tokenFile := readCacheToken(tokenFilePath)

	var resultToken *Token

	expireTime := tokenFile.Refreshed.Add(time.Second * time.Duration(tokenFile.Token.ExpireIn))
	if expireTime.Before(time.Now()) {
		resultToken = doRefresh(tokenFile.Token.RefreshToken)
		writeCacheToken(tokenFilePath, resultToken)
	} else {
		if Debug {
			log.Println("skip refreshing token")
		}
		resultToken = tokenFile.Token
	}

	return resultToken
}
