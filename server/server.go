package main

import (
	"context"
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

//TODO
//Remove comments
//Change error status
//Empty is valid return
//Stop checking results each call
//The tickets can never be empty

//variables
const (
	NUM_VALUES  = 3
	SERVER_PORT = ":8080"
	PUT         = "PUT"
	POST        = "POST"
	GET         = "GET"
	DELETE      = "DELETE"
)

// Move to DB if I have enough time
var existingTickets []Ticket
var server *http.Server

//Main runner
func main() {
	//Randomise functions are seeded with current time to guarantee randomness
	rand.Seed(time.Now().UnixNano())

	//Test ticket
	testTicket := generateTicket()
	testTicket.ID = "12345678"
	values := []int{0, 1, 1}
	testTicket.Lines = append(testTicket.Lines, Line{ID: strconv.Itoa(rand.Intn(999999)), Values: values})
	existingTickets = append(existingTickets, testTicket)
	server = startServer()
}

func startServer() *http.Server {
	router := mux.NewRouter()

	router.HandleFunc("/ticket", CreateTicket).Methods(POST)
	router.HandleFunc("/ticket", GetTickets).Methods(GET)
	router.HandleFunc("/ticket/{id}", GetTicket).Methods(GET)
	router.HandleFunc("/ticket/{id}/{lines}", UpdateTicket).Methods(PUT)
	router.HandleFunc("/status/{id}", GetTicketStatus).Methods(PUT)

	srv := &http.Server{
		Addr: SERVER_PORT,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 2,
		ReadTimeout:  time.Second * 2,
		IdleTimeout:  time.Second * 4,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
	return srv
}

func ShutDown() {
	fmt.Println("Shutdown called")
	var wait time.Duration = 10
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	if server != nil {
		server.Shutdown(ctx)
	}
}

//Public Methods

//Generates a new ticket with single line and returns ID to user
func CreateTicket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create ticket called")

	ticket := generateTicket()
	existingTickets = append(existingTickets, ticket)

	json.NewEncoder(w).Encode([...]Ticket{ticket})
}

//Returns all tickets
func GetTickets(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get tickets called")

	if existingTickets == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Ticket does not exist"))
		return
	} else {
		json.NewEncoder(w).Encode(existingTickets)
	}

}

//Return the specified ticket
func GetTicket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get ticket called")
	parameters := mux.Vars(r)
	ticketID := parameters["id"]
	ticket := findTicket(ticketID)
	if ticket == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Ticket does not exist"))
		return
	} else {
		json.NewEncoder(w).Encode(([...]Ticket{*ticket}))
	}
}

//Amend ticket lines
func UpdateTicket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update ticket called")
	parameters := mux.Vars(r)
	ticketID := parameters["id"]
	linesString := parameters["lines"]
	lineNum, err := strconv.Atoi(linesString)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("Invalid amount of lines"))
		return
	}
	ticket := findTicket(ticketID)
	if ticket == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Ticket does not exist"))
		return
	} else {
		if ticket.IsChecked == true {
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("Ticket status has been checked, no changes can be made"))
			return
		} else {
			for x := 1; x <= lineNum; x++ {
				ticket.Lines = append(ticket.Lines, generateTicketLine())
			}
			json.NewEncoder(w).Encode(([...]Ticket{*ticket}))
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
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Ticket does not exist"))
		return
	} else {
		if !ticket.IsChecked {
			//Calculate results
			for index, value := range ticket.Lines {
				result := calculateLineResult(value)
				ticket.Lines[index].Result = result
			}
			sort.Sort(ByResult(ticket.Lines)) //#TODO check if this modifies the original, otherwise it will be returned without being sorted
			ticket.IsChecked = true
		}
		json.NewEncoder(w).Encode(([...]Ticket{*ticket}))
	}
}

//Private methods

//Generates and returns a new ticket
func generateTicket() Ticket {
	lines := make([]Line, 1)
	lines[0] = generateTicketLine()
	newTicket := Ticket{ID: strconv.Itoa(rand.Intn(999999)), Lines: lines, IsChecked: false}
	return newTicket
}

//Generate and returns a new ticket line with randomised values
func generateTicketLine() Line {
	values := make([]int, NUM_VALUES)
	for index, _ := range values {
		values[index] = rand.Intn(3) // 0 - 2 inclusive
	}
	newLine := Line{ID: strconv.Itoa(rand.Intn(999999)), Values: values}
	return newLine
}

//Looks for the ticket and returns if found
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
	firstValue := line.Values[0]

	for x := 1; x < len(line.Values); x++ {
		total += line.Values[x]
		if line.Values[x] != firstValue {

			areSame = false
		} else {
			areTrailingUnique = false
		}
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

//Objects
type Ticket struct {
	ID        string `json:"id,omitempty"`
	Lines     []Line `json:"lines,omitempty"`
	IsChecked bool   `json:"isChecked,omitempty"`
}

type Line struct {
	ID     string `json:"id,omitempty"`
	Values []int  `json:"values,omitempty"`
	Result int    `json:"result,omitempty"`
}

//Sorter
type ByResult []Line

func (p ByResult) Len() int {
	return len(p)
}
func (p ByResult) Less(i, j int) bool {
	return p[i].Result > p[j].Result
}
func (p ByResult) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
