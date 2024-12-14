FROM golang:1.22-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod tidy
COPY . .
RUN go build -o main cmd/app/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 80
CMD ["./main"]