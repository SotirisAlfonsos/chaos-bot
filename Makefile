test:
	go test ./... -count=1

test-short:
	go test ./... -count=1 -v -short

build:
	go build -o bin/chaos-bot main.go

run:
	go run main.go

run-tls:
	go run main.go --config.file=config/example/example_tls_config.yml