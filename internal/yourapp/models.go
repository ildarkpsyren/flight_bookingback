package yourapp

import "time"

type Ticket struct {
	ID               int       `json:"id"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	BookingID        string    `json:"booking_id"`
	IIN              int64     `json:"iin"`
	CreatedAt        time.Time `json:"created_at"`
	DepartureTime    string    `json:"departure_time"`
	ArrivalTime      string    `json:"arrival_time"`
	DepartureAirport string    `json:"departure_airport"`
	ArrivalAirport   string    `json:"arrival_airport"`
}
