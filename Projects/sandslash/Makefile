# Define variables for both PostgreSQL URLs
DEV_URL := docker+postgres://_/postgres:alpine/dev
SCHEMA_FILE := store/schema/schema.sql
FILE_SCHEMA_FILE := file://store/schema/schema.sql
SCHEMA_DIR := file://store/schema

# Command to list all defined targets with descriptions
list-commands:
	@echo "Available targets:"
	@echo "  schema-inspect-local:     		Inspect schema with local PostgreSQL URL"
	@echo "  schema-apply-local:       		Apply schema with local PostgreSQL URL"
	@echo "  migrate-push-local:       		Push migrations with local PostgreSQL URL"
	@echo "  migrate-hash-local:       		Hash migrations with local PostgreSQL URL"
	@echo "  schema-inspect-web-local: 		Inspect schema via web with local PostgreSQL URL"
	@echo "  schema-inspect-remote:    		Inspect schema with remote PostgreSQL URL"
	@echo "  schema-apply-remote:      		Apply schema with remote PostgreSQL URL"
	@echo "  migrate-push-remote:      		Push migrations with remote PostgreSQL URL"
	@echo "  migrate-hash-remote:      		Hash migrations with remote PostgreSQL URL"
	@echo "  schema-inspect-web-remote:		Inspect schema via web with remote PostgreSQL URL"
	@echo "  all-local:                		Run all steps for local PostgreSQL"
	@echo "  all-remote:               		Run all steps for remote PostgreSQL"

# Rule to inspect the schema with the local PostgreSQL URL
schema-inspect-local:
	atlas schema inspect \
		-u $(POSTGRES_URL_LOCAL) \
		--format '{{ sql . }}' > $(SCHEMA_FILE)

# Rule to apply the schema with the local PostgreSQL URL
schema-apply-local:
	atlas schema apply \
		-u $(POSTGRES_URL_LOCAL) \
		--to $(FILE_SCHEMA_FILE) \
		--dev-url "$(DEV_URL)"

# Rule to push migrations with the local PostgreSQL URL
migrate-push-local:
	atlas migrate push star \
		--dev-url "$(DEV_URL)" \
		--dir $(SCHEMA_DIR)

# Rule to hash migrations with the local PostgreSQL URL
migrate-hash-local:
	atlas migrate hash --dir $(SCHEMA_DIR)

# Rule to inspect schema via web with the local PostgreSQL URL
schema-inspect-web-local:
	atlas schema inspect \
		-u $(POSTGRES_URL_LOCAL) \
		--web

# Rule to inspect the schema with the remote PostgreSQL URL
schema-inspect-remote:
	atlas schema inspect \
		-u $(POSTGRES_URL_REMOTE) \
		--format '{{ sql . }}' > $(SCHEMA_FILE)

# Rule to apply the schema with the remote PostgreSQL URL
schema-apply-remote:
	atlas schema apply \
		-u $(POSTGRES_URL_REMOTE) \
		--to $(FILE_SCHEMA_FILE) \
		--dev-url "$(DEV_URL)"

# Rule to push migrations with the remote PostgreSQL URL
migrate-push-remote:
	atlas migrate push star \
		--dev-url "$(DEV_URL)" \
		--dir $(SCHEMA_DIR)

# Rule to hash migrations with the remote PostgreSQL URL
migrate-hash-remote:
	atlas migrate hash --dir $(SCHEMA_DIR)

# Rule to inspect schema via web with the remote PostgreSQL URL
schema-inspect-web-remote:
	atlas schema inspect \
		-u $(POSTGRES_URL_REMOTE) \
		--web

# Default target to run all steps for local PostgreSQL
all-local: schema-inspect-local schema-apply-local migrate-push-local migrate-hash-local schema-inspect-web-local

# Default target to run all steps for remote PostgreSQL
all-remote: schema-inspect-remote schema-apply-remote migrate-push-remote migrate-hash-remote schema-inspect-web-remote

