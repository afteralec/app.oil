build:
	go build -o bin/ptrcr main.go

test:
	SERVER_READ_TIMEOUT=10 DATABASE_URL=root:pass@/test?parseTime=true SENDING_STONE_URL=localhost:8010 REDIS_ADDR=127.0.0.1:6379 PETRICHOR_PLAY_URL=http://localhost:5173 DISABLE_SENDING_STONE=true DISABLE_CSRF=true DISABLE_LOGGING=true go test -v ./...
	npm test

test-ci:
	SERVER_READ_TIMEOUT=10  REDIS_ADDR=127.0.0.1:6379 DISABLE_SENDING_STONE=true DISABLE_CSRF=true DISABLE_LOGGING=true go test -v ./...
	npm test

dev:
	BASE_URL=http://localhost:8008 PETRICHOR_PLAY_URL=http://localhost:5173 SERVER_READ_TIMEOUT=10 DATABASE_URL=root:pass@/test?parseTime=true SENDING_STONE_URL=localhost:8010 REDIS_ADDR=127.0.0.1:6379 DISABLE_SENDING_STONE=false go run main.go run

alpine:
	curl -o \
		web/static/alpine-focus.min.js \
		https://cdn.jsdelivr.net/npm/@alpinejs/focus@3.x.x/dist/cdn.min.js
	curl -o \
		web/static/alpine-morph.min.js \
		https://cdn.jsdelivr.net/npm/@alpinejs/morph@3.x.x/dist/cdn.min.js
	curl -o \
		web/static/alpine.min.js \
		https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js

htmx:
	curl -o \
		web/static/htmx.min.js \
		https://unpkg.com/htmx.org@1.9.6/dist/htmx.min.js
	curl -o \
		web/static/htmx-morph.js \
		https://unpkg.com/htmx.org@1.9.6/dist/ext/alpine-morph.js

icons:
	curl -o \
		web/static/iconify-icon.min.js \
		https://code.iconify.design/iconify-icon/1.0.7/iconify-icon.min.js

mainjs:
	bun build web/scripts/main.js \
		--outdir web/static \
		--minify-whitespace \
		--minify-syntax \
		--entry-naming "[dir]/[name].min.[ext]"

js:
	make alpine
	make htmx
	make icons
	make mainjs

css:
	npx postcss web/styles/styles.css -o web/static/styles.min.css

redis:
	docker run --name app-redis -p 6379:6379 -d --rm redis

proto:
	git clone git@github.com:petrichormud/proto.git ./proto

protoc:
	protoc --proto_path=proto proto/*.proto --go_out=./ --go-grpc_out=.
