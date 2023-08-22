.PHONY: clean critic security lint test build run
include .env

APP_NAME = main
BUILD_DIR = $(PWD)/build
PSQL=postgres://${DB_USER}:${DB_SECRET}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL}

clean:
	rm -rf ./build

critic:
	gocritic check -enableAll ./...

security:
	tmp/gosec ./...

lint:
	tmp/golangci-lint run ./...

test: clean critic security lint
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

# change this
build: clean critic security lint
# build: test
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) main.go

run: build
	tmp/air

docker.run: docker.network docker.chi

docker.network:
	docker network inspect dev-network >/dev/null 2>&1 || \
	docker network create -d bridge dev-network

docker.chi.build:
	docker build -t chi .

docker.chi: docker.chi.build
	docker run --rm -d \
		--name gochi \
		--network dev-network \
		-p 8000:8000 \
		chi

docker.stop: docker.stop.chi

docker.stop.chi:
	docker stop gochi

cmgr:
	tmp/migrate create -ext sql -dir db/migrations -seq ${name}

migup:
	tmp/migrate -path db/migrations -database "${PSQL}" -verbose up

migdown:
	tmp/migrate -path db/migrations -database "${PSQL}" -verbose down

setupair:
	curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b tmp/

setupmigrate:
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz -C tmp/

setupgosec:
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b tmp/ v2.16.0

setupgocritic:
	go install -v github.com/go-critic/go-critic/cmd/gocritic@latest

setuplint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b tmp/ v1.54.1

setup: setupair setupmigrate setupgosec setupgocritic setuplint