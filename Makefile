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

test:
	DATABASE_URL=root:pass@/test RUN_INTEGRATION_TESTS=true go test -v ./tests/...

twcss:
	tailwindcss -i ./web/styles/styles.css -o ./web/static/styles.css

dev:
	DATABASE_URL=root:pass@/test go run main.go

alpine:
	curl -o \
		web/static/alpine-focus.min.js \
		https://cdn.jsdelivr.net/npm/@alpinejs/focus@3.x.x/dist/cdn.min.js
	curl -o \
		web/static/alpine.min.js \
		https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js

htmx:
	curl -o \
		web/static/htmx.min.js \
		https://unpkg.com/htmx.org@1.9.6/dist/htmx.min.js

iconify:
	curl -o \
		web/static/iconify-icon.min.js \
		https://code.iconify.design/iconify-icon/1.0.7/iconify-icon.min.js

minmain:
	uglifyjs --module --webkit web/scripts/main.mjs -o web/static/main.min.js

postcss:
	bunx postcss web/styles/styles.css -o web/static/styles.min.css

bunmain:
	bun build web/scripts/main.js \
		--outdir web/static \
		--minify-whitespace \
		--minify-syntax \
		--entry-naming "[dir]/[name].min.[ext]"
