# Flight Booking Application

Welcome to the Flight Booking Application! This backend service is designed to manage flight ticket bookings efficiently and is built using Go.

---

## Features

### Database Initialization
- Connects to a MySQL database using provided credentials.
- Initializes the database connection for seamless operation.

### Ticket Management
- **Get All Tickets**: Retrieve a list of all flight tickets.
- **Create Ticket**: Add a new flight ticket to the database.
- **Get Ticket by ID**: Retrieve a specific flight ticket by its ID.
- **Delete Ticket**: Remove a specific flight ticket by its ID.

---

## API Endpoints

- `GET /tickets` - Fetches all flight tickets.
- `POST /tickets` - Creates a new flight ticket.
- `GET /tickets/{id}` - Fetches a flight ticket by its ID.
- `DELETE /tickets/{id}` - Deletes a flight ticket by its ID.

---

## Usage

### 1. Run the Application

Execute the following command to start the application:

```sh
go run cmd/yourapp/main.go
```

### 2. Environment Configuration

Ensure the following prerequisites are met:

- A running and accessible MySQL database.
- Update the database connection string in the `main.go` file if necessary. Specifically, replace `User` and `Password` in the following line with your database username and password:
  ```go
  db, err := yourapp.InitDB("User:Password@tcp(127.0.0.1:3306)/flight_booking")
  ```

---

## Dependencies

This application leverages the following Go libraries:

- **[github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql)**: MySQL driver for Go.
- **[github.com/gorilla/mux](https://github.com/gorilla/mux)**: HTTP router and dispatcher for Go.

Install dependencies using:

```sh
go get github.com/go-sql-driver/mysql github.com/gorilla/mux
```

---

## How to Contribute

We welcome contributions! Follow these steps to get started:

1. Fork the repository.
2. Create a new branch:
   ```sh
   git checkout -b feature-branch
   ```
3. Commit your changes:
   ```sh
   git commit -am 'Add new feature'
   ```
4. Push to the branch:
   ```sh
   git push origin feature-branch
   ```
5. Open a Pull Request for review.

---

Thank you for using the Flight Booking Application. Feel free to contribute or raise issues to help us improve!
