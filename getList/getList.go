package GetList

import (
	"ZOHO-GO/Maria"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

type AccessToken struct {
	Access_token string `json:"access_token"`
	Api_domain   string `json:"api_domain"`
	Token_type   string `json:"token_type"`
	Expires_in   int    `json:"expires_in"`
}

type DownloadList struct {
	request_id   string
	program_code string
	student_id   string
}

var Download_list [][]string

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
	body, err := io.ReadAll(resp.Body)
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

func HttpRequest(search_key string, row_count, start_index int) map[string]any {
	accesstoken := RetriveToken()
	// Set the URL endpoint
	apiurl := "https://sign.zoho.com/api/v1/requests"

	// Set the request parameters
	params := url.Values{}
	content := `{"page_context":{"row_count":` + strconv.Itoa(row_count) + `,"start_index":` + strconv.Itoa(start_index) + `,"search_columns":{"request_name": ` + search_key + `},"sort_column":"created_time","sort_order":"DESC"}}`
	params.Set("data", content)

	// Create a new GET request with the authorization header
	req, err := http.NewRequest("GET", apiurl, nil)
	if err != nil {
		panic(err)
	}
	zohotoken := "Zoho-oauthtoken " + accesstoken
	req.Header.Set("Authorization", zohotoken)

	// Add the request parameters to the URL query string
	req.URL.RawQuery = params.Encode()

	// Send the request and check for errors
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Convert the resbonse into json
	m := map[string]any{}
	if err := json.Unmarshal(body, &m); err != nil {
		panic(err)
	}
	return m

}

func GetThreadnumber() float64 {
	page_context := HttpRequest("23W", 10, 1)["page_context"]

	var total_count float64 = page_context.(map[string]interface{})["total_count"].(float64)
	pages := total_count / 10

	thread_number := math.Ceil(pages)
	// fmt.Println(thread_number)
	return thread_number
}

func GetDownloadList(row_count int, keyword string) [][]string {
	// thread := GetThreadnumber()
	db := Maria.InitMaria()
	for i := 1; i <= 1; i++ {
		requests := HttpRequest(keyword, row_count, i)["requests"]
		for _, element := range requests.([]interface{}) {
			if element.(map[string]interface{})["request_status"].(string) == "completed" {
				request_id := element.(map[string]interface{})["request_id"].(string)
				camp_email, _ := element.(map[string]interface{})["actions"].([]interface{})[0].(map[string]interface{})["recipient_email"].(string)

				studentID, Proname := Maria.GetStudentInfo(camp_email, "Spring", 2023, db)
				// fmt.Print(request_id, " ", Proname, " ", studentID)
				// fmt.Print("\n")
				Download_list = append(Download_list, []string{request_id, Proname, studentID})

			}

		}
	}

	return Download_list

}
