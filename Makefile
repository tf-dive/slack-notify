
.PHONY: build
build:
	mkdir -p bin/
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/slack-notify ./main.go
