run:
	go run cmd/main/main.go

test:
	DATABASE_URL=root:pass@/test RUN_INTEGRATION_TESTS=true go test -v ./tests/...

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

icons:
	curl -o \
		web/static/iconify-icon.min.js \
		https://code.iconify.design/iconify-icon/1.0.7/iconify-icon.min.js

main:
	bun build web/scripts/main.js \
		--outdir web/static \
		--minify-whitespace \
		--minify-syntax \
		--entry-naming "[dir]/[name].min.[ext]"

js:
	make alpine
	make htmx
	make icons
	make main

css:
	bunx postcss web/styles/styles.css -o web/static/styles.min.css

redis:
	docker run --name app-redis -p 6379:6379 -d --rm redis
