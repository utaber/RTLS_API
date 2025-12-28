package config

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret string
	JWTAlgo   jwt.SigningMethod
	JWTExpire time.Duration
	DBURL     string
}

func Load() *Config {
	_ = godotenv.Load()

	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		log.Fatal("JWT_SECRET_KEY is not set")
	}

	return &Config{
		JWTSecret: secret,
		JWTAlgo:   jwt.SigningMethodHS256,
		JWTExpire: time.Hour,
		DBURL:     "https://simple-app-3ed62-default-rtdb.asia-southeast1.firebasedatabase.app/",
	}
}
