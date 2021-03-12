FROM golang:1.14
WORKDIR /app
ENV GO111MODULE=on
COPY ./ ./
RUN GOMAXPROCS=4 go test -timeout 90s -race ./...
RUN go get -u github.com/pefish/go-build-tool@0.0.3
RUN make
ENV GO_CONFIG /app/config/pom.yaml
ENV GO_SECRET /app/secret/pom.yaml
CMD ["./build/bin/linux/main", "--help"]

# docker build -t pefish/telegram-bot-manager:v1.2.4 .