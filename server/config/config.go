package config

import "os"

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Download DownloadConfig
	JMComic  JMComicConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Path string
}

type DownloadConfig struct {
	Dir string
}

type JMComicConfig struct {
	BaseURL string
	Auth    string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "5000"),
		},
		Database: DatabaseConfig{
			Path: getEnv("DB_PATH", "./data/jmcomic.db"),
		},
		Download: DownloadConfig{
			Dir: getEnv("DOWNLOAD_DIR", "./data/downloads"),
		},
		JMComic: JMComicConfig{
			BaseURL: getEnv("JM_BASE_URL", "https://www.cdnhjk.net"),
			Auth:    getEnv("JM_AUTH", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
