PROJECT := star
POSTGRES_URL := postgres://superuser:secret@localhost:5432/test?sslmode=disable
DEV_URL := docker+postgres://_/postgres:alpine/dev
SCHEMA_FILE := store/schema/schema.sql
FILE_SCHEMA_FILE := file://store/schema/schema.sql
SCHEMA_DIR := file://store/schema

help:
	@echo "Available commands:"
	
	@echo "\nAtlas"
	@echo "  ainspect       - Inspect schema and output to $(SCHEMA_FILE)"
	@echo "  aapply         - Apply schema to database using $(FILE_SCHEMA_FILE) and $(DEV_URL)"
	@echo "  apush          - Push migrations to $(DEV_URL) from $(SCHEMA_DIR)"
	@echo "  ahash          - Hash migrations from $(SCHEMA_DIR)"
	@echo "  aweb           - Inspect schema using web interface"
	@echo "  all            - Run all commands: ainspect, aapply, apush, ahash, aweb"

	@echo "\nGoose"	
	@echo "  gup            - Apply Goose migrations"
	@echo "  gdo            - Rollback Goose migrations"
	@echo "  gdown          - Rollback Goose migrations to version 0"
	@echo "  gstatus        - Display the status of Goose migrations"
	@echo "  help           - Display this help message"

ainspect:
	atlas schema inspect \
		-u $(POSTGRES_URL) \
		--format '{{ sql . }}' > $(SCHEMA_FILE)
aapply:
	atlas schema apply \
		-u $(POSTGRES_URL) \
		--to $(FILE_SCHEMA_FILE) \
		--dev-url "$(DEV_URL)"
apush:
	atlas migrate push $(PROJECT) \
		--dev-url "$(DEV_URL)" \
		--dir $(SCHEMA_DIR)
ahash:
	atlas migrate hash --dir $(SCHEMA_DIR)
aweb:
	atlas schema inspect \
		-u $(POSTGRES_URL) \
		--web
all: schema-inspect schema-apply migrate-push migrate-hash schema-inspect-web

gup:
	goose -dir store/schema/goose postgres $(POSTGRES_URL) up
gdo:
	goose -dir store/schema/goose postgres $(POSTGRES_URL) down
gdown:
	goose -dir store/schema/goose postgres $(POSTGRES_URL) down-to 0
gstatus:
	goose -dir store/schema/goose postgres $(POSTGRES_URL) status
