run:
	cd backend && go run main.go
worker:
	cd mailer && go run server.go
build:
	go build ./...
test:
	cd backend && go test ./tests/...