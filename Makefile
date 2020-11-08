.PHONY: migrate migrate_down migrate_up migrate_version debug hot_reload prod delve check_install, swagger, local

force:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations force 1

version:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations version

migrate_up:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations up 1

migrate_down:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations down 1

hot_reload:
	 docker-compose -f docker-compose.HotReload.yml up --build

delve:
	 docker-compose -f docker-compose.DelveHotReload.yml up --build

debug:
	 docker-compose -f docker-compose.debug.yml up --build

prod:
	 docker-compose -f docker-compose.prod.yml up --build


check_install:
	which swagger || GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: check_install
	swagger generate spec -o ./swagger/swagger.yaml --scan-models

# linter
run-linter:
	golangci-lint -c .golangci.yml run ./...

local:
	 docker-compose -f docker-compose.local.yml up --build

swaggo:
	swag init -g **/**/*.go