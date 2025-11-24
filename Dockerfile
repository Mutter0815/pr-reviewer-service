FROM golang:1.24-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o pr-reviewerservice ./cmd/app

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache curl
COPY --from=builder /app/pr-reviewerservice .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/internal/transport/http/swagger ./internal/transport/http/swagger
EXPOSE 8080
CMD ["./pr-reviewerservice"]
