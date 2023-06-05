# Stage 1: Download dependencies
FROM golang:1.19 AS builder

WORKDIR /go/src
COPY go.mod .
COPY go.sum .
ARG GITHUB_TOKEN
ENV GO_ENVIRONMENT="production"
ENV CONF_DIR=/go/src/conf
ENV GOPRIVATE=github.com/rromero96/*
RUN go env -w GOPRIVATE=github.com/rromero96/*
RUN git config --global url."https://github.com/".insteadOf git://github.com/
RUN git config --global credential.helper 'store --file ~/.git-credentials'
RUN echo "https://github.com:${GITHUB_TOKEN}@github.com" >> ~/.git-credentials
RUN GIT_TERMINAL_PROMPT=1 go mod download -x

# Stage 2: Build the application
FROM builder AS final

COPY . .
RUN go build -o app cmd/api/main.go

EXPOSE 8080

# Set the entrypoint to run the application
ENTRYPOINT ["./app"]

