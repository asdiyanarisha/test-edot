#ENV_LOCAL_TEST=\
#  MYSQL_DB=test-edot-new\
#  MYSQL_USER=edot\
#  MYSQL_PWD=edot123

include .env

initial-setup:
	#$(ENV_LOCAL_TEST) \
	docker compose up --build --remove-orphans

test-up:
	migrate -path database/migration/ -database "mysql://$(MYSQL_USER):$(MYSQL_PWD)@tcp(localhost:3306)/$(MYSQL_DB)?multiStatements=true" -verbose up

test-down:
	migrate -path database/migration/ -database "mysql://$(MYSQL_USER):$(MYSQL_PWD)@tcp(localhost:3306)/$(MYSQL_DB)?multiStatements=true" -verbose down -all
	docker compose down --volumes

migrate-up:
	migrate -path database/migration/ -database "mysql://$(MYSQL_USER):$(MYSQL_PWD)@tcp(localhost:3306)/$(MYSQL_DB)?multiStatements=true" -verbose up

migrate-down:
	migrate -path database/migration/ -database "mysql://$(MYSQL_USER):$(MYSQL_PWD)@tcp(localhost:3306)/$(MYSQL_DB)?multiStatements=true" -verbose down -all