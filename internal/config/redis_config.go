package config

type RedisConfig interface {
	Host() string
	Port() int
	Password() string
	DB() int
}

type RedisConfigImpl struct {
	HostVal     string
	PortVal     int
	PasswordVal string
	DBVal       int
}

func (c RedisConfigImpl) Host() string     { return c.HostVal }
func (c RedisConfigImpl) Port() int        { return c.PortVal }
func (c RedisConfigImpl) Password() string { return c.PasswordVal }
func (c RedisConfigImpl) DB() int          { return c.DBVal }

func LoadRedisConfigFromEnv() RedisConfig {
	return RedisConfigImpl{
		HostVal:     GetEnv("REDIS_HOST", "localhost"),
		PortVal:     GetEnvAsInt("REDIS_PORT", 6379),
		PasswordVal: GetEnv("REDIS_PASSWORD", ""),
		DBVal:       GetEnvAsInt("REDIS_DB", 0),
	}
}
