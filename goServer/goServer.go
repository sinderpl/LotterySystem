package main

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//variables
const numOfValues = 3

// Move to DB if have enough time
var existingTickets []Ticket

//Main runner
func main() {

	//Randomise functions are seeded with current time
	//to guarantee randomness
	rand.Seed(time.Now().UnixNano())

	//Test ticket
	testTicket := generateTicket()
	testTicket.ID = "12345678"
	values := []int{0, 1, 1}
	testTicket.Lines = append(testTicket.Lines, Line{ID: strconv.Itoa(rand.Intn(999999)), Values: values})
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
				result := calculateLineResult(value)
				ticket.Lines[index].Result = &result
				fmt.Println("line result check : ", ticket.Lines[index])
			}
			sort.Sort(ByResult(ticket.Lines))
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
	areTrailingUnique := true
	total := line.Values[0]
	//We begin the checks at the first value
	// for comparison checks
	firstValue := line.Values[0]

	// This can be calculated in a single array pass O(n)
	// By checking different situations and then ,
	// starting checks on highest value outcome first
	for x := 1; x < len(line.Values); x++ {
		//If total is 2 = 10pts
		total += line.Values[x]
		//If any values are different than the initial one
		// We can discard the possibility of the 5 pts outcome
		if line.Values[x] != firstValue {
			areSame = false
		} else {
			// If either 2nd or 3rd value are the same as first
			//we discard the 1 pt outcome
			areTrailingUnique = false
		}
		//If all these checks fail the 0pt outcome is applied
	}

	if total == 2 {
		return 10
	} else if areSame {
		return 5
	} else if areTrailingUnique {
		return 1
	} else {
		return 0
	}
}

//Structs
type Ticket struct {
	ID        string `json:"id,omitempty"`
	Lines     []Line `json:"Lines,omitempty"`
	IsChecked bool   `json:"IsChecked,omitempty"`
}

type Line struct {
	ID     string `json:"id,omitempty"`
	Values []int  `json:"values,omitempty"`
	Result *int   `json:"Result,omitempty"`
}

//Sorter
type ByResult []Line

func (p ByResult) Len() int {
	return len(p)
}
func (p ByResult) Less(i, j int) bool {
	return *p[i].Result > *p[j].Result
}
func (p ByResult) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
