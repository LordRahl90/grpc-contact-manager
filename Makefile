generate:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative contact/contact.proto

test:
	go test ./... -v --cover

test-service:
	go test ./services/$s -v --cover

build:
	go build -o ./cmd/grpc-contact ./cmd

docker-build:
	docker build -f .docker/Dockerfile -t lordrahl/grpc-contact:latest .

ts: test-service
db: docker-build