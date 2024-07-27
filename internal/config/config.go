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

	return &Config{
		ServerHost: serverHost,
		Storage: storage,
		DeleteAfter: deleteAfter,
	}, nil
}