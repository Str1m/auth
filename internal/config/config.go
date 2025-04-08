package config

import "github.com/joho/godotenv"

func MustLoad(path string) {
	if err := godotenv.Load(path); err != nil {
		panic("failed to load config")
	}
}
