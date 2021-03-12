
DEFAULT: build-cur

GORUN = go-build-tool


build-cur:
	make test
	$(GORUN)

test:
	mkdir -p mock/mock-go-http
	mockgen github.com/pefish/go-http IHttp > mock/mock-go-http/i_http.go
	go test -cover ./...
