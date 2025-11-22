# Golang Integration Testing Example

This repository accompanies the article [Integration Testing in Go with Testcontainers](https://blog.mohammedadekunle.com.ng/integration-testing-in-golang-with-docker-and-testcontainers) and demonstrates how to build a simple CRUD application in Golang with robust testing practices. The project covers unit testing and integration testing against real infrastructure using [Docker](https://docker.com) and [Testcontainers](https://testcontainers.com/).

## Features

- Create, Read, Update, Delete (CRUD) operations for notes with an `id`, `title` and `description`.
- Layered architecture:
  - **Database Access Layer (DBAL)**: Handles all database interactions using [pgx](https://github.com/jackc/pgx).
  - **Service Layer**: Contains business logic.
  - **HTTP Layer**: Exposes REST API endpoints.
- Integration tests with real PostgreSQL using Testcontainers.
- Unit tests with mocked dependencies.

## Prerequisites

- `Go 1.20` or higher.
- Docker.
- `make` (optional, for running migrations).
- PostgreSQL (if running outside Docker).

## Getting Started

1. **Clone the repository**

```bash
git clone https://github.com/the-code-genin/golang_integration_testing.git
cd golang_integration_testing
```

2. **Initialize Go modules**

```bash
go mod download
```

3. **Set up PostgreSQL**

(Optional) Using Docker:

```bash
docker run -p 5432:5432 -e POSTGRES_PASSWORD=password -d postgres:16.11
```

Running migrations:

```bash
make migrate-up
```

This ensure we have a database with all migrations applied.

4. **Running the Server**

Start the application:

```bash
go run .
```

The server should start on port `8080` (or the port specified via `env` variables):

```bash
# [GIN-debug] Listening and serving HTTP on :8080
```

## Project Structure

- `migrations/` - SQL migration files.
- `repository/` - Database Access Layer (DBAL).
- `service/` - Business logic layer.
- `http/` - REST API layer.
- `tests/`- Test helpers.
- `main.go`- Application entry point.
- `Makefile` - Optional automation commands.

## API Endpoints

- `POST /notes` - Create a note with title and description.
- `GET /notes/:id` - Fetch a single note by ID.
- `GET /notes` - Fetch all notes.
- `PUT /notes/:id` - Update a note by ID.
- `DELETE /notes/:id` - Delete a note by ID.
