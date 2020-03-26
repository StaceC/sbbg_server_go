# Simple Browser Based Game - Go Server

The demo Vue server can be built using **Docker** or your **Local** development environment.

Using either method below, the server should become available on:

http://localhost:8089

*Please ensure you have no other services running on 8089 on your machine*

## Docker

With docker installed on your systems command line do the following:

#### Build Docker

```
docker-compose build
```

#### Run Docker

```
docker-compose up
```

## Local

Without docker, you can build locally if golang 1.1+ is installed on your system

#### Run Local

```
go run ./cmd/sbbg
```

#### Test Local
```
go test ./pkg/[NAME]
```

#### Example Engine Test

```
go test ./pkg/engine
```
