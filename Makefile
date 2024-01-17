dev:
	air -c .air.dev.toml

prod:
	air -c .air.prod.toml

db:
    docker run --name cicero_db_dev -e POSTGRES_USER=cicero -e POSTGRES_PASSWORD=123456 -p 4444:5432 -d postgres:alpine

into_db:
	docker exec -it cicero_db_dev bash -c 'psql -U cicero'

create_db:
	CREATE DATABASE cicero_db_dev;

drop_db:
	DROP DATABASE cicero_db_dev;

migrate_up:
	migrate -database 'postgres://cicero:123456@localhost:4444/cicero_db_dev?sslmode=disable' -source file://C:/go/Cicero-Backend/pkg/databases/migrations -verbose up

migrate_down:
	migrate -database 'postgres://cicero:123456@localhost:4444/cicero_db_dev?sslmode=disable' -source file://C:/go/Cicero-Backend/pkg/databases/migrations -verbose down
