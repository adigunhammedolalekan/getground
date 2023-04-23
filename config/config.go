package config

import "os"

type Config struct {
	DatabaseUser, DatabasePassword, DatabaseName, DatabaseAddr string
	Port                                                       string
}

func New() Config {
	return Config{
		DatabaseUser:     os.Getenv("DATABASE_USER"),
		DatabasePassword: os.Getenv("DATABASE_PASSWORD"),
		DatabaseName:     os.Getenv("DATABASE_NAME"),
		DatabaseAddr:     os.Getenv("DATABASE_ADDR"),
		Port:             os.Getenv("PORT"),
	}
}
