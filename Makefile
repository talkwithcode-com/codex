test:
	go test --coverprofile=coverage.out ./... -v -short

coverage:
	go tool cover -func=coverage.out  