.PHONY: migrate demigrate goosestatus connect

migrate:
	@cd sql/schema/ && goose postgres "postgres://user:mysecretpassword@localhost:5432/test" up

demigrate:
	@cd sql/schema/ && goose postgres "postgres://user:mysecretpassword@localhost:5432/test" down

goosestatus:
	@cd sql/schema/ && goose postgres "postgres://user:mysecretpassword@localhost:5432/test" status

connect:
	@docker-compose exec goddess psql -U user test

up: 
	@docker-compose up -d

kill: 
	@docker-compose down

run:
	go run .