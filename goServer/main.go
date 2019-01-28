package main

import (
	// "encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

//Main runner
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/ticket", CreateTicket).Methods("POST")
	router.HandleFunc("/ticket", GetTickets).Methods("GET")
	router.HandleFunc("/ticket/{id}", GetTicket).Methods("GET")    // add method to retrieve multiple tickets by id
	router.HandleFunc("/ticket/{id}", UpdateTicket).Methods("PUT") // add n for number of lines
	router.HandleFunc("/status/{id}", GetTicketStatus).Methods("PUT")

	log.Fatal(http.ListenAndServe(":8080", router))

}

//variables
const numOfValues = 3

//Methods
func CreateTicket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("createTicketCalled")
}
func GetTickets(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, " ", "Hello World")
	fmt.Println("getTicketsCalled")
}
func GetTicket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getTicketCalled")
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
	ID    int    `json:"id,omitempty"`
	Lines []Line `json:"Lines,omitempty"`
}

type Line struct {
	ID     int   `json:"id,omitempty"`
	Values []int `json:"values,omitempty"`
}

// var tickets []Ticket

// tickets = append(tickets, Ticket{ID:"1"})
