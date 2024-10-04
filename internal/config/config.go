package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost string
	Storage string
	DeleteAfter int
	EnableHttps bool
}

func Init() (*Config, error) {
	err := godotenv.Load("./.env")
	if err != nil {
		return &Config{}, err
	}

	serverHost := os.Getenv("SERVER_HOST")
	storage := os.Getenv("STORAGE")
	deleteAfter, err :=  strconv.Atoi(os.Getenv("DELETE_AFTER"))
	if err != nil {
		deleteAfter = 3
	}

	enableHttps, err := strconv.ParseBool(os.Getenv("ENABLE_HTTPS"))
	if err != nil {
		enableHttps = false
	}

	return &Config{
		ServerHost: serverHost,
		Storage: storage,
		DeleteAfter: deleteAfter,
		EnableHttps: enableHttps,
	}, nil
}