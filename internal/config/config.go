package config

import "github.com/joho/godotenv"

func MustLoad(path string) error {
	if err := godotenv.Load(path); err != nil {
		return err
	}
	return nil
}
