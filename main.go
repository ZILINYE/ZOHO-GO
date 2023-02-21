package main

import (
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
	"time"
)

type RequestInfo struct {
	File_prefix string
	Request_ID  string
	Student_ID  string
	AccessToken string
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
	params.Add("refresh_token", "") # Refresh Token from ZOHO
	params.Add("client_id", "") # Client ID From ZOHO
	params.Add("client_secret", "") # Client Secret
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
	readFile, err := os.Open("DownloadList.txt")
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
}

func CreateWokerPool(noOfWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go Processor(&wg)
	}
	wg.Wait()

}

func Processor(wg *sync.WaitGroup) {
	for job := range jobs {
		out, err := os.Create(job.File_prefix + job.Student_ID + ".zip")
		if err != nil {
			fmt.Println(err)
		}

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

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			fmt.Println(err)
		}

	}
	wg.Done()

}

func main() {

	startTime := time.Now()
	noOfWorker := 5
	// get Access Token
	accessToken := RetriveToken()

	// read download list into channel
	go AddQueue(accessToken)

	// multithread worker pool and assign job from channel
	CreateWokerPool(noOfWorker)
	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("total time taken ", diff.Seconds(), "seconds")

}
