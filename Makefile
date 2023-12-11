start:
	docker-compose up --abort-on-container-exit --force-recreate --build server --build client

test:
	go clean --testcache
	go test ./...

test-short:
	go clean --testcache
	go test -short ./...

start-server:
	go run cmd/server/main.go

start-client:
	go run cmd/client/main.go
