# Appointment Service API

This is a microservice for managing patient appointments on the medical platform. It provides a REST API for creating and retrieving appointments, and updating their statuses. It uses synchronous communication with the Doctor Service to validate data.

## Technologies

* Language: Go 1.25.5
* Framework: Gin for HTTP routing
* Database: MongoDB
* Architecture: Clean Architecture (Model, Repository, Usecase, Transport, Client)

## How to Run the Project

### Requirements
1. Go installed (version 1.25.5 or higher).
2. A running MongoDB server on the default port (localhost:27017). The service will automatically create the appointment_db database and the appointments collection.
3. The Doctor Service must be running simultaneously on http://localhost:8081.

### Running the service
In the root directory of the project, run the following command in your terminal:
```bash
  go run ./cmd/appointment-s/main.go
```
The service will start on port 8082. It supports Graceful Shutdown for safe termination.

---

## REST API Endpoints

### 1. Create an Appointment
Creates a new appointment. The service will verify if the provided doctor_id exists by calling the Doctor Service over HTTP.

* URL: /appointments
* Method: POST
* Request Body (JSON Example):
```json
  {
      "title": "Annual Checkup",
      "description": "Patient feels pain in the left arm",
      "doctor_id": "60d5ec49f1b2c8a14b53ef12"
  }
```
**Business Rules (Validation):**
* The title field is required.
* The doctor_id field is required.
* The doctor must exist in the Doctor Service.
* Newly created appointments automatically receive the status "new".

**Responses:**
* 200 OK: Appointment created successfully.
* 400 Bad Request: Missing fields or the referenced doctor does not exist.
* 500 Internal Server Error: Failed to communicate with Doctor Service or database error.

---

### 2. Get Appointment by ID
Retrieves information about a specific appointment.

* URL: /appointments/:id
* Method: GET

**Responses:**
* 200 OK: Successful request.
* 404 Not Found: If the appointment is not found ("there is no appointment like this").
* 500 Internal Server Error: Server error.

---

### 3. Get All Appointments
Returns an array of all appointments.

* URL: /appointments
* Method: GET

**Responses:**
* 200 OK: Successful request. Returns an empty array "[]" if there are no records.
* 500 Internal Server Error: Database error.

---

### 4. Update Appointment Status
Updates the status of an existing appointment.

* URL: /appointments/:id/status
* Method: PATCH
* Request Body (JSON Example):
```json
  {
      "status": "in_progress"
  }
```
**Business Rules:**
* Valid statuses are: new, in_progress, done.
* You cannot transition a status from done back to new.

**Responses:**
* 200 OK: Status updated successfully.
* 400 Bad Request: Invalid status value or forbidden transition.
* 404 Not Found: Appointment not found.
* 500 Internal Server Error: Server error.

---

## Project Structure (Clean Architecture)

* `cmd/appointment-s/main.go`: Entry point, setups dependencies, HTTP router, and initializes connections.
* `internal/client/doctor_client.go`: HTTP client for synchronous communication with the Doctor Service.
* `internal/model: Data structures` (Appointment) and domain errors.
* `internal/repository`: Interacts with the MongoDB database (appointment_db).
* `internal/usecase`: Business logic layer managing validation and status transition rules.
* `internal/transport/http`: Handlers mapping HTTP requests to usecase methods using the Gin framework.