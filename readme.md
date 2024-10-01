Test Senior Backend Engineer E-dot
-------------------------

Requirements:
- Golang 1.19
- MySQL
- Migrate
- Docker
- Makefile

### Services in this app:
- rest api
- scheduler worker (for release stock)

### How To Run

#### run script
```shell
go get
go run main.go
```

#### run unittest
```shell
go test ./... -v
```

### Run Integration Test
first set config env `.env` to set configuration for testing, and next run initial dependency integration db such as database:
```shell
make initial-setup
```

this command for run integration test
```shell
make integration-test
```


### How Install golang migrate
To install the migrate CLI tool using curl on Linux, you can follow these steps:
```shell
$ curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey| apt-key add -
$ echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
$ apt-get update
$ apt-get install -y migrate
## install dependency client library of database
$ go install -tags mysql github.com/golang-migrate/migrate/v4/cmd/migrate
```

Next, create migration files using the following command:

```shell
$ migrate create -ext sql -dir database/migration/ -seq create_table_card
```

#### Run Migration Up

```shell
$ migrate -path database/migration/ -database "mysql://root:rahasia123@tcp(localhost:3306)/komcard_dev?multiStatements=true" -verbose up
```

#### Run Migration Down

```shell
$ migrate -path database/migration/ -database "mysql://root:rahasia123@tcp(localhost:3306)/komcard_dev?multiStatements=true" -verbose down
```
