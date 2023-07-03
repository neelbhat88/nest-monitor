.PHONY: build application test

default: build

all: clean run

clean: 
	go clean

run: db-up build
	./application

mocks:
	go generate ./...

application: cmd/*.go
	go mod tidy -compat=1.19
	go build -o $@ $^

build: mocks application

db-up:
	docker-compose up -d nest-monitor-db

create-migration:
	@echo "Please enter a filename..."; \
	read FILE; \
	migrate create -ext sql -dir internal/data/postgres/migrations -seq $$FILE;

test:
	time go test -v ./...

start: docker
	docker-compose up

docker: mocks
	docker build --progress=plain . --tag nest-monitor --build-arg SSH_PRIVATE_KEY="`cat $$SSH_KEY_FILEPATH`"

migrate-up:
	 migrate -path internal/data/postgres/migrations -database 'postgres://nest-monitor:nest-monitor@localhost:5432/nest-monitor?sslmode=disable&search_path=nest-monitor' up

migrate-down:
	 migrate -path internal/data/postgres/migrations -database 'postgres://nest-monitor:nest-monitor@localhost:5432/nest-monitor?sslmode=disable&search_path=nest-monitor' down 1