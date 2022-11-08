package main

import (
	// "encoding/json"

	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

type RequestInfo struct {
	File_prefix  string
	Request_ID   string
	Student_ID   string
	Program_code string
	AccessToken  string
}

type AccessToken struct {
	Access_token string `json:"access_token"`
	Api_domain   string `json:"api_domain"`
	Token_type   string `json:"token_type"`
	Expires_in   int    `json:"expires_in"`
}

var jobs = make(chan RequestInfo)

func RetriveToken() string {

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

	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Reading body failed: %s", err)

	}

	// Unmarshal result
	post := AccessToken{}
	err = json.Unmarshal([]byte(body), &post)
	if err != nil {
		log.Printf("Reading body failed: %s", err)

	}
	return post.Access_token
}
func AddQueue(accessToken string) {
	readFile, err := os.Open("test.txt")
	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	// var fileLines []RequestInfo
	for fileScanner.Scan() {
		downloadslice := strings.Split(fileScanner.Text(), " ")
		prefix := downloadslice[0]
		requestID := downloadslice[1]
		StuID := downloadslice[2]
		newitem := &RequestInfo{Student_ID: StuID, Request_ID: requestID, File_prefix: prefix, AccessToken: accessToken}
		// fileLines = append(fileLines, *newitem)
		jobs <- *newitem
	}

	readFile.Close()
	close(jobs)
	// return &fileLines
}

// func (R RequestInfo) AddQueue(AccessToken string, c chan [5]string) {
// 	requeslist := [5]string{R.File_prefix, R.Program_code, R.Request_ID, R.Student_ID, AccessToken}

// 	c <- requeslist

// }
func CreateWokerPool(noOfWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go Downloader(&wg)
	}
	wg.Wait()
	// close(results)

}

func Downloader(wg *sync.WaitGroup) {
	for job := range jobs {
		out, err := os.Create(job.File_prefix + job.Student_ID + ".zip")
		if err != nil {
			fmt.Println(err)
		}
		defer out.Close()
		fmt.Println(job.Student_ID)
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://sign.zoho.com/api/v1/requests/"+job.Request_ID+"/pdf", nil)
		if err != nil {
			fmt.Println(err)
		}
		req.Header.Set("Authorization", "Zoho-oauthtoken "+job.AccessToken)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			fmt.Println(err)
		}

	}
	wg.Done()

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
	noOfWorker := 10
	accessToken := RetriveToken()
	// done := make(chan bool)
	go AddQueue(accessToken)

	CreateWokerPool(noOfWorker)
	fmt.Println("main function")
	// <-done
}
