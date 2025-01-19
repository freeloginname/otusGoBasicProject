FROM golang:1.23.1 AS builder

WORKDIR /build/nodes/

COPY .dockerignore ./

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server_api ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /build/nodes/server_api ./
COPY --from=builder /build/nodes/.env ./
COPY --from=builder /build/nodes/sql/migrations ./sql/migrations
COPY --from=builder /build/nodes/templates ./templates
COPY --from=builder /build/nodes/static ./static

ARG CONFIG_ENV=".env"
ENV CONFIG_ENV=$CONFIG_ENV

EXPOSE 8080
EXPOSE 5432

CMD ["./server_api"]
