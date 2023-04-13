package GetList

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type AccessToken struct {
	Access_token string `json:"access_token"`
	Api_domain   string `json:"api_domain"`
	Token_type   string `json:"token_type"`
	Expires_in   int    `json:"expires_in"`
}

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

func HttpRequest(accesstoken string) {
	// Set the URL endpoint
	apiurl := "https://sign.zoho.com/api/v1/requests"

	// Set the request parameters
	params := url.Values{}
	params.Set("data", `{"page_context":{"row_count":10,"start_index":1,"search_columns":{},"sort_column":"created_time","sort_order":"DESC"}}`)

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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Print the response body
	fmt.Println(string(body))
}
