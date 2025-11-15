.PHONY: migrate-up
migrate-up:
	@migrate -path ./migrations -database "postgres://postgres:password@localhost:5432/postgres?sslmode=disable" up

.PHONY: migrate-down
migrate-down:
	@migrate -path ./migrations -database "postgres://postgres:password@localhost:5432/postgres?sslmode=disable" down
