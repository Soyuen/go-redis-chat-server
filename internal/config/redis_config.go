package config

type RedisConfig struct {
	HostVal     string
	PortVal     int
	PasswordVal string
	DBVal       int
}

func (c RedisConfig) Host() string     { return c.HostVal }
func (c RedisConfig) Port() int        { return c.PortVal }
func (c RedisConfig) Password() string { return c.PasswordVal }
func (c RedisConfig) DB() int          { return c.DBVal }

func LoadRedisConfigFromEnv() RedisConfig {
	return RedisConfig{
		HostVal:     GetEnv("REDIS_HOST", "localhost"),
		PortVal:     GetEnvAsInt("REDIS_PORT", 6379),
		PasswordVal: GetEnv("REDIS_PASSWORD", ""),
		DBVal:       GetEnvAsInt("REDIS_DB", 0),
	}
}
