FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o user_app ./cmd/...

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/user_app .

EXPOSE 8888

CMD ["./user_app"]