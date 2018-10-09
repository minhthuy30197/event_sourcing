# Gin template REST API

## CSDL

- PostgreSQL: Cài qua docker `docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=123456 -d postgres:10-alpine`

- TimescaleDB: Cài qua docker `docker run -d --name timescaledb -p 5433:5432 -e POSTGRES_PASSWORD=123456 timescale/timescaledb`

## Run

```go
  go run main.go
```
