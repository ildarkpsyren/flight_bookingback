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

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetTickets(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, first_name, last_name, booking_id, iin, created_at, departure_time, arrival_time, departure_airport, arrival_airport FROM tickets")
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
			var createdAt string
			if err := rows.Scan(&ticket.ID, &ticket.FirstName, &ticket.LastName, &ticket.BookingID, &ticket.IIN, &createdAt, &ticket.DepartureTime, &ticket.ArrivalTime, &ticket.DepartureAirport, &ticket.ArrivalAirport); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ticket.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
			if err != nil {
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
		var createdAt string
		query := `SELECT id, first_name, last_name, booking_id, iin, created_at, departure_time FROM tickets WHERE id = ?`
		err := db.QueryRow(query, id).Scan(&ticket.ID, &ticket.FirstName, &ticket.LastName, &ticket.BookingID, &ticket.IIN, &createdAt, &ticket.DepartureTime)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Ticket not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		ticket.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ticket); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func CheckTicket(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			LastName  string `json:"last_name"`
			BookingID string `json:"booking_id"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var ticket Ticket
		query := `SELECT id, first_name, last_name, booking_id, iin, created_at, departure_time FROM tickets WHERE last_name = ? AND booking_id = ?`
		err := db.QueryRow(query, request.LastName, request.BookingID).Scan(&ticket.ID, &ticket.FirstName, &ticket.LastName, &ticket.BookingID, &ticket.IIN, &ticket.CreatedAt, &ticket.DepartureTime)
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

		// Validate IIN length
		if len(ticket.IIN) != 12 {
			http.Error(w, "IIN must be 12 digits long", http.StatusBadRequest)
			return
		}

		// Parse the input times in the format "2006-01-02 15:04"
		departureTime, err := time.Parse("2006-01-02 15:04", ticket.DepartureTime)
		if err != nil {
			http.Error(w, "Invalid departure time format", http.StatusBadRequest)
			return
		}
		arrivalTime, err := time.Parse("2006-01-02 15:04", ticket.ArrivalTime)
		if err != nil {
			http.Error(w, "Invalid arrival time format", http.StatusBadRequest)
			return
		}

		ticket.DepartureTime = departureTime.Format("2006-01-02 15:04")
		ticket.ArrivalTime = arrivalTime.Format("2006-01-02 15:04")
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

func DeleteAllTickets(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := `DELETE FROM tickets`
		_, err := db.Exec(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func SearchTickets(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")
		iin := r.URL.Query().Get("iin")

		var tickets []Ticket
		query := `SELECT id, first_name, last_name, booking_id, iin, created_at, departure_time, arrival_time, departure_airport, arrival_airport FROM tickets WHERE first_name = ? AND last_name = ? AND iin = ?`
		rows, err := db.Query(query, firstName, lastName, iin)
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

		for rows.Next() {
			var ticket Ticket
			var createdAt string
			if err := rows.Scan(&ticket.ID, &ticket.FirstName, &ticket.LastName, &ticket.BookingID, &ticket.IIN, &createdAt, &ticket.DepartureTime, &ticket.ArrivalTime, &ticket.DepartureAirport, &ticket.ArrivalAirport); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ticket.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
			if err != nil {
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
