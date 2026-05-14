FROM golang:1.22-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o domloc ./cmd/domloc
RUN go test ./...

FROM ubuntu:22.04
COPY --from=builder /app/domloc /usr/local/bin/domloc
RUN apt-get update && apt-get install -y dnsmasq ca-certificates && rm -rf /var/lib/apt/lists/*
ENTRYPOINT ["domloc"]
CMD ["--version"]
