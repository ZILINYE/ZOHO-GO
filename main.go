package main

import (
	// "encoding/json"

	"bufio"
	"fmt"
	"os"
	"strings"
)

type RequestInfo struct {
	file_prefix  string
	request_ID   string
	student_ID   string
	Program_code string
}

func main() {
	// client := &http.Client{}
	// payload := url.Values{}
	// payload.Add("data", "{'page_context': {'row_count':10 , 'start_index': 1, 'search_columns': {'request_name': '22F'}, 'sort_column': 'created_time', 'sort_order': 'DESC'}}")
	// req, _ := http.NewRequest("GET", "https://sign.zoho.com/api/v1/requests?"+payload.Encode(), nil)
	// req.Header.Add("Authorization", "Zoho-oauthtoken 1000.86165fe89600a0e4c54d924790cdfbf0.15ba5bc4b94b00b5fe46737c25fc302f")
	// response, err := client.Do(req)
	// if err != nil {
	// 	fmt.Println("wrong")
	// }
	// defer response.Body.Close()
	// body, _ := ioutil.ReadAll(response.Body) // response body is []byte
	// result := gjson.Get(string(body), `requests.#(request_status="completed").`)

	// // println(result.String())
	// for _, s := range result.Array() {
	// 	fmt.Println(s)
	// }

	readFile, err := os.Open("test.txt")
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []RequestInfo
	for fileScanner.Scan() {
		downloadslice := strings.Split(fileScanner.Text(), " ")
		prefix := downloadslice[0]
		requestID := downloadslice[1]
		StuID := downloadslice[2]
		newitem := &RequestInfo{student_ID: StuID, request_ID: requestID, file_prefix: prefix}
		fileLines = append(fileLines, *newitem)
	}

	readFile.Close()

	for _, line := range fileLines {
		line.Downloader()

	}

	// fmt.Println(fileLines)

}

func (R RequestInfo) Downloader() {
	fmt.Println(R.file_prefix)
}
