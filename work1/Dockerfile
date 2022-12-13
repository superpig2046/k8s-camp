FROM golang:1.16-alpine
WORKDIR /build
ENV GOPROXY https://goproxy.cn
COPY . .
RUN \
    go mod download && \
    CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -o httpserver .
ENTRYPOINT ["/build/httpserver"]