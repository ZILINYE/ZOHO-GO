package main

import (
	"ZOHO-GO/FileProcess"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

// Define each request information
type RequestInfo struct {
	File_prefix string
	Request_ID  string
	Student_ID  string
	AccessToken string
}

var jobs = make(chan RequestInfo)

// var extactlist = make(chan string)

func AddQueue(accessToken, prefix string, downloadList [][]string) {

	for _, element := range downloadList {

		prefix := prefix + element[1] + "-"
		requestID := element[0]
		StuID := element[2]
		newitem := &RequestInfo{Student_ID: StuID, Request_ID: requestID, File_prefix: prefix, AccessToken: accessToken}

		jobs <- *newitem
	}

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
		foldername := "StudentContract"
		filename := job.File_prefix + job.Student_ID + ".zip"
		out, err := os.Create(foldername + "/" + filename)
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
	// prefix := "2023-Spring-"
	noOfWorker := 100
	// get Access Token
	// accessToken := GetList.RetriveToken()
	// DownloadList := GetList.GetDownloadList(100, "23S")

	// read download list into channel
	// go AddQueue(accessToken, prefix, DownloadList)

	// multithread worker pool and assign job from channel
	// CreateWorkerPool(noOfWorker)
	go FileProcess.LoopFile()
	FileProcess.CreateWorkerPool(noOfWorker)
	// FileProcess.MergePDFs()
	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("The download total time taken ", diff.Seconds(), "seconds")

	// Test.Test()

}
