package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/elafont/LanaChallenge/products"
	"github.com/elafont/LanaChallenge/tickets"
)

func NotFoundHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		NewResponse(http.StatusNotFound, http.StatusText(http.StatusNotFound), nil).WriteTo(w)
	}
}

func Root(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		NewResponse(
			http.StatusOK,
			fmt.Sprintf("Ticket Server V%s", version),
			&Data{Type: "Root Version", Content: version}).WriteTo(w)
	}
}

func NewTicket(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := len(s.Tickets)
		ticket := tickets.NewTicket(tickets.Promotions)

		s.Tickets = append(s.Tickets, ticket)
		NewResponse(
			http.StatusOK,
			fmt.Sprintf("New Ticket: %d", id),
			&Data{Type: "Ticket Status", Content: ticket.Status()}).WriteTo(w)
		return
	}
}

func ListTickets(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tks = make([]string, 0, len(s.Tickets))

		for ix, t := range s.Tickets {
			tks = append(tks, fmt.Sprintf("[%d] %s", ix, t.Status()))
		}

		NewResponse(
			http.StatusOK,
			fmt.Sprintf("Ticket List"),
			&Data{Type: "Tickets", Content: &tks}).WriteTo(w)
		return
	}
}

func GetTicket(s *Server) http.HandlerFunc { // This functions reacts to GET, to read a ticket
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getInt64Param(r, "ticket_id")
		if err != nil {
			s.Respond(w, paramError("ticket_id"))
			return
		}

		if int(id) >= len(s.Tickets) {
			// error, id not available  http.StatusBadRequest
			NewResponse(
				http.StatusBadRequest,
				fmt.Sprintf("ID: %d not available of %d", id, len(s.Tickets)),
				nil).WriteTo(w)
			return
		}

		ticket := s.Tickets[id]

		NewResponse(
			http.StatusOK,
			fmt.Sprintf("Ticket ID: %d", id),
			&Data{Type: "Ticket String", Content: ticket.String()}).WriteTo(w)
		return

	}
}

func AddProduct(s *Server) http.HandlerFunc { // This functions reacts to POST, to add a product to a ticket
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getInt64Param(r, "ticket_id")
		if err != nil {
			s.Respond(w, paramError("ticket_id"))
			return
		}

		if int(id) >= len(s.Tickets) {
			// error, id not available  http.StatusBadRequest
			NewResponse(
				http.StatusBadRequest,
				fmt.Sprintf("ID: %d not available of %d", id, len(s.Tickets)),
				nil).WriteTo(w)
			return
		}

		ticket := s.Tickets[id]
		r.ParseForm()

		quantity := r.FormValue("quantity")
		qty, err := strconv.Atoi(quantity)
		if err != nil {
			NewResponse(
				http.StatusBadRequest,
				fmt.Sprint("Quantity is not a number", quantity),
				nil).WriteTo(w)
			return
		}

		code := r.FormValue("code")
		product, err := products.GetProduct(code)
		if err != nil {
			NewResponse(
				http.StatusBadRequest,
				fmt.Sprintf("Product: %s unknown", code),
				nil).WriteTo(w)
			return
		}

		ticket.Add(qty, product)

		NewResponse(
			http.StatusOK,
			fmt.Sprintf("Ticket ID: %d", id),
			&Data{Type: "Ticket String", Content: ticket.String()}).WriteTo(w)
		return
	}
}
