# Remitly-Internship-Exercise
This is a repository for Remitly Summer Internship 2025 task

# SWIFT Code REST API

A RESTful API built with Go and Gin that allows users to manage SWIFT codes. It supports retrieving, adding, and deleting SWIFT codes from a PostgreSQL database.

## Features

- Retrieve SWIFT codes
- Add new SWIFT codes
- Delete existing SWIFT codes
- Structured JSON responses
- Containerized setup for easy deployment

## Tech Stack

- **Language:** Go
- **Framework:** Gin
- **Database:** PostgreSQL 17
- **Containerization:** Docker (assumed, based on standard practices)

## Getting Started

### Prerequisites

- Go installed
- PostgreSQL 17 installed and running
- Docker (optional, for containerized setup)

### Database setup
- Create database swiftdb
- create table swift-codes
  
<pre>CREATE TABLE swift_codes (
    id SERIAL PRIMARY KEY,
    swift_code VARCHAR NOT NULL,
    bank_name TEXT NOT NULL,
    is_headquarter BOOLEAN NOT NULL,
    branch_of VARCHAR,
    country_iso2 VARCHAR NOT NULL,
    country_name TEXT NOT NULL,
    city TEXT,
    address TEXT,
    time_zone TEXT
);</pre>

### Installation

1. Clone the repository or unzip the project folder:
   ```bash
   git clone https://github.com/yourusername/swift-api.git
   cd swift-api
2. Create your swift_codes table in your database
3. Configure your database connection in config/db.go
4. Run the application:
   - open in swift-api repository
   ```bash
   go run main.go
### API ENDPOINTS
Endpoint 0: Import SWIFT Codes from CSV

POST /v1/import-swift?file={filename}.csv

Imports a batch of SWIFT codes from a CSV file located in the server's import/ directory.
Accepts a .csv filename via the file query parameter.
The file must be placed in the local import/ directory.
Parses and inserts all SWIFT code entries from the file into the database.
Successful Response:

{
  "message": "SWIFT codes imported successfully!"
}</pre>

Error Responses:

-Missing query parameter:

{
  "error": "Missing 'file' query parameter"
}</pre>

-Invalid file type:

{
  "error": "Only .csv files are allowed"
}</pre>

-CSV parsing or database insertion failure:

{
  "error": "Failed to parse CSV",
  "details": "detailed error message here"
}</pre>

{
  "error": "Failed to insert SWIFT codes",
  "details": "detailed error message here"
}</pre>

Endpoint 1: Retrieve SWIFT Code Details

GET /v1/swift-codes/{swift-code}

Returns information about a specific SWIFT code, whether it's a headquarter or a branch.

Response for a Headquarter SWIFT Code:

<pre>{
  "address": "string",
  "bankName": "string",
  "countryISO2": "string",
  "countryName": "string",
  "isHeadquarter": true,
  "swiftCode": "string",
  "branches": [
    {
      "address": "string",
      "bankName": "string",
      "countryISO2": "string",
      "isHeadquarter": false,
      "swiftCode": "string"
    }
    // ... more branches
  ]
}</pre>

Response for a Branch SWIFT Code:

<pre>{
  "address": "string",
  "bankName": "string",
  "countryISO2": "string",
  "countryName": "string",
  "isHeadquarter": false,
  "swiftCode": "string"
}</pre>

Endpoint 2: Get All SWIFT Codes for a Country

GET /v1/swift-codes/country/{countryISO2code}

Returns all SWIFT codes (both headquarters and branches) for the given 2-letter country code.

Response:

<pre>{
  "countryISO2": "string",
  "countryName": "string",
  "swiftCodes": [
    {
      "address": "string",
      "bankName": "string",
      "countryISO2": "string",
      "isHeadquarter": true,
      "swiftCode": "string"
    },
    {
      "address": "string",
      "bankName": "string",
      "countryISO2": "string",
      "isHeadquarter": false,
      "swiftCode": "string"
    }
    // ... more SWIFT codes
  ]
}</pre>

Endpoint 3: Add a New SWIFT Code

POST /v1/swift-codes

Adds a new SWIFT code entry to the database.

Request:

<pre>{
  "address": "string",
  "bankName": "string",
  "countryISO2": "string",
  "countryName": "string",
  "isHeadquarter": true,
  "swiftCode": "string"
}</pre>

Response:

{
  "message": "string"
}

Endpoint 4: Delete a SWIFT Code

DELETE /v1/swift-codes/{swift-code}

Deletes the SWIFT code from the database if it exists.

Response:

{
  "message": "string"
}

### USE EXAMPLES USING curl

Start with importing 

ENDPOINT 0 IMPORT (change example.csv to the name of the file you inserted in import folder)

curl -X POST "http://localhost:8080/v1/import-swift?file=example.csv"

ENDPOINT 1 (change CODE for your swift-code)

curl -X GET http://localhost:8080/v1/swift-codes/CODE

headquarter without branches

http://localhost:8080/v1/swift-codes/SUSRPLP1XXX

headquarter with one branch

http://localhost:8080/v1/swift-codes/UBPGMCMXXXX

headquarter with many branches

http://localhost:8080/v1/swift-codes/TPEOPLPWXXX

ENDPOINT 2

Working country code

curl -X GET http://localhost:8080/v1/swift-codes/country/PL

Not existing country code

curl -X GET http://localhost:8080/v1/swift-codes/country/PLL

ENDPOINT 3 (insert using curl)

curl -X POST http://localhost:8080/v1/swift-codes -H "Content-Type: application/json" -d "{\"swiftCode\":\"TESTPLP1XXX\",\"bankName\":\"Test Bank\",\"address\":\"Test Address 1\",\"countryISO2\":\"PL\",\"countryName\":\"Poland\",\"isHeadquarter\":true}"

ENDPOINT 4
curl -X DELETE http://localhost:8080/v1/swift-codes/TESTPLP1XXX
