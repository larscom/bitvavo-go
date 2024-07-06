run:
	@clear && printf '\e[3J' && go run ./example/main.go
build:
	go build -o ./bin/main ./example/main.go
test:
	go test -v ./.../ --race
