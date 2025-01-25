package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
	"yourapp/internal/yourapp"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)
		log.Printf("Method: %s, URL: %s, Status: %d, Duration: %s, Comment: %s",
			r.Method, r.URL.Path, rw.statusCode, time.Since(start), "Request processed")
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

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
	router.Use(yourapp.CorsMiddleware) // Use the CORS middleware
	router.Use(loggingMiddleware)      // Use the logging middleware

	router.HandleFunc("/api/tickets/search", yourapp.SearchTickets(db)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/tickets/delall", yourapp.DeleteAllTickets(db)).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/tickets/viewall", yourapp.GetTickets(db)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/tickets/createticket", yourapp.CreateTicket(db)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/tickets/viewby/{id}", yourapp.GetTicketByID(db)).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/tickets/delby/{id}", yourapp.DeleteTicket(db)).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/tickets/check", yourapp.CheckTicket(db)).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/tickets", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
