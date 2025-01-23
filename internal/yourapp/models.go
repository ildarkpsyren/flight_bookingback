package yourapp

type Ticket struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Age       int    `json:"age"`
	Departure string `json:"departure"`
	Arrival   string `json:"arrival"`
	Date      string `json:"date"`
	Time      string `json:"time"`
}
