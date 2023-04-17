package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"ZOHO-GO/GetList"
)

// Define each request information
type RequestInfo struct {
	File_prefix string
	Request_ID  string
	Student_ID  string
	AccessToken string
}

var jobs = make(chan RequestInfo)

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

// Define worker pool
func CreateWorkerPool(noOfWorkers int) {
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
	// // Init the database connection
	// // db := Maria.InitMaria()
	// // Getting the student information by passing the student campus email address
	// // ? param campus email address
	// Maria.GetStudentInfo("zhonghao.gu01@stclairconnect.ca", db)
	// startTime := time.Now()
	// noOfWorker := 100
	// get Access Token
	accessToken := GetList.RetriveToken()
	fmt.Println(GetList.HttpRequest(accessToken, "{}", 10))

	// read download list into channel
	// go AddQueue(accessToken)

	// // multithread worker pool and assign job from channel
	// CreateWorkerPool(noOfWorker)
	// endTime := time.Now()
	// diff := endTime.Sub(startTime)
	// fmt.Println("total time taken ", diff.Seconds(), "seconds")

}
