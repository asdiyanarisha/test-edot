package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"test-edot/util"
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
		Host: util.GetEnv("MYSQL_HOST", "localhost"),
		User: util.GetEnv("MYSQL_USER", "root"),
		Pass: util.GetEnv("MYSQL_PWD", ""),
		Port: util.GetEnv("MYSQL_PORT", "3306"),
		Name: util.GetEnv("MYSQL_DB", "test-edot"),
	}

	once.Do(func() {
		ConnMYSQL(conf)
	})
}

func GetConn() *gorm.DB {
	return dbConn
}

func ConnMYSQL(conf dbConfig) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		conf.User,
		conf.Pass,
		conf.Host,
		conf.Port,
		conf.Name,
		"utf8mb4",
		"True",
		"Local",
	)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DriverName:           "mysql",
		DisableWithReturning: true,
		DSN:                  dsn,
	}), &gorm.Config{
		SkipDefaultTransaction:   true,
		DisableNestedTransaction: true,
	})
	if err != nil {
		panic(err)
	}

	dbConn = db
	fmt.Println("success connect to MYSQL")
}
