FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download || true
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /kencat-api ./cmd/api

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /kencat-api /app/kencat-api
COPY . /app
EXPOSE 8080
CMD ["/app/kencat-api"]
