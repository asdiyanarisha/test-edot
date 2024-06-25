package util

import (
	"github.com/joho/godotenv"
	"os"
)

func GetEnv(key string, fallback string) string {
	// Godotenv read the .env file on the root folder
	a, _ := godotenv.Read()
	var (
		val     string
		isExist bool
	)
	// Check the key of the env using Hashmap
	// if exist return the actual value, if !exist return the fallback value
	val, isExist = a[key]
	if !isExist {
		val = os.Getenv(key)
		if val == "" {
			val = fallback
		}
	}
	return val
}
