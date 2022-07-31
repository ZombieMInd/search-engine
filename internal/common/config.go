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

type SearchDBConfig struct {
	Host         string `envconfig:"SEARCH_DB_URL" default:"http://127.0.0.1:7700"`
	APIKey       string `envconfig:"SEARCH_DB_API_KEY" default:"masterKey"`
	SitemapIndex string `envconfig:"SEARCH_DB_SITEMAP_INDEX" default:"sitemaps"`
}
