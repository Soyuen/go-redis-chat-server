package redis

type Config interface {
	Host() string
	Port() int
	Password() string
	DB() int
}
