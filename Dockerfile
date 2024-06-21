FROM golang:latest

LABEL maintainer="Ümit Demir <umitde296@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main cmd/main/main.go

EXPOSE 8080

CMD ["./main"]