### Golang REST API example 🚀🤩🙌👍
Golang REST API Clean Architecture similar to production example, using Postgresql for a database, Redis for sessions and caching,
AWS S3 images uploading, logging, Jaeger OpenTracing, Prometheus and Grafana metrics,
Swagger documentation, Docker infrastructure for development.

#### 👨‍💻 Used:

* [echo](https://github.com/labstack/echo) - Web framework
* [sqlx](https://github.com/jmoiron/sqlx) - Extensions to database/sql.
* [pgx](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit for Go
* [viper](https://github.com/spf13/viper) - Go configuration with fangs
* [go-redis](https://github.com/go-redis/redis) - Type-safe Redis client for Golang
* [zap](https://github.com/uber-go/zap) - Logger
* [validator](https://github.com/go-playground/validator) - Go Struct and Field validation
* [jwt-go](https://github.com/dgrijalva/jwt-go) - JSON Web Tokens (JWT)
* [uuid](https://github.com/google/uuid) - UUID
* [migrate](https://github.com/golang-migrate/migrate) - Database migrations. CLI and Golang library.
* [minio-go](https://github.com/minio/minio-go) - AWS S3 MinIO Client SDK for Go
* [bluemonday](https://github.com/microcosm-cc/bluemonday) - HTML sanitizer
* [swag](https://github.com/swaggo/swag) - Swagger
* [testify](https://github.com/stretchr/testify) - Testing toolkit
* [gomock](https://github.com/golang/mock) - Mocking framework
* [CompileDaemon](https://github.com/githubnemo/CompileDaemon) - Compile daemon for Go
* [Docker](https://www.docker.com/) - Docker
* [opentracing](https://github.com/opentracing/opentracing-go) - OpenTracing API for Go
* [jaeger](https://github.com/jaegertracing/jaeger-client-go) - Jaeger Bindings for Go OpenTracing API.
* [jaeger-lib](https://github.com/jaegertracing/jaeger-lib) - Different components of Jaeger

#### 🙌👨‍💻🚀 Docker-compose files:

    docker-compose.local.yml - run postgresql, redis, aws, prometheus, grafana containrs
    docker-compose.dev.yml - run docker development environment
    docker-compose.delve.yml run development environment with delve debug

### Local development usage:

    make local
    make run

### Docker development usage:

    make docker

### SWAGGER UI:

https://localhost:5000/swagger/index.html

### Jaeger UI:

http://localhost:16686

### Prometheus UI:

http://localhost:9090

### Grafana UI:

http://localhost:3000