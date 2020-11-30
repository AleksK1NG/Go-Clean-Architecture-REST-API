.PHONY: migrate migrate_down migrate_up migrate_version docker prod docker_delve local swaggo test

# ==============================================================================
# Go migrate postgresql

force:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations force 1

version:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations version

migrate_up:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations up 1

migrate_down:
	migrate -database postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable -path migrations down 1


# ==============================================================================
# Docker compose commands

docker:
	echo "Starting docker environment"
	docker run -d --name jaeger \
					   -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
					   -p 5775:5775/udp \
					   -p 6831:6831/udp \
					   -p 6832:6832/udp \
					   -p 5778:5778 \
					   -p 16686:16686 \
					   -p 14268:14268 \
					   -p 14250:14250 \
					   -p 9411:9411 \
					   jaegertracing/all-in-one:1.21
	docker-compose -f docker-compose.dev.yml up --build

docker_delve:
	echo "Starting docker debug environment"
	docker run -d --name jaeger \
    			   -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
    			   -p 5775:5775/udp \
    			   -p 6831:6831/udp \
    			   -p 6832:6832/udp \
    			   -p 5778:5778 \
    			   -p 16686:16686 \
    			   -p 14268:14268 \
    			   -p 14250:14250 \
    			   -p 9411:9411 \
    			   jaegertracing/all-in-one:1.21
	docker-compose -f docker-compose.delve.yml up --build

prod:
	echo "Starting docker prod environment"
	docker run -d --name jaeger \
    			   -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
    			   -p 5775:5775/udp \
    			   -p 6831:6831/udp \
    			   -p 6832:6832/udp \
    			   -p 5778:5778 \
    			   -p 16686:16686 \
    			   -p 14268:14268 \
    			   -p 14250:14250 \
    			   -p 9411:9411 \
    			   jaegertracing/all-in-one:1.21
	docker-compose -f docker-compose.prod.yml up --build

local:
	echo "Starting local environment"
	docker run -d --name jaeger \
			   -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
			   -p 5775:5775/udp \
			   -p 6831:6831/udp \
			   -p 6832:6832/udp \
			   -p 5778:5778 \
			   -p 16686:16686 \
			   -p 14268:14268 \
			   -p 14250:14250 \
			   -p 9411:9411 \
			   jaegertracing/all-in-one:1.21
	docker-compose -f docker-compose.local.yml up --build


jaeger:
	echo "Starting jaeger containers"
	docker run --name jaeger \
      -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
      -p 5775:5775/udp \
      -p 6831:6831/udp \
      -p 6832:6832/udp \
      -p 5778:5778 \
      -p 16686:16686 \
      -p 14268:14268 \
      -p 14250:14250 \
      -p 9411:9411 \
      jaegertracing/all-in-one:1.21


# ==============================================================================
# Tools commands

run-linter:
	echo "Starting linters"
	golangci-lint -c .golangci.yml run ./...

swaggo:
	echo "Starting swagger generating"
	swag init -g **/**/*.go


# ==============================================================================
# Main

run:
	go run ./cmd/api/main.go

build:
	go build ./cmd/api/main.go

test:
	go test -cover ./...


# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache


# ==============================================================================
# Docker support

FILES := $(shell docker ps -aq)

down-local:
	docker stop $(FILES)
	docker rm $(FILES)

clean:
	docker system prune -f

logs-local:
	docker logs -f $(FILES)