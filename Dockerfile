FROM golang:latest AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /go/src
COPY go.mod .
RUN go mod download
COPY . .
RUN go build cmd/main.go

FROM golang:latest
COPY --from=builder /go/src .
EXPOSE 8080
ENTRYPOINT ["./main"]
