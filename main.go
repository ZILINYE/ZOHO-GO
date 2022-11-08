package main

import (
	// "encoding/json"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type RequestInfo struct {
	File_prefix  string
	Request_ID   string
	Student_ID   string
	Program_code string
}

type AccessToken struct {
	Access_token string `json:"access_token"`
	Api_domain   string `json:"api_domain"`
	Token_type   string `json:"token_type"`
	Expires_in   int    `json:"expires_in"`
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
	RetriveToken()

	// readFile, err := os.Open("test.txt")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fileScanner := bufio.NewScanner(readFile)
	// fileScanner.Split(bufio.ScanLines)
	// var fileLines []RequestInfo
	// for fileScanner.Scan() {
	// 	downloadslice := strings.Split(fileScanner.Text(), " ")
	// 	prefix := downloadslice[0]
	// 	requestID := downloadslice[1]
	// 	StuID := downloadslice[2]
	// 	newitem := &RequestInfo{student_ID: StuID, request_ID: requestID, file_prefix: prefix}
	// 	fileLines = append(fileLines, *newitem)
	// }

	// readFile.Close()

	// for _, line := range fileLines {
	// 	line.Downloader()

	// }

}

func RetriveToken() {

	params := url.Values{}
	params.Add("refresh_token", "1000.23be290456580cd7378b94f2eb3d2334.8ed115d741371a6cf5ada13b2903819e")
	params.Add("client_id", "1000.6C4D4C3LQS1XV9BVF70PS55G3PELTK")
	params.Add("client_secret", "211e6b7d3395fd9e8f7d67df464a884e0f573c6079")
	params.Add("redirect_uri", "https%3A%2F%2Fsign.zoho.com")
	params.Add("grant_type", "refresh_token")
	resp, err := http.PostForm("https://accounts.zoho.com/oauth/v2/token?",
		params)
	if err != nil {
		log.Printf("Request Failed: %s", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Printf("Reading body failed: %s", err)
	// 	return
	// }

	// Log the request body
	bodyString := string(body)
	log.Print(bodyString)
	// Unmarshal result
	post := AccessToken{}
	err = json.Unmarshal([]byte(body), &post)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return
	}
	fmt.Println(post.Access_token)
}
func (R RequestInfo) Downloader() {
	fmt.Println(R.File_prefix)
}
