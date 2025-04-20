MAKEFLAGS += --no-print-directory

# MAIN
run:
	air

debug:
	dlv debug ./cmd/app/main.go

check-env:
	@go run ./cmd/check/main.go

kill-server:
	@pkill -f "air" || true
	@pkill -f "main" || true

build:
	@go run ./cmd/build/main.go

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

# DOCS
docs-go: 
	@swag init -g ./cmd/app/main.go -o ./docs

# BUNDLER
bundler:
	@bun build ./src/scripts/bundlers/css.ts --compile --outfile ./bin/css-compiler
	@bun build ./src/scripts/bundlers/javascript.ts --compile --outfile ./bin/javascript-compiler

## run the uncompiled versions of each bundler
bundler-js-test:
	@bun ./src/scripts/bundlers/javascript.ts

bundler-css-test:
	@bun ./src/scripts/bundlers/css.ts


# MIGRATIONS
# Config
MIGRATE_BIN=migrate
MIGRATIONS_DIR=./internal/db/migrations
DB_URL=postgres://postgres:psql@localhost:5432/goth-todo?sslmode=disable

# Usage:
# make migrate-create name=create_users_table
migrate-create:
	$(MIGRATE_BIN) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up:
	$(MIGRATE_BIN) -database "$(DB_URL)" -path $(MIGRATIONS_DIR) up

migrate-down:
	$(MIGRATE_BIN) -database "$(DB_URL)" -path $(MIGRATIONS_DIR) down

migrate-force:
	$(MIGRATE_BIN) -database "$(DB_URL)" -path $(MIGRATIONS_DIR) force $(version)

migrate-version:
	$(MIGRATE_BIN) -database "$(DB_URL)" -path $(MIGRATIONS_DIR) version

migrate-drop:
	$(MIGRATE_BIN) -database "$(DB_URL)" -path $(MIGRATIONS_DIR) drop -f

migrate-clean:
	@go run ./cmd/migrate/main.go --migrate-clean

migrate-seed:
	@go run ./cmd/seed/main.go

# DOCKER
docker-build:
	docker build -t goth-todo .

docker-run:
	docker run -p 8080:8080 goth-todo

docker-dev:
	docker compose up --build


