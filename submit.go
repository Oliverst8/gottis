package main

import (
	"bytes"
	"fmt"
	"net/http"
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
	body := []byte(fmt.Sprintf(`{
		"user": %s,
		"token": %s,
		"script": "true"
	}`, username, token))
	request, err := http.NewRequest(http.MethodPost, loginurl, bytes.NewReader(body))
	request.Header.Add(_headerKey, _headerValue)
	client := &http.Client{}
	res, err := client.Do(request)
	fmt.Println(res.Cookies())
	return res
}

func Submit() {
	Login()

}
