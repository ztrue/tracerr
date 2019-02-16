default: test

.PHONY: lint
lint:
	go fmt . && \
		go fmt ./examples && \
		go vet && \
		golint $$(go list ./...)

.PHONY: doc
doc:
	@echo GoDoc link: http://localhost:6060/pkg/github.com/ztrue/tracerr
	godoc -http=:6060

.PHONY: test
test:
	go test -cover -v

.PHONY: coverage
coverage:
	go test -coverprofile=coverage.out && \
	go tool cover -func=coverage.out && \
	go tool cover -html=coverage.out

.PHONY: bench
bench:
	GOMAXPROCS=1 go test -bench=. -benchmem
