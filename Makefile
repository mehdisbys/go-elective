deps:
		go mod vendor

test:
		go test -cover -failfast ./...