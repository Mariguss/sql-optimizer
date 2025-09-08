FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /sql-optimizer ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /sql-optimizer .
COPY ./web ./web

EXPOSE 8080

CMD ["./sql-optimizer"]