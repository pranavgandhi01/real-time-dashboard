FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY ./backend/services/mock-data-service .
RUN CGO_ENABLED=0 go build -o mock-data-service .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/mock-data-service .
COPY ./docs/openskyapi_response.json ./docs/

EXPOSE 8083
CMD ["./mock-data-service"]