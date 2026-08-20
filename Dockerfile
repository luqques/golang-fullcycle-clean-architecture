FROM golang:1.23-bookworm AS builder

WORKDIR /app

RUN apt-get update \
    && apt-get install -y --no-install-recommends protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.2 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/orders.proto

RUN CGO_ENABLED=0 GOOS=linux go build -o /orders ./cmd/orders

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /orders /orders
COPY migrations ./migrations

EXPOSE 8000 8080 50051

ENTRYPOINT ["/orders"]
