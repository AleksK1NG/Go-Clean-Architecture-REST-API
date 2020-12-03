### Golang REST API example 🚀
I write my first [Go](https://golang.org/) project, i try to create example REST API similar to real world production projects using [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html), [Postgresql](https://www.postgresql.org/) for a database, [Redis](https://redis.io/) for sessions and caching,
[AWS S3](https://aws.amazon.com/) images uploading, logging, [Jaeger](https://www.jaegertracing.io/) [OpenTracing](https://opentracing.io/), [Prometheus](https://prometheus.io/) and [Grafana](https://grafana.com/) metrics,
[Swagger](https://swagger.io/about/) documentation, [Docker](https://www.docker.com/) infrastructure for development.
<br/>

For web framework i have chosen most popular for today [Echo](https://github.com/labstack/echo), another possible good alternatives [Gin](https://github.com/gin-gonic/gin), [Chi](https://github.com/go-chi/chi), [Mux](https://github.com/gorilla/mux).
<br/>
SQL db standard solution on my opinion is combination of [sqlx](https://github.com/jmoiron/sqlx) and [pgx](https://github.com/jackc/pgx).
<br/>
Redis for sessions and caching [go-redis](https://github.com/go-redis/redis)
<br/>
[Viper](https://github.com/spf13/viper) for configuration.
<br/>
[Zap](https://github.com/uber-go/zap) for logger, but this two is good too [zerolog](https://github.com/rs/zerolog) and [logrus](https://github.com/sirupsen/logrus).
<br/>
Swagger for api documentation, i used [swag](https://github.com/swaggo/swag), another good one is [go-swagger](https://github.com/go-swagger/go-swagger)
<br/>
[MinIO](https://min.io/) for AWS S3 good variant [minio-go](https://github.com/minio/minio-go)
<br/>
[Validator](https://github.com/go-playground/validator) is good solution for validation.
<br/>
For testing and mocking [testify](https://github.com/stretchr/testify) and [gomock](https://github.com/golang/mock) very good tools.
<br/>
Session and Token Based Authentication, i implement both for this example project, used [jwt-go](https://github.com/dgrijalva/jwt-go) for JWT.
<br/>



#### 👨‍💻 Full list what has been used:
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

#### Recomendation for local development most comfortable usage:
    make local // run all containers
    make run // it's easier way to attach debugger or rebuild/rerun project

#### 🙌👨‍💻🚀 Docker-compose files:
    docker-compose.local.yml - run postgresql, redis, aws, prometheus, grafana containrs
    docker-compose.dev.yml - run docker development environment
    docker-compose.delve.yml run development environment with delve debug

### Docker development usage:
    make docker

### Local development usage:
    make local
    make run

### SWAGGER UI:

https://localhost:5000/swagger/index.html

### Jaeger UI:

http://localhost:16686

### Prometheus UI:

http://localhost:9090

### Grafana UI:

http://localhost:3000