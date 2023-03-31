package config

type Config struct {
	HttpPort    string
	PostgresURL string
}

func NewConfig() Config {
	return Config{
		HttpPort:    "8080",
		PostgresURL: "postgresql://myuser:password@localhost:5432/postgresDB",
	}
}
