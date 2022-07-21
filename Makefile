run:
	cd backend && go run main.go
build:
	cd backend && go build && cd .. && cd mailer && go build
test:
	cd backend && go test ./tests/...