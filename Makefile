include .env
build:
	@go build -o bin/message-processor cmd/main.go 

run: build
	@sudo docker-compose up -d
	@$(MAKE) migrations_up
	@./bin/message-processor --config=.env 

test:
	@go test ./... -v

start_db:
	@sudo docker-compose up -d

migrations_up:
	@cd migrations && goose postgres "host=${DB_HOST} port=${DB_PORT} \
	user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=${DB_SSLMODE}" up

migrations_down:
	@cd migrations && goose postgres "host=${DB_HOST} port=${DB_PORT} \
	user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=${DB_SSLMODE}" down

migrations_status:
	@cd migrations && goose postgres "host=${DB_HOST} port=${DB_PORT} \
	user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=${DB_SSLMODE}" status

