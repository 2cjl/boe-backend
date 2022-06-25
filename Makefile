run: build

build:
	set GOOS=linux GOARCH=amd64
	go build -o boe-backend main.go
