# EMI Tracker

EMI Tracker is a web-based application for managing personal EMI (Equated Monthly Installment) records. It provides a RESTful API for user registration, authentication, and EMI record management, built with Go and MySQL.

## Features

- User registration and login with password hashing
- Add, update, delete, and view EMI records
- Track total loaned, paid, active, and completed EMIs per user
- Secure password storage using bcrypt
- JSON-based API responses

## Tech Stack

- Go (Golang)
- MySQL
- Standard library net/http for routing
- bcrypt for password hashing
- godotenv for environmental variables

## Project Structure

```
├── main.go                # Application entry point
├── database/              # Database connection logic
├── handlers/              # HTTP handlers for users and EMI records
├── models/                # Data models
├── repo/                  # Database access/repository layer
├── router/                # HTTP router setup
├── utils/                 # Utility functions (JSON, password)
├── go.mod, go.sum         # Go module files
```

## Setup & Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/sajidzamanme/emi-tracker.git
   cd emi-tracker
   ```

2. **Configure the database:**
   - Create a MySQL database named `emiTracker`.
   - Use a `.env` file for variables: "ROOT", "PASSWORD", "PORT".
   - Run the SQL in `emiTrackerConfig.session.sql` to create tables.

3. **Install dependencies:**
   ```bash
   go mod tidy
   ```

4. **Run the application:**
   ```bash
   go run main.go
   ```
   The server will start on `localhost:8080`.

## API Endpoints

### User Endpoints

- `POST   /users/signup`         — Register a new user
- `POST   /users/login`          — User login
- `GET    /users`                — List all users
- `GET    /users/{userID}`       — Get user by ID
- `PUT    /users/{userID}`       — Update user
- `DELETE /users/{userID}`       — Delete user
- `GET    /users/{userID}/emirecords` — List all EMI records for a user

### EMI Record Endpoints

- `GET    /emirecords/{recordID}`         — Get EMI record by ID
- `POST   /emirecords/{userID}`           — Add EMI record for a user
- `PUT    /emirecords/{recordID}`         — Update EMI record
- `DELETE /emirecords/{recordID}`         — Delete EMI record
- `GET    /emirecords/{recordID}/payinstallment` — Pay an installment

## Contributing

1. Fork the repository
2. Create a new branch (`git checkout -b feature-branch`)
3. Commit your changes
4. Push to your fork and open a pull request
