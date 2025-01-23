package yourapp

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func GetTickets(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, first_name, last_name, booking_id, iin, created_at, departure_time FROM tickets")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}(rows)

		var tickets []Ticket
		for rows.Next() {
			var ticket Ticket
			if err := rows.Scan(&ticket.ID, &ticket.FirstName, &ticket.LastName, &ticket.BookingID, &ticket.IIN, &ticket.CreatedAt, &ticket.DepartureTime); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tickets = append(tickets, ticket)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tickets); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetTicketByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var ticket Ticket
		query := `SELECT id, first_name, last_name, booking_id, iin, created_at, departure_time FROM tickets WHERE id = ?`
		err := db.QueryRow(query, id).Scan(&ticket.ID, &ticket.FirstName, &ticket.LastName, &ticket.BookingID, &ticket.IIN, &ticket.CreatedAt, &ticket.DepartureTime)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Ticket not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ticket); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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

		// Parse the input times without seconds
		departureTime, err := time.Parse("2006-01-02T15:04", ticket.DepartureTime.Format("2006-01-02T15:04"))
		if err != nil {
			http.Error(w, "Invalid departure time format", http.StatusBadRequest)
			return
		}
		arrivalTime, err := time.Parse("2006-01-02T15:04", ticket.ArrivalTime.Format("2006-01-02T15:04"))
		if err != nil {
			http.Error(w, "Invalid arrival time format", http.StatusBadRequest)
			return
		}

		// Format the times to include seconds
		ticket.DepartureTime = departureTime
		ticket.ArrivalTime = arrivalTime

		ticket.CreatedAt = time.Now()
		ticket.BookingID = generateBookingID() // Generate a unique booking ID

		query := `INSERT INTO tickets (first_name, last_name, booking_id, iin, departure_time, arrival_time, departure_airport, arrival_airport) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
		result, err := db.Exec(query, ticket.FirstName, ticket.LastName, ticket.BookingID, ticket.IIN, ticket.DepartureTime, ticket.ArrivalTime, ticket.DepartureAirport, ticket.ArrivalAirport)
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
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(ticket); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// generateBookingID generates a unique booking ID
func generateBookingID() string {
	return fmt.Sprintf("BKG-%d", time.Now().UnixNano())
}

func DeleteTicket(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		query := `DELETE FROM tickets WHERE id = ?`
		_, err := db.Exec(query, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func GetTicketByNameAndBookingID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		firstName := vars["first_name"]
		lastName := vars["last_name"]
		bookingID := vars["booking_id"]

		var ticket Ticket
		query := `SELECT id, first_name, last_name, booking_id, iin, created_at, departure_time FROM tickets WHERE first_name = ? AND last_name = ? AND booking_id = ?`
		err := db.QueryRow(query, firstName, lastName, bookingID).Scan(&ticket.ID, &ticket.FirstName, &ticket.LastName, &ticket.BookingID, &ticket.IIN, &ticket.CreatedAt, &ticket.DepartureTime)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Ticket not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ticket); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
