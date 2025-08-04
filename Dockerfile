# Сперва сборка
FROM golang:1.24.1-alpine AS builder  

WORKDIR /app  

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/valservice main.go
RUN ls -l /app
 
# Потом образ
FROM alpine:latest
WORKDIR /app
RUN ls -l /app

COPY --from=builder /app/valservice /app/valservice
RUN ls -l /app

EXPOSE 8080
CMD ["/app/valservice"]  

