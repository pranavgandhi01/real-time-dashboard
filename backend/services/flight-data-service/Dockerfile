FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download
COPY ./backend ./
RUN CGO_ENABLED=0 GOOS=linux go build -o flight-data-service ./services/flight-data-service/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/flight-data-service .
EXPOSE 8081
USER 1000:1000
CMD ["./flight-data-service"]