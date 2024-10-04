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
)

var submiturl = "https://open.kattis.com/submit"
var loginurl = "https://open.kattis.com/login"
var _headerKey = "User-Agent"
var _headerValue = "kattis-cli-submit"

func Login(config Config) *http.Response {
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
	res, _ := client.Do(request)
	fmt.Println(res.StatusCode)
	return res
}

func Submit() {
	config, err := GetConfig()
	if err != nil {
		fmt.Println(err)
	}
	response := Login(config)
	//submitFilesToKattis(getProjectFilenames(config), response.Cookies())
	Submit2(submiturl, response.Cookies(), "twosum", "java", "Twosum.java", "", nil, nil, getProjectFilenames(config))
}

func getProjectFilenames(config Config) []string {
	language, err := GetLanguage(config.DefaultLang)
	if err != nil {
		log.Fatal("error getting the language")
	}
	filenames := getAllFileNames(language.Extensions)
	return filenames
}

func submitFilesToKattis(files []string, cookies []*http.Cookie) *http.Response {
	args := url.Values{}
	request, _ := http.NewRequest(http.MethodPost, submiturl, bytes.NewBufferString(args.Encode()))
	request.Header.Add(_headerKey, _headerValue)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return nil
}

// Submit makes a submission to the given submit_url.
func Submit2(submitURL string, cookies []*http.Cookie, problem, language, mainClass, tag string, assignment, contest *string, files []string) (*http.Response, error) {
	// Prepare form data
	data := make(map[string]string)
	data["submit"] = "true"
	data["submit_ctr"] = "2"
	data["language"] = language
	data["mainclass"] = mainClass
	data["problem"] = problem
	data["tag"] = tag
	data["script"] = "true"

	if assignment != nil {
		data["assignment"] = *assignment
	}
	if contest != nil {
		data["contest"] = *contest
	}

	// Create a buffer to hold our form data
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

	// Create the HTTP request
	req, err := http.NewRequest("POST", submitURL, &buf)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Add cookies to the request
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %v", err)
	}

	return resp, nil
}
