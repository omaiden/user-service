TEST_FLAGS:=-race -timeout 10s

export TEST_DB_URL=postgres://user:password@localhost:5432/%s?sslmode=disable

run:
	go run main.go
.PHONY: run

test:
	@go test -v ./...
.PHONY: test

test.unit: export SKIP_INTEGRATION=1
test.unit:
	go test ./... ${TEST_FLAGS}
.PHONY: test.unit

vendor:
	go mod tidy
	go mod vendor
.PHONY: vendor
