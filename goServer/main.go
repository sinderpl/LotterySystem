package main

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

//variables
const numOfValues = 3

var existingTickets []Ticket

//Main runner
func main() {
	//Test ticket
	testTicket := generateTicket()
	testTicket.ID = 12345678
	existingTickets = append(existingTickets, testTicket)

	router := mux.NewRouter()
	router.HandleFunc("/ticket", CreateTicket).Methods("POST")
	router.HandleFunc("/ticket", GetTickets).Methods("GET")        //can this be merged with retrieve all
	router.HandleFunc("/ticket/{id}", GetTicket).Methods("GET")    // add method to retrieve multiple tickets by id
	router.HandleFunc("/ticket/{id}", UpdateTicket).Methods("PUT") // add n for number of lines
	router.HandleFunc("/status/{id}", GetTicketStatus).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", router))

}

//Methods

//Generates a new ticket with single line and returns ID to user
func CreateTicket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("createTicketCalled")

	ticket := generateTicket()
	existingTickets = append(existingTickets, ticket)
	json.NewEncoder(w).Encode(ticket.ID)
}

//Returns all ticket ID's
func GetTickets(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getTicketsCalled")

	var ticketList []int
	if existingTickets != nil {
		for _, value := range existingTickets {
			ticketList = append(ticketList, value.ID)
		}
	}

	json.NewEncoder(w).Encode(ticketList)
}

//Return the specified ticket
func GetTicket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getTicketCalled")

	fmt.Println(r.GetBody())
}
func UpdateTicket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("updateTicket")
}
func GetTicketStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getTicketStatus")
}

func generateTicket() Ticket {
	lines := make([]Line, 1)
	lines[0] = generateTicketLine()
	newTicket := Ticket{ID: rand.Intn(999999), Lines: lines}
	return newTicket
}

func generateTicketLine() Line {
	values := make([]int, numOfValues)
	for index, _ := range values {
		values[index] = rand.Intn(3) // 0 - 2 inclusive
	}
	newLine := Line{ID: rand.Intn(999999), Values: values}
	return newLine
}

func calculateResult() {
	fmt.Println("Generating results")
}

//Structs
type Ticket struct {
	ID        int    `json:"id,omitempty"`
	Lines     []Line `json:"Lines,omitempty"`
	IsChecked bool   `json:"IsChecked,omitempty"`
}

type Line struct {
	ID     int   `json:"id,omitempty"`
	Values []int `json:"values,omitempty"`
}
