package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var submiturl = "https://open.kattis.com/submit"
var loginurl = "https://open.kattis.com/login"
var _headerKey = "User-Agent"
var _headerValue = "kattis-cli-submit"

func Login(config Config) (*http.Response, error) {
	username := config.User.Username
	token := config.User.Token
	loginArgs := url.Values{}
	loginArgs.Set("user", username)
	loginArgs.Set("script", "true")
	loginArgs.Set("token", token)

	request, _ := http.NewRequest(http.MethodPost, loginurl, bytes.NewBufferString(loginArgs.Encode()))
	request.Header.Add(_headerKey, _headerValue)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	res, err := client.Do(request)
	return res, err
}

func Submit() {
	config, err := GetConfig()
	if err != nil {
		fmt.Println(err)
	}
	loginRes, err := Login(config)
	if err != nil {
		log.Fatal("Login failed")
	}
	fmt.Printf("Login response: %d\n", loginRes.StatusCode)
	submitRes, err := submitFilesToKattis(getProjectFilenames(config), loginRes.Cookies())
	if err != nil {
		log.Fatal("Error submitting to kattis")
	}
	fmt.Printf("Submit response: %d\n", submitRes.StatusCode)
	response := make([]byte, 128)
	submitRes.Body.Read(response)
	fmt.Println(response)
	lastLine := strings.Split(string(response), "\n")[1]
	submissionURL := strings.Split(lastLine, " ")[2]
	fmt.Println(submissionURL)
}

func getProjectFilenames(config Config) []string {
	language, err := GetLanguage(config.DefaultLang)
	if err != nil {
		log.Fatal("error getting the language")
	}
	filenames := getAllFileNames(language.Extensions)
	return filenames
}

func submitFilesToKattis(files []string, cookies []*http.Cookie) (*http.Response, error) {
	projectConfig, err := GetProjectConfig()
	if err != nil {
		log.Fatal("Could not get project config")
	}
	//Prepare form data
	data := make(map[string]string)
	data["submit"] = "true"
	data["submit_ctr"] = "2"
	data["language"] = projectConfig.Language
	data["mainclass"] = strings.TrimSuffix(projectConfig.MainFile, filepath.Ext(projectConfig.MainFile))
	data["problem"] = projectConfig.Problem
	data["tag"] = ""
	data["script"] = "true"

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add form fields
	for key, value := range data {
		if err := writer.WriteField(key, value); err != nil {
			return nil, fmt.Errorf("error writing field %s: %v", key, err)
		}
	}

	// Add files
	for _, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("error opening file %s: %v", filePath, err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile("sub_file[]", file.Name())
		if err != nil {
			return nil, fmt.Errorf("error creating form file for %s: %v", filePath, err)
		}

		if _, err = io.Copy(part, file); err != nil {
			return nil, fmt.Errorf("error copying file %s: %v", filePath, err)
		}
	}

	// Close the writer to finalize the form
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("error closing writer: %v", err)
	}

	request, _ := http.NewRequest(http.MethodPost, submiturl, &buf)
	request.Header.Add(_headerKey, _headerValue)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// Add cookies to the request
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %v", err)
	}

	return resp, nil
}

func showJudgement(submissionURL string) {
	
}
