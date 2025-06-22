FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY ./backend/pkg ./pkg
COPY go.mod go.sum ./
RUN go mod download
COPY ./backend/services/api-gateway .
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/api-gateway .
EXPOSE 8080
USER 1000:1000
CMD ["./api-gateway"]