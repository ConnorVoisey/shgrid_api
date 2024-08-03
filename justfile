list:
    just --list

dev:
    docker compose up -d
    air

pre-commit:
    go fmt ./...
    golangci-lint run
    go test ./...
    curl 0.0.0.0:3000/openapi.json | jq > openapi.json

init:
    go install github.com/air-verse/air@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
    go install -tags 'postgres' 'github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1'
    go install github.com/danielgtaylor/restish@latest
    cp -i .env.example .env
    restish api configure shgrid_api :3000

seed:
    go run cmd/seed/main.go
