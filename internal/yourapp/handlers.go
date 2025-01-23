package yourapp

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetTickets(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, surname, age, departure, arrival, date, time FROM tickets")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {

			}
		}(rows)

		var tickets []Ticket
		for rows.Next() {
			var ticket Ticket
			if err := rows.Scan(&ticket.ID, &ticket.Name, &ticket.Surname, &ticket.Age, &ticket.Departure, &ticket.Arrival, &ticket.Date, &ticket.Time); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tickets = append(tickets, ticket)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(tickets)
		if err != nil {
			return
		}
	}
}

func CreateTicket(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ticket Ticket
		if err := json.NewDecoder(r.Body).Decode(&ticket); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := db.Exec("INSERT INTO tickets (name, surname, age, departure, arrival, date, time) VALUES (?, ?, ?, ?, ?, ?, ?)",
			ticket.Name, ticket.Surname, ticket.Age, ticket.Departure, ticket.Arrival, ticket.Date, ticket.Time)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ticket.ID = int(id)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(ticket)
		if err != nil {
			return
		}
	}
}

func GetTicketByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
			return
		}

		var ticket Ticket
		err = db.QueryRow("SELECT id, name, surname, age, departure, arrival, date, time FROM tickets WHERE id = ?", id).Scan(
			&ticket.ID, &ticket.Name, &ticket.Surname, &ticket.Age, &ticket.Departure, &ticket.Arrival, &ticket.Date, &ticket.Time)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Ticket not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(ticket)
		if err != nil {
			return
		}
	}
}

func DeleteTicket(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("DELETE FROM tickets WHERE id = ?", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
