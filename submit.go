package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

var loginurl = "https://open.kattis.com/login"
var _headerKey = "User-Agent"
var _headerValue = "kattis-cli-submit"

func Login() *http.Response {
	config, err := GetConfig()
	if err != nil {
		fmt.Println(err)
	}
	username := config.User.Username
	token := config.User.Token
	loginArgs := url.Values{}
	loginArgs.Set("user", username)
	loginArgs.Set("script", "true")
	loginArgs.Set("token", token)

	request, err := http.NewRequest(http.MethodPost, loginurl, bytes.NewBufferString(loginArgs.Encode()))
	request.Header.Add(_headerKey, _headerValue)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	res, err := client.Do(request)
	fmt.Println(res.StatusCode)
	return res
}

func Submit() {
	Login()

}
