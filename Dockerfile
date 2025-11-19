FROM golang:1.24-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o pr-reviewerservice ./cmd/app

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/pr-reviewerservice .
EXPOSE 8080
CMD ["./pr-reviewerservice"]