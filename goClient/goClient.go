package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func main() {
	// MakeRequest("POST", "ticket")
	// MakeRequest("GET", "ticket")
	MakeRequest("PUT", "ticket/12345678")
	MakeRequest("PUT", "ticket/12345678")
	MakeRequest("PUT", "ticket/12345678")
	MakeRequest("PUT", "ticket/12345678")
	MakeRequest("PUT", "ticket/12345678")
	MakeRequest("PUT", "ticket/12345678")
	MakeRequest("PUT", "ticket/12345678")
	MakeRequest("PUT", "ticket/12345678")
	MakeRequest("PUT", "status/12345678")
	// MakeRequest("GET", "ticket/12345678")
}

func MakeRequest(requestType string, requestContents string) {
	requestContents = "http://localhost:8080/" + requestContents
	fmt.Println(requestContents)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(requestType, requestContents, strings.NewReader((url.Values{}).Encode()))
	if err != nil {
		log.Fatalln("Error while creating request", err)
	}

	resp, err := client.Do(req)

	if err != nil { //Try - Catch - Finally block
		log.Fatalln("Error while executing request", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
