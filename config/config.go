package config

type Config struct {
	Port      string `mapstructure:"port"`
	Copyright string `mapstructure:"copyright"`
	Server    struct {
		Port         string `mapstructure:"port"`
		Host         string `mapstructure:"host"`
		ReadTimeout  string `mapstructure:"read_timeout"`
		WriteTimeout string `mapstructure:"write_timeout"`
	} `mapstructure:"server"`
	Database struct {
		FilePath       string `mapstructure:"file_path"`
		BackupEnabled  bool   `mapstructure:"backup_enabled"`
		BackupInterval string `mapstructure:"backup_interval"`
	} `mapstructure:"database"`
	Cache struct {
		RedisURL       string `mapstructure:"redis_url"`
		DefaultTTL     string `mapstructure:"default_ttl"`
		MaxConnections int    `mapstructure:"max_connections"`
	} `mapstructure:"cache"`
	Logging struct {
		Level  string `mapstructure:"level"`
		Format string `mapstructure:"format"`
		Output string `mapstructure:"output"`
	} `mapstructure:"logging"`
	Admin struct {
		Enabled      bool   `mapstructure:"enabled"`
		Username     string `mapstructure:"username"`
		PasswordHash string `mapstructure:"password_hash"`
	} `mapstructure:"admin"`
}
