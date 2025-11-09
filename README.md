# Yard Planning System

A straightforward yard planning system designed with a clean and modular architecture using Go. It allows for efficient management of yard and block managements.

---

## Project Structure

This project follows a clear and modular architecture, separating concerns into distinct directories for better organization and maintainability.

```
.
├── app/
│   ├── controller/
│   ├── model/
│   ├── repository/
│   └── service/
├── database/
├── helper/
├── response/
├── web/
├── .env.example
├── go.mod
├── go.sum
└── main.go
```

---

## Getting Started

Follow these steps to get the project up and running on your local machine.

### Prerequisites

- [Go](https://golang.org/dl/) 1.23 or newer  
- [Postgres](https://www.postgresql.org/) database

---

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/kukuhpranata/yard-planning.git
cd yard-planning
```

### 2. Set Up Environment Variables

Copy the example `.env` file:

```bash
cp .env.example .env
```

Then open `.env` and configure your database connection and other required variables.

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Run Postgres Docker Container 
(Optional), you can use your own database.

```bash
cd /postgres-docker
docker compose up -d
```
Then run/import ```dbdump.sql```

---

## Running the Application

Start the application with:

```bash
go run main.go
```

---
