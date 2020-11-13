# Initial stage: download modules
FROM golang:1.15.5 as modules

ADD ./go.mod go.sum /m/
RUN cd /m && go mod download
RUN go get github.com/githubnemo/CompileDaemon

# Intermediate stage: Build the binary
FROM golang:1.15.5 as builder

COPY --from=modules /go/pkg /go/pkg

ENV config=docker

# add a non-privileged user
RUN useradd -u 10001 myapp

RUN mkdir -p /app
ADD docker /app
WORKDIR /app

EXPOSE 5000
EXPOSE 5555
EXPOSE 7070

# Build the binary with go build
RUN GOOS=linux GOARCH=amd64

ENTRYPOINT CompileDaemon --build="go build cmd/api/main.go" --command=./main




