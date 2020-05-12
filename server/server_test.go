package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	serverLocation = "http://localhost:8080/"
)

func TestMain(m *testing.M) {
	main()
	retCode := m.Run()
	ShutDown()
	os.Exit(retCode)
}

//Tests the creation of a new ticket
func TestCreateTicket(t *testing.T) {
	createNewTicket(t)
}

//Tests the get all tickets function
func TestGetAllTickets(t *testing.T) {
	//Create two tickets on the server
	tickets := make([]*Ticket, 2)
	tickets = append(tickets, createNewTicket(t))
	tickets = append(tickets, createNewTicket(t))

	//Check the amount of tickets
	ticketsOnServer := makeRequest(GET, serverLocation+"ticket", t) //Return all ticket id's
	if len(tickets) != len(ticketsOnServer) {
		log.Fatalln("The amount of tickets does not match the amount created")
		t.Fail()
	}
	//Verify ID's
	for index, _ := range tickets {
		if tickets[index].ID != ticketsOnServer[index].ID {
			log.Fatalln("Ticket ID's do not match")
			t.Fail()
		}
	}
}

//Test retrieving a single ticket
func TestGetSingleTicket(t *testing.T) {
	ticket := createNewTicket(t)
	ticketOnServer := makeRequest(GET, serverLocation+"ticket/"+ticket.ID, t) // Get ticket
	if ticket != &ticketOnServer[0] {
		log.Fatalln("Tickets do not match")
		t.Fail()
	}
}

func createNewTicket(t *testing.T) *Ticket {
	fmt.Println("Creating new Ticket")
	ticket := makeRequest(POST, serverLocation+"ticket", t) // Create new ticket
	if len(ticket) >= 1 && ticket[0].ID != "" {
		fmt.Println("Ticket created, ID : ", ticket[0].ID)
		return &ticket[0]
	} else {
		log.Fatalln("Ticker was not created correctly")
		t.Fail()
		return nil // will never be reached
	}
}

func makeRequest(requestType string, requestContents string, t *testing.T) []Ticket {

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(requestType, requestContents, strings.NewReader((url.Values{}).Encode()))
	if err != nil {
		log.Fatalln("Error while creating request: ", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Error while executing request: ", err)
		t.Fail()
	}
	defer resp.Body.Close()
	var tickets []Ticket
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 500 {
		err = json.Unmarshal([]byte(body), &tickets)
		if err != nil {
			log.Fatalln("Error while unmarshalling: ", err)
			t.Fail()
		}
	} else {
		//Print any errors
		log.Fatalln(string(body))
		t.Fail()
	}
	return tickets
}
