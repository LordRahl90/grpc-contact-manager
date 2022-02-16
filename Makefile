generate:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative contact/contact.proto

test:
	go test ./... -v --cover

test-service:
	go test ./services/$s -v --cover

ts: test-service