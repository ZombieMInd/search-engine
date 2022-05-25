package server

type Config struct {
	Name      string `envconfig:"NAME" default:"logger"`
	BindAddr  string `envconfig:"BIND_ADDR" default:"127.0.0.1:8080"`
	DBURL     string `envconfig:"DB_URL" default:"host=localhost dbname=restapi_dev sslmode=disable"`
	StoreMode string `envconfig:"STORE_MODE" default:"psql"`
	Redis     RedisConfig
}

type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST" default:"localhost"`
	Password string `envconfig:"REDIS_PASSWORD" default:""`
	DB       int    `envconfig:"REDIS_DB"`
}
