package config

type AppConfig struct {
	Port string
}

func LoadEnvConfig() AppConfig {
	port := GetEnv("PORT", "8000")

	return AppConfig{
		Port: port,
	}
}
