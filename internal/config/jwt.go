package config

type JWTConfig struct {
	Secret     string
	Expiration int // часов
}
