package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	// MakeRequest("POST", "ticket")
	// MakeRequest("GET", "ticket")
	MakeRequest("PUT", "ticket/12345678")
	MakeRequest("PUT", "status/12345678")
	// MakeRequest("GET", "ticket/12345678")
}

func MakeRequest(requestType string, requestContents string) {
	requestContents = "http://localhost:8080/" + requestContents
	fmt.Println(requestContents)

	// //Populate body of request if necessary
	// requestData := url.Values{}
	// if additionalInformation != nil {
	// 	for key, value := range additionalInformation {
	// 		requestData.Set(key, value)
	// 	}
	// }
	client := &http.Client{}
	req, err := http.NewRequest(requestType, requestContents, strings.NewReader((url.Values{}).Encode()))
	if err != nil {
		log.Fatalln("Error while creating request", err)
	}
	// fmt.Println("Request: ")
	// fmt.Println(req)
	resp, err := client.Do(req)
	// fmt.Println("Response: ")
	// fmt.Println(resp)
	if err != nil {
		log.Fatalln("Error while executing request", err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// // log.Println(resp)
	// log.Println("Body back :\n", string(body))

	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr)

	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// // log.Println(resp)
	// log.Println("Body back :\n", string(body))
}

// func main() {

// }
