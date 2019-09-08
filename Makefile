deps:
		go mod vendor

lint:
		golangci-lint run --config=.golangci.yml ./...

test:	lint
		go test -cover -failfast ./...