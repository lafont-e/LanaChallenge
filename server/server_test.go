package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/lafont-e/LanaChallenge/tickets"
)

var s *Server

// CreateRequest creates the request and response for the given path, method and body
func CreateRequest(method string, target string, body io.Reader) (response *httptest.ResponseRecorder, request *http.Request) {
	response = httptest.NewRecorder()
	request = httptest.NewRequest(method, target, body)

	return
}

// Equals performs a deep equal comparison against two
// values and fails if they are not the same.
func Equals(tb testing.TB, expected, actual interface{}) {
	tb.Helper()

	//log.Printf("Equals %[1]v :: %[1]T\n\tgot: %[2]v :: %[2]T\n", expected, actual)
	if !reflect.DeepEqual(expected, actual) {
		tb.Fatalf(
			"\n\texp: %#[1]v (%[1]T)\n\tgot: %#[2]v (%[2]T)\n",
			expected,
			actual,
		)
	}
}

func BindJSON(r io.Reader, target interface{}) error {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

func TestMain(m *testing.M) {
	// Create Server
	s = &Server{
		Router:  mux.NewRouter(),
		Tickets: make([]*tickets.Ticket, 0, 8),
		Logger:  log.New(os.Stdout, "", log.LstdFlags),
	}

	s.RegisterRoutes()
	os.Exit(m.Run())
}

/*
	Test Guess
*/

func testNotFound(t *testing.T) {
	w, r := CreateRequest(
		http.MethodGet,
		"/nonexistant",
		nil,
	)

	s.Router.ServeHTTP(w, r)
	Equals(t, http.StatusNotFound, w.Code)
}

func testNewTicket(t *testing.T) {
	w, r := CreateRequest(
		http.MethodGet,
		"/newticket",
		nil,
	)

	s.Router.ServeHTTP(w, r)
	Equals(t, http.StatusOK, w.Code)

	var body Response

	if err := BindJSON(w.Body, &body); err != nil {
		t.Fatalf("error reading response %v", err)
	}

	Equals(t, fmt.Sprintf("New Ticket: %d", len(s.Tickets)-1), body.Message)
}

func testGetTicket(t *testing.T) {
	w, r := CreateRequest(
		http.MethodGet,
		"/ticket/0",
		nil,
	)

	s.Router.ServeHTTP(w, r)
	Equals(t, http.StatusOK, w.Code)

	var body Response

	if err := BindJSON(w.Body, &body); err != nil {
		t.Fatalf("error reading response %v", err)
	}

	Equals(t, "Ticket ID: 0", body.Message)
	Equals(t, "Ticket Status", body.Data.Type)
}

func testListTickets(t *testing.T) {
	testNewTicket(t)
	w, r := CreateRequest(
		http.MethodGet,
		"/tickets",
		nil,
	)

	s.Router.ServeHTTP(w, r)
	Equals(t, http.StatusOK, w.Code)

	var body Response

	if err := BindJSON(w.Body, &body); err != nil {
		t.Fatalf("error reading response %v", err)
	}

	Equals(t, "Ticket List", body.Message)
	Equals(t, "Tickets", body.Data.Type)
	if len(body.Data.Content.([]interface{})) != 2 {
		t.Fatal("error getting list, expecting 2 games")
	}
}

func testAddProduct(t *testing.T) {
	form := url.Values{}
	form.Set("code", "PEN")
	form.Set("quantity", "3")

	w, r := CreateRequest(
		http.MethodPost,
		"/ticket/0",
		strings.NewReader(form.Encode()),
	)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	s.Router.ServeHTTP(w, r)
	Equals(t, http.StatusOK, w.Code)

	var body Response

	if err := BindJSON(w.Body, &body); err != nil {
		t.Fatalf("error reading response %v", err)
	}

	Equals(t, "Ticket ID: 0", body.Message)
	Equals(t, "Ticket String", body.Data.Type)
	ticketStr := strings.Split(body.Data.Content.(string), "\n")
	lenTk := len(ticketStr)
	if lenTk < 3 {
		t.Fatal("error Ticket String is not right.")
	}

	tkstr := strings.TrimSpace(ticketStr[lenTk-3])
	Equals(t, tkstr, "Total:   10.00")
}

func TestServer(t *testing.T) {
	t.Run("Not Found", testNotFound)
	t.Run("New Ticket", testNewTicket)
	t.Run("Get Ticket", testGetTicket)
	t.Run("List Tickets", testListTickets)
	t.Run("Add Product", testAddProduct)
}
