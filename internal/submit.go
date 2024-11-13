package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var _statusMap = map[int]string{
	0: "New",
	1: "New",
	2: "Waiting for compile",
	3: "Compiling",
	4: "Waiting for run",
	5: "Running",
	6: "Judge Error",
	// 7: "<invalid value>",
	8:  "Compile Error",
	9:  "Run Time Error",
	10: "Memory Limit Exceeded",
	11: "Output Limit Exceeded",
	12: "Time Limit Exceeded",
	13: "Illegal Function",
	14: "Wrong Answer",
	// 15: "<invalid value>",
	16: "Accepted",
}

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
	projectConfig, err := GetProjectConfig()
	if err != nil {
		HandleError("error getting project config", err)
	}
	loginRes, err := Login(config)

	if err != nil {
		HandleError("error logging into kattis", err)
	}

	submitRes, err := submitFilesToKattis(getProjectFilenames(config), loginRes.Cookies(), projectConfig)

	if err != nil {
		HandleError("error submitting to kattis", err)
	}

	response, err := io.ReadAll(submitRes.Body)

	if err != nil {
		HandleError("error trying to read submit response", err)
	}

	lastLine := strings.Split(string(response), "\n")[1]
	submissionURL := strings.Trim(strings.Split(lastLine, " ")[2], "\x00")

	showJudgement(submissionURL, loginRes.Cookies(), projectConfig)

}

func getProjectFilenames(config Config) []string {
	language, err := GetLanguage(config.DefaultLang)
	if err != nil {
		HandleError("error getting the language", err)
	}
	filenames := getAllFileNames(language.Extensions)
	return filenames
}

func submitFilesToKattis(files []string, cookies []*http.Cookie, projectConfig ProjectConfig) (*http.Response, error) {

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
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				HandleError("error closing file", err)
			}
		}(file)

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

func getSubmissionStatus(submissionURL string, cookies []*http.Cookie) (map[string]interface{}, error) {

	request, err := http.NewRequest(http.MethodGet, submissionURL+"?json", bytes.NewBufferString(url.Values{}.Encode()))

	if err != nil {
		HandleError("could not generate http request", err)
	}

	request.Header.Set(_headerKey, _headerValue)
	request.Header.Set("Accept", "application/json")

	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)

	result := make(map[string]interface{})

	err = json.Unmarshal(body, &result)
	if err != nil {
		HandleError("error putting result into map", err)
	}

	return result, nil
}

func showJudgement(submissionURL string, loginCookies []*http.Cookie, projectConfig ProjectConfig) {
	updatingProgressbar := false
	for {

		status, err := getSubmissionStatus(submissionURL, loginCookies)
		if err != nil {
			log.Fatal("Could not retrieve submission status")
		}

		status_id, err := getIntFromJsonField(status, "status_id")
		if err != nil {
			HandleError("could not get status_id from map", err)
		}

		testcase_index, err := getIntFromJsonField(status, "testcase_index")
		if err != nil {
			HandleError("could not get testcase_index from map", err)
		}

		htmlCode, err := getHtmlCodeFromJsonField(status, "row_html")

		if err != nil {
			HandleError("Could not get htmlcode from field row_html", err)
		}

		testcase_total, err := getTotalTestCaseAmount(htmlCode)

		printSubmissionProgressBar(projectConfig, status_id, testcase_index, htmlCode, testcase_total, updatingProgressbar)
		updatingProgressbar = true
		// stop if submission has finished running on kattis
		if status_id > 5 {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func printSubmissionProgressBar(projectConfig ProjectConfig, status_id int, testcase_index int, htmlCode string, testcase_total int, updatingProgressbar bool) {
	whiteColor := "\033[0m"
	redColor := "\033[0;31m"
	greenColor := "\033[0;32m"
	yellowColor := "\033[0;33m"
	if updatingProgressbar {
		fmt.Printf("\033[F\033[F\033[F\033[F") // Move up 3 lines (one for each line to be updated)
		fmt.Printf("\r\033[K")                 // Move to start of the line and clear it
	}
	fmt.Printf("Problem name: %s\n", projectConfig.Problem)
	fmt.Printf("\r\033[K") // Move to start of the line and clear it
	var color string
	if status_id < 5 {
		color = whiteColor
	} else if status_id == 5 {
		color = yellowColor
	} else if status_id < 16 {
		color = redColor
	} else {
		color = greenColor
	}
	fmt.Printf("Status: %s%s%s\n", color, _statusMap[status_id], whiteColor)
	fmt.Printf("\r\033[K") // Move to start of the line and clear it
	fmt.Printf("Succeded: %d/%d\n", countAcceptedTestCases(htmlCode), testcase_total)

	wrongCharacter := "X"
	rightCharacter := "#"
	runningCharacter := "-"

	testString := ""
	for index := 1; index < testcase_total; index++ {
		status := getTestCaseStatus(htmlCode, index)
		if status == "not checked" {
			testString += runningCharacter
		} else if status == "Wrong Answer" {
			testString += redColor + wrongCharacter + whiteColor
		} else {
			testString += greenColor + rightCharacter + whiteColor
		}

	}
	fmt.Printf("[%s] %d/%d\n", testString, testcase_index, testcase_total)
}

func getTestCaseStatus(htmlCode string, testCaseIndex int) string {
	re := regexp.MustCompile(fmt.Sprintf("Test case %d/[0-9]*:\\s(\\w+ ?\\w*|\\w*)", testCaseIndex))
	return re.FindStringSubmatch(htmlCode)[1]
}

func countAcceptedTestCases(htmlCode string) int {
	re := regexp.MustCompile("Test case [0-9]*/[0-9]*: Accepted")
	return len(re.FindAll([]byte(htmlCode), -1))
}

func getIntFromJsonField(status map[string]interface{}, field string) (int, error) {

	if status_id_value, ok := status[field].(float64); ok {
		// Convert float64 to int
		return int(status_id_value), nil
	} else {
		return 0, errors.New("could not convert field to int")

	}
}

func getHtmlCodeFromJsonField(status map[string]interface{}, field string) (string, error) {

	if status_id_value, ok := status[field].(string); ok {
		// Convert float64 to int
		return string(status_id_value), nil
	} else {
		return "", errors.New("could not convert field to string")

	}
}

func getTotalTestCaseAmount(htmlCode string) (int, error) {
	testcase_amount := strings.Count(htmlCode, "Test case")
	return testcase_amount, nil
}
