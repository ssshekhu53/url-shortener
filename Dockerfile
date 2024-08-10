FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN GOOS=linux GOARCH=amd64 go build -o main

FROM alpine:3.15

COPY --from=builder /app/main /main

RUN chmod 777 /main

EXPOSE 9000

CMD ["/main"]
