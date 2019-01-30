package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	MakeRequest("ticket")
}

func MakeRequest(requestContents string) {
	requestContents = "http://localhost:8080/" + requestContents
	resp, err := http.Get(requestContents)

	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// log.Println(resp)
	log.Println("Body back :\n", string(body))
}
