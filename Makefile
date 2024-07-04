run:
	@clear && printf '\e[3J' && go run ./cmd/main.go
build:
	go build -o ./bin/main ./cmd/main.go
test:
	go test -v ./.../ --race
