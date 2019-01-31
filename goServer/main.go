package main

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//variables
const numOfValues = 3

var existingTickets []Ticket

//Main runner
func main() {
	//Test ticket
	testTicket := generateTicket()
	testTicket.ID = "12345678"
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
	fmt.Println("Create tickets called")

	ticket := generateTicket()
	existingTickets = append(existingTickets, ticket)
	json.NewEncoder(w).Encode(ticket.ID)
}

//Returns all ticket ID's
func GetTickets(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get tickets called")

	var ticketList []string
	if existingTickets != nil {
		for _, value := range existingTickets {
			ticketList = append(ticketList, value.ID)
		}
	}

	json.NewEncoder(w).Encode(ticketList)
}

//Return the specified ticket
func GetTicket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get ticket called")
	parameters := mux.Vars(r)
	ticketID := parameters["id"]
	ticket := findTicket(ticketID)
	if ticket == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ticket does not exist"))
	} else {
		json.NewEncoder(w).Encode(ticket)
	}
}

//Amend ticket lines
func UpdateTicket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update ticket called")
	parameters := mux.Vars(r)
	ticketID := parameters["id"]
	ticket := findTicket(ticketID)
	if ticket == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ticket does not exist"))
	} else {
		if ticket.IsChecked == true {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Ticket status has been checked, no changes can be made"))
		} else {
			ticket.Lines = append(ticket.Lines, generateTicketLine())
			json.NewEncoder(w).Encode(ticket.ID)
		}
	}
}

//Retrieves the status of a ticket
func GetTicketStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get Ticket status called")
	parameters := mux.Vars(r)
	ticketID := parameters["id"]
	ticket := findTicket(ticketID)
	if ticket == nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ticket does not exist"))
	} else {
		if ticket.IsChecked {
			w.Write([]byte("Ticket already checked"))
		} else {
			ticket.IsChecked = true

			//Calculat results
			for index, value := range ticket.Lines {
				fmt.Println("line result check : ", value)
				ticket.Lines[index].Result = calculateLineResult(value)
			}

			// sort.Sort(sort.Reverse(ticket.Lines))
		}
		json.NewEncoder(w).Encode(ticket)
	}
}

//Private methods
func generateTicket() Ticket {
	lines := make([]Line, 1)
	lines[0] = generateTicketLine()
	newTicket := Ticket{ID: strconv.Itoa(rand.Intn(999999)), Lines: lines, IsChecked: false}
	return newTicket
}

func generateTicketLine() Line {
	values := make([]int, numOfValues)
	for index, _ := range values {
		values[index] = rand.Intn(3) // 0 - 2 inclusive
	}
	newLine := Line{ID: strconv.Itoa(rand.Intn(999999)), Values: values}
	return newLine
}

func findTicket(ticketID string) *Ticket {
	var ticket *Ticket
	for index, _ := range existingTickets {
		if existingTickets[index].ID == ticketID {
			ticket = &existingTickets[index]
		}
	}
	return ticket
}

func calculateLineResult(line Line) int {
	fmt.Println("Generating results for line")

	areSame := true
	areUnique := true
	total := line.Values[0]
	firstValue := line.Values[0]
	for x := 1; x < len(line.Values); x++ {
		total += line.Values[x] //IF total is 2 = 10pts
		if line.Values[x] != firstValue {
			areSame = false // If values are equal = 5pts
		} else {
			areUnique = false // If not unique to first number = 0 pts
		}
	}

	if total == 2 {
		return 10
	} else if areSame {
		return 5
	} else if areUnique {
		return 1
	} else {
		return 0
	}
}

//Sorter
// type Lines []Line

// func (p Lines) Len() int           { return len(p) }
// func (p Lines) Less(i, j int) bool { return p[i].Result < p[j].Result }
// func (p Lines) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

//Structs
type Ticket struct {
	ID        string `json:"id,omitempty"`
	Lines     []Line `json:"Lines,omitempty"`
	IsChecked bool   `json:"IsChecked,omitempty"`
}

type Line struct {
	ID     string `json:"id,omitempty"`
	Values []int  `json:"values,omitempty"`
	Result int    `json:"Result,omitempty"`
}
