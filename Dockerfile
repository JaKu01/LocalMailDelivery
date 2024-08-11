FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd

RUN go build -o imap

EXPOSE 143
EXPOSE 8080

WORKDIR /app

CMD ["./cmd/imap", "test", "test"]