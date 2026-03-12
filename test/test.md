## Check the test coverage
go test -v -coverpkg=./... -coverprofile=coverage.out ./...

## Check test coverage without main
go test -v \
-coverpkg=./internal/...,./router/...,./database/...,./config/... \
-coverprofile=coverage.out ./...

## Check result in browser
go tool cover -html=coverage.out