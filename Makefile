.PHONY:run
run: test
	go run main.go

.PHONY:build
build: test
	go build

.PHONY:test
test: vet fmt
	go test ./...

.PHONY:vet
vet:
	go vet ./...

.PHONY:fmt
fmt:
	go fmt ./...
