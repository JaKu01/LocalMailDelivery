FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd

RUN go build -o imap

EXPOSE 143

CMD ["./imap"]