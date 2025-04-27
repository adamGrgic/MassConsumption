include .env
export $(shell sed 's/=.*//' .env)
MAKEFLAGS += --no-print-directory

# MAIN
run:
	air

import:
	@go run ./cmd/import/main.go

debug:
	dlv debug ./cmd/app/main.go

check-env:
	@go run ./cmd/check/main.go

kill-server:
	@pkill -f "air" || true
	@pkill -f "main" || true

# CLIENT
css:
	@./bin/css-compiler

css-watch:
	@./bin/css-compiler --watch

js:
	@./bin/javascript-compiler

js-watch:
	@./bin/javascript-compiler --watch

htmx:
	@go run ./cmd/htmx/main.go

templ:
	@templ generate

templ-watch:
	@clear
	@templ generate --watch

# BUNDLER
bundler:
	@bun build ./src/scripts/bundlers/css.ts --compile --outfile ./bin/css-compiler
	@bun build ./src/scripts/bundlers/javascript.ts --compile --outfile ./bin/javascript-compiler

# DOCKER
docker-build:
	docker build -t web-scraper .

docker-run:
	docker run -p 8080:8080 web-scraper

docker-dev:
	docker compose up --build


