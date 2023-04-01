FROM golang:1.20-alpine

WORKDIR /app

COPY . /app

RUN go mod init && go mod tidy

EXPOSE 8080

CMD go run .