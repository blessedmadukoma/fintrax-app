include .env

DB_DRIVER := $(DB_DRIVER)
DB_USER := $(DB_USER)
DB_PASSWORD := $(DB_PASSWORD)
DB_HOST := $(DB_HOST)
DB_PORT := $(DB_PORT)
DB_NAME := $(DB_NAME)

db/createmigration:
	migrate create -ext sql -dir db/migrations -seq $(name)

db/migrateup:
	# migrate up "$(DB_URL)"
	migrate -path db/migrations -database "$(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up

db/migratedown:
	# migrate down
	migrate -path db/migrations -database "$(DB_DRIVER)://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down

pg_up:
# postgres up - create a postgres container
	docker-compose up -d

pg_down:
# postgres down - delete a postgres container
	docker-compose down

db_up:
	# create database
	docker exec -it fintrax_db createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

db_down:
	# drop database
	docker exec -it fintrax_db dropdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

sqlc:
	sqlc generate

.PHONY:	db/createmigration