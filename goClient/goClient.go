package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	serverLocation = "http://localhost:8080/"
	PUT            = "PUT"
	POST           = "POST"
	GET            = "GET"
	DELETE         = "DELETE"
)

var reader = bufio.NewReader(os.Stdin)

func main() {
	command := ""
	fmt.Println("Welcome to the lottery system!")

	//Main text loop
	for command != "6" {
		fmt.Println("\nPlease choose what you would like to do: ")
		fmt.Println("(1) Create new ticket")
		fmt.Println("(2) Return current ticket ID's")
		fmt.Println("(3) Get ticket by ID")
		fmt.Println("(4) Add a new line to a ticket by ID")
		fmt.Println("(5) Check the result of a ticket")
		fmt.Println("(6) Exit the program \n")

		command, _ = reader.ReadString('\n')
		command = strings.TrimRight(command, "\n")
		switch command {
		case "1":
			createNewTicket()
		case "2":
			getAllTickets()
		case "3":
			getSingleTicket()
		case "4":
			addLinesToTicket()
		case "5":
			checkTicketStatus()
		case "6":
		default:
			fmt.Println("Command not recognised, please try again")
		}

	}
}

func createNewTicket() {
	fmt.Println("Creating new ticket")
	ticket := MakeRequest(POST, serverLocation+"ticket") // Create new ticket
	if ticket[0].ID == "" {
		fmt.Println("Error occured, the ticket could not be created")
	} else {
		fmt.Println("Ticket created, ID : ", ticket[0].ID)
	}
}

func getAllTickets() {
	fmt.Println("Creating new ticket")
	ticket := MakeRequest(GET, serverLocation+"ticket") //Return all ticket id's
	if ticket[0].ID == "" {
		fmt.Println("Error occured, the tickets could not be retrieved")
	} else {
		fmt.Println("The following are current ticket ID's : ")
		for index, value := range ticket {
			fmt.Println(index, ". ", value.ID)
		}
	}
}

func getSingleTicket() {
	fmt.Println("Retrieving ticket")
	ticketNumber, _ := reader.ReadString('\n')
	ticket := MakeRequest(GET, serverLocation+"ticket/"+ticketNumber) // Get ticket
	if ticket[0].ID == "" {
		fmt.Println("The ticket could not be found")
	} else {
		fmt.Println("Ticket found :")
		printTickets(ticket)
	}
}

func addLinesToTicket() {
	fmt.Println("Retrieving ticket:")
	ticketNumber, _ := reader.ReadString('\n')
	fmt.Println("How many lines:")
	numLines, _ := reader.ReadString('\n')
	ticket := MakeRequest(PUT, serverLocation+"ticket/?id="+ticketNumber+"&lines="+numLines) // Add ticket line
	if ticket[0].ID == "" {
		fmt.Println("The ticket could not be found")
	} else {
		fmt.Println("Ticket found, line has been appended. Ticket ID :", ticket[0].ID)
	}
}

func checkTicketStatus() {
	fmt.Println("Retrieving ticket")
	ticketNumber, _ := reader.ReadString('\n')
	ticket := MakeRequest(PUT, serverLocation+"status/"+ticketNumber) //Calculate result and retrieve ticket
	if ticket[0].ID == "" {
		fmt.Println("The ticket could not be found")
	} else {
		fmt.Println("Ticket found, results calculated.")
		printTickets(ticket)
	}
}

func printTickets(tickets []Ticket) {
	for _, value := range tickets {
		fmt.Println("---------------")
		fmt.Println("ID				:	", value.ID)
		fmt.Println("IsChecked		:	", value.IsChecked)
		fmt.Println("Ticket Lines	: ")
		printLines(value.Lines)
		fmt.Println("---------------")
	}
}
func printLines(lines []Line) {
	for index, value := range lines {
		fmt.Println(index, ". ", "Values : ", value.Values, " Result : ", value.Result)
	}
}

func MakeRequest(requestType string, requestContents string) []Ticket {

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
	var tickets []Ticket
	err = json.Unmarshal([]byte(body), &tickets)
	fmt.Println("Body : ", string(body))
	if err != nil {
		fmt.Println("Error while unmarshalling", err)
	}
	return tickets
}

//Objects
type Ticket struct {
	ID        string `json:"id,omitempty"`
	Lines     []Line `json:"lines,omitempty"`
	IsChecked bool   `json:"isChecked,omitempty"`
}

type Line struct {
	ID     string `json:"id,omitempty"`
	Values []int  `json:"values,omitempty"`
	Result *int   `json:"result,omitempty"`
}
