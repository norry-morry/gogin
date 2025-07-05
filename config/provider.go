// Package config provides MySQL connection configuration settings.
package config

// ProvideMySQLSettings は指定された Config から MySQLSettings を返します。
func ProvideMySQLSettings(cfg Config) MySQLSettings {
	return cfg.MySQL
}
