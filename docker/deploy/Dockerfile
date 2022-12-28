FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o bin/app cmd/app/main.go

EXPOSE 3000

CMD ["./bin/app"]
