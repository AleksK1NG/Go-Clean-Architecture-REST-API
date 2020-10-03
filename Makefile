.PHONY: migrate migrate_down migrate_up migrate_version compose_debug compose_hot__reload compose_prod compose_dev_db compose_dlv_reload

force:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations force 1

version:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations version

migrate_up:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations up 1

migrate_down:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations down 1

compose_hot__reload:
	 docker-compose -f docker-compose.HotReload.yml up --build

compose_dlv_reload:
	 docker-compose -f docker-compose.DelveHotReload.yml up --build

compose_debug:
	 docker-compose -f docker-compose.debug.yml up --build

compose_prod:
	 docker-compose -f docker-compose.prod.yml up --build

compose_dev_db:
	 docker-compose -f docker-compose.db.yml up --build