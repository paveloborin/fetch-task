APP?=fetchtasker

#RUN LOCAL
run:
	go run cmd/fetchtask.go

build:
	go build -o ./bin/${APP} ./cmd/fetchtask.go

#TESTS
test_unit:
	go test `go list ./... | grep -e ./cmd -e ./pkg` -v -coverprofile .testCoverage.txt
	go tool cover -func=.testCoverage.txt

#LINTERS
fmt:
	go fmt ./...

errcheck:
	errcheck ./...

lint:
	go fmt $$(go list ./... | grep -v ./vendor/)
	goimports -d -w $$(find . -type f -name '*.go' -not -path './vendor/*')
	golangci-lint run --skip-dirs tmp