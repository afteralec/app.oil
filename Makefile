ifneq (,$(wildcard ./.env))
    include .env
    export
endif

env:
	cp .env.example .env

run:
	go run cmd/main/main.go

setup_test_db:
	chmod u+x ./scripts/setup_test_db.sh
	./scripts/setup_test_db.sh

fetch_migrations:
	chmod u+x ./scripts/fetch_migrations.sh
	./scripts/fetch_migrations.sh

cleanup_migrations:
	chmod u+x ./scripts/cleanup_migrations.sh
	./scripts/cleanup_migrations.sh

twcss:
	tailwindcss -i ./web/styles/styles.css -o ./web/static/styles.css

dev:
	DATABASE_URL=root:pass@/test go run main.go
