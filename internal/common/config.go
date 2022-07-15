package common

type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST" default:"localhost:6379"`
	Password string `envconfig:"REDIS_PASSWORD" default:""`
	DB       int    `envconfig:"REDIS_DB"`
}

type DetaBaseConfig struct {
	Key string `envconfig:"DETA_BASE_KEY" required:"true"`
}

type DBConfig struct {
	DBURL     string `envconfig:"DB_URL" default:"host=localhost dbname=restapi_dev sslmode=disable"`
	StoreMode string `envconfig:"STORE_MODE" default:"psql"`
}
