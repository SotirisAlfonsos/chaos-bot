test:
	go test ./... -count=1 -v

test-integration:
	go test ./... -count=1 -v -tags=integration

build:
	go build -o bin/chaos-bot main.go

run:
	go run main.go

run-tls:
	go run main.go --config.file=config/example/example_tls_config.yml