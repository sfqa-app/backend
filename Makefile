build:
	go build -o ./tmp/sfqa .

run:
	go run .

dev:
	air

test:
	go test ./...
