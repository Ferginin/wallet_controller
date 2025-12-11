FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /app

COPY go.mod go.sum config.env ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /wallet_controller ./main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /wallet_controller .

EXPOSE 8080

CMD ["./wallet_controller"]
