package yourapp

import (
	"time"
)

type IIN string

type Ticket struct {
	ID               int       `json:"id"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	BookingID        string    `json:"booking_id"`
	IIN              string    `json:"iin"`
	CreatedAt        time.Time `json:"created_at"`
	DepartureTime    time.Time `json:"departure_time"`
	ArrivalTime      time.Time `json:"arrival_time"`
	DepartureAirport string    `json:"departure_airport"`
	ArrivalAirport   string    `json:"arrival_airport"`
}
