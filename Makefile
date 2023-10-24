ifneq (,$(wildcard ./.env))
    include .env
    export
endif

run:
	go run cmd/main/main.go

sqlc:
	./scripts/fetch_migrations.sh
	sqlc generate
	./scripts/cleanup_migrations.sh

setup_test_db:
	./scripts/setup_test_db.sh

fetch_migrations:
	./scripts/fetch_migrations.sh

cleanup_migrations:
	./scripts/cleanup_migrations.sh

twcss:
	tailwindcss -i ./web/styles/styles.css -o ./web/static/styles.css

dev:
	DATABASE_URL=root:pass@/test go run main.go
