.PHONY: run go-mod

run:
	go run ./cmd/server

go-mod:
	go mod tidy && go mod vendor