FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go build -o bin/message-processor cmd/main.go

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

CMD ["sh", "-c", "goose -dir ./migrations postgres \"host=$DB_HOST port=$DB_PORT user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=$DB_SSLMODE\" up && ./bin/message-processor --config=.env"]