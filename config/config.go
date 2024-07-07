package config

// responsible for getting the following config
//   Database
//   JWT


type Config struct {
	Database *DatabaseConfig
	Jwt *JwtConfig
}

func NewConfig() *Config {
	return &Config{
		Database: NewDatabase(),
		Jwt: NewJwtConfig(),
	}
}


