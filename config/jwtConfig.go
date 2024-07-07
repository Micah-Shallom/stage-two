package config

import "os"

type JwtConfig struct {
	Secret []byte

}

func NewJwtConfig() *JwtConfig {
	return &JwtConfig{
		Secret: []byte(os.Getenv("JWT_SECRET")),
	}

}