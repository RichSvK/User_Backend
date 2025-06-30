Repository Source: https://github.com/golang-migrate/migrate

Run command after install golang migrate:
- go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
- migrate -source file://path/to/migrations -database postgres://username:password@localhost:5432/database?sslmode=disable up `N`