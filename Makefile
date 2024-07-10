build:
	@go build -o bin/lb main.go

run: build
	@./bin/lb
