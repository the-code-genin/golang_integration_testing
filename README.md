# Golang Integration Testing Example

This repository accompanies the article series [A Comprehensive Guide to Software Testing in Golang](https://the-code-genin.medium.com/list/a-comprehensive-guide-to-software-testing-in-golang-b1cd48050813) and demonstrates how to build a simple CRUD application in Golang with robust testing practices.

The project covers [unit](https://the-code-genin.medium.com/unit-testing-in-golang-writing-fast-reliable-tests-ee5f14ce5b06), [integration](https://the-code-genin.medium.com/integration-testing-in-golang-with-docker-and-testcontainers-a-practical-guide-ad654508284a) and [system](https://the-code-genin.medium.com/system-testing-in-golang-end-to-end-testing-in-practice-65e16a6eee85) testing against real infrastructure using [Docker](https://docker.com) and [Testcontainers](https://testcontainers.com/).

## Features

- Create, Read, Update, Delete (CRUD) operations for notes with an `id`, `title`, `description` and timestamps.
- Layered architecture:
  - **Database Access Layer (DBAL)**: Handles all database interactions using [pgx](https://github.com/jackc/pgx).
  - **Service Layer**: Contains business logic.
  - **HTTP Layer**: Exposes REST API endpoints.
- Integration tests against ephemeral PostgreSQL containers using Testcontainers.
- Unit tests with mocked dependencies.
- System tests with `httptest`.

## Prerequisites

- `Go 1.20` or higher.
- Docker.
- `make` (optional).

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

This ensures we have a database with all migrations applied.

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

- `POST /v1/notes` - Create a note with title and description.
- `GET /v1/notes/:id` - Fetch a single note by ID.
- `GET /v1/notes` - Fetch all notes.
- `PUT /v1/notes/:id` - Update a note by ID.
- `DELETE /v1/notes/:id` - Delete a note by ID.
