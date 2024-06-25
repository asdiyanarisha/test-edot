package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
	"test-asset-fendr/util"
)

var (
	dbConn *gorm.DB
	once   sync.Once
)

type (
	dbConfig struct {
		Host string
		User string
		Pass string
		Port string
		Name string
	}
)

func CreateConn() {
	conf := dbConfig{
		Host: util.GetEnv("PSQL_HOST", "localhost"),
		User: util.GetEnv("PSQL_USER", "postgres"),
		Pass: util.GetEnv("PSQL_PWD", "postgres"),
		Port: util.GetEnv("PSQL_PORT", "5432"),
		Name: util.GetEnv("PSQL_DB", "postgres"),
	}

	once.Do(func() {
		ConnPsql(conf)
	})
}

func GetConn() *gorm.DB {
	return dbConn
}

func ConnPsql(conf dbConfig) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		conf.Host, conf.User, conf.Pass, conf.Name, conf.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("cant connect postgresl")
		return
	}

	dbConn = db
	fmt.Println("success connect to psql")
}
