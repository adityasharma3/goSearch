FROM golang:1.23-alpine
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY ./cmd/search/.env ./

RUN go build -o main ./cmd/search/main.go

EXPOSE 8080
CMD ["./main"]