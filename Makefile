build:
	go build -o bin ./cmd

start: build
	./bin

live_reload:
	air

tidy:
	go mod tidy

vendor:
	go mod vendor

install: tidy vendor

migrate:
	@echo "Migrating database"
	migrate -path internal/pkg/postgres/migration -database $(url) -verbose up
	@echo "Migrate database completed"

migrate_down:
	@echo "Migrating database down"
	migrate -path internal/pkg/postgres/migration -database $(url) -verbose down
	@echo "Migrate database down completed"

migrate_down_count:
	@echo "Migrating database down"
	migrate -path internal/pkg/postgres/migration -database $(url) -verbose down $(count)
	@echo "Migrate database down completed"

migrate_force:
	@echo "Forcing migration"
	migrate -path internal/pkg/postgres/migration -database $(url) -verbose force $(version)
	@echo "Migration forced"


migrateup:
	@echo "Migrating up"
	$(MAKE) migrate url=$(url)
	@echo "Migrate up completed"

migrate_ci:
	@echo "Migrating CI database"
	$(MAKE) migrate url=$(url)
	@echo "Migrate CI database completed"

seed_db:
	@echo "Seeding database from seed.sql"
	psql $(url) -f internal/pkg/postgres/seed.sql
	@echo "Database seeded from seed.sql"

migratedown:
	@echo "Migrating down"
	$(MAKE) migrate_down url=$(url)
	@echo "Migrate down completed"

forcedown:
	@echo "Forcing down migration"
	$(MAKE) migrate_force url=$(url) version=$(version)
	@echo "Forced down migration completed"

force_fix_migration:
	@echo "Forcing migration fix"
	$(MAKE) migrate_force url=$(url) version=$(version)
	@echo "Migration fixed"

sqlc:
	@echo "Generating sqlc"
	sqlc generate
	@echo "Sqlc generated"

generate_migration:
	migrate create -ext sql -dir internal/pkg/postgres/migration $(name)

generate_32_bit_key:
	openssl rand -base64 32

static_analysis:
	@echo "Running go vet ./..."
	go vet ./...
	@echo "Running staticcheck ./..."
	staticcheck ./...
	@echo "Running errcheck ./..."
	# errcheck ./...
	@echo "static analysis done"


.PHONY: build start tidy vendor clean migrate migrate_down migrate_force migrateup migrate_ci migratedown forcedown force_fix_migration sqlc generate_migration generate_32_bit_key static_analysis
