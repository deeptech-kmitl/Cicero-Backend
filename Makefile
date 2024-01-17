DB_URL=postgres://cicero:123456@localhost:4444/cicero_db_dev?sslmode=disable
PATH_MIGRATE ?= file://C:/go/Cicero-Backend/pkg/databases/migrations
PROJECT_ID ?= blabla
IMAGE_NAME ?= cicero_api

run_dev:
	air -c .air.dev.toml

run_prod:
	air -c .air.prod.toml

init_db:
    docker run --name cicero_db_dev -e POSTGRES_USER=cicero -e POSTGRES_PASSWORD=123456 -p 4444:5432 -d postgres:alpine

into_db:
    docker exec -it cicero_db_dev bash -c 'psql -U cicero'

create_db:
    docker exec -it cicero_db_dev bash -c 'psql -U cicero -c "CREATE DATABASE cicero_db_dev;"'

drop_db:
    docker exec -it cicero_db_dev bash -c 'psql -U cicero -c "DROP DATABASE cicero_db_dev;"'

db: init_db into_db create_db

run_db:
	docker start cicero_db_dev

migrate_up:
	migrate -database '$(DB_URL)' -source $(PATH_MIGRATE) -verbose up

migrate_down:
	migrate -database '$(DB_URL)' -source $(PATH_MIGRATE) -verbose down

build: 
	docker build -t asia.gcr.io/$(PROJECT_ID)/$(IMAGE_NAME) .

push:
	docker push asia.gcr.io/$(PROJECT_ID)/$(IMAGE_NAME)

.PHONY: init_db into_db create_db drop_db db run_db migrate_up migrate_down build push
