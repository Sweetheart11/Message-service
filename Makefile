build:
	@go build -o bin/timetracker cmd/main.go 

run: build
	@sudo docker-compose up -d
	@$(MAKE) migrations_up
	@./bin/timetracker --config=.env 

test:
	@go test ./... -v

start_db:
	@sudo docker-compose up -d

migrations_up:
	@cd migrations && goose postgres "host=localhost port=5432 \
	user=postgres password=example dbname=postgres sslmode=disable" up

migrations_down:
	@cd migrations && goose postgres "host=localhost port=5432 \
	user=postgres password=example dbname=postgres sslmode=disable" down

migrations_status:
	@cd migrations && goose postgres "host=localhost port=5432 \
	user=postgres password=example dbname=postgres sslmode=disable" status

