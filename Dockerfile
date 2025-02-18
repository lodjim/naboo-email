FROM golang:1.22.4

WORKDIR /app

COPY . .

RUN go mod download

RUN go build cmd/main.go

EXPOSE 50051

ENTRYPOINT ["./main"]
