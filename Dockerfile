# BUILD
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux go build -o main .

# RUN
FROM alpine:3.19

RUN adduser -D appuser
USER appuser

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 12345

CMD ["./main"]
