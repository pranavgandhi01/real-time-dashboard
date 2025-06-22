FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download
COPY ./backend ./
RUN CGO_ENABLED=0 GOOS=linux go build -o mock-data-service ./services/mock-data-service/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/mock-data-service .
EXPOSE 8083
USER 1000:1000
CMD ["./mock-data-service"]