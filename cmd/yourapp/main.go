package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"yourapp/internal/yourapp"
)

func main() {
	db, err := yourapp.InitDB("User:Password@tcp(127.0.0.1:3306)/flight_booking")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}(db)

	router := mux.NewRouter()
	router.HandleFunc("/tickets", yourapp.GetTickets(db)).Methods("GET")
	router.HandleFunc("/tickets", yourapp.CreateTicket(db)).Methods("POST")
	router.HandleFunc("/tickets/{id}", yourapp.GetTicketByID(db)).Methods("GET")
	router.HandleFunc("/tickets/{id}", yourapp.DeleteTicket(db)).Methods("DELETE")
	router.HandleFunc("/tickets/{first_name}/{last_name}/{booking_id}", yourapp.GetTicketByNameAndBookingID(db)).Methods("GET")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
