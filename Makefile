all: osx

osx:
	env GOOS=darwin GOARCH=amd64 go build -o chicka main.go

linux:
	env GOOS=linux GOARCH=amd64 go build -o chicka main.go
