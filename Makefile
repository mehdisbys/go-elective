deps:
		export GO111MODULE=on ; go mod vendor

lint:
		golangci-lint run --config=.golangci.yml ./...

test:	lint
		go test -cover -failfast ./...

linux-binary:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main -ldflags "-w -s"