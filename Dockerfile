FROM golang:alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/server

# -------------------------------------------------------------------------------

FROM alpline:3.21

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]
