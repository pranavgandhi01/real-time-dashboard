FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY ./backend/pkg ./pkg
COPY ./backend/services/flight-data-service/go.mod ./backend/services/flight-data-service/go.sum ./
RUN go mod download
COPY ./backend/services/flight-data-service .
RUN CGO_ENABLED=0 GOOS=linux go build -o flight-data-service .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/flight-data-service .
EXPOSE 8081
USER 1000:1000
CMD ["./flight-data-service"]