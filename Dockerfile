FROM golang:1.23.3-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY migrations ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go
RUN go build -o migrate ./cmd/migrator/migrate.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY --from=builder /app/configForDocker.env .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

# there's three teams:
# - first we wait on purpose, this is to allow postgres in docker to wake up
# - next, we make migrations to the database, this will happen every time the container is started, 
# 	but in fact, the migration will take place only the first time to create tables, and then
#   the binary for migrations will respond that there are no new migrations to apply
# - and at the end, the application itself will run
CMD ["sh", "-c", "sleep 3 && ./migrate --storage-path postgres://evans:evans@storage:8081/postgres?sslmode=disable --migrations-path ./migrations && ./main -config ./configForDocker.env"]