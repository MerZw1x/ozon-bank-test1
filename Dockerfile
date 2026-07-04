FROM golang:1.26.1-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./src/cmd/main.go

EXPOSE 80

CMD ["./app"] 