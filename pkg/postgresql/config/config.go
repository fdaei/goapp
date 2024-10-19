package postgresql

type Config struct {
	Host            string `koanf:"host"`
	Port            int    `koanf:"port"`
	User            string `koanf:"user"`
	Password        string `koanf:"password"`
	DBName          string `koanf:"dbName"`
	SSLMode         string `koanf:"sslMode"`
	MaxIdleConns    int    `koanf:"maxIdleConns"`
	MaxOpenConns    int    `koanf:"maxOpenConns"`
	ConnMaxLifetime int    `koanf:"connMaxLifetime"`
}
