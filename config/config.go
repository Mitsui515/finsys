package config

type AppConfig struct {
	Server ServerConfig `json:"server"`
	Log    LogConfig    `json:"log"`
}

type ServerConfig struct {
	Port    int    `json:"port"`
	Host    string `json:"host"`
	Mode    string `json:"mode"`
	Version string `json:"version"`
}

type LogConfig struct {
	Level      string `json:"level"`
	FilePath   string `json:"file_path"`
	MaxSize    int    `json:"max_size"`
	MaxBackups int    `json:"max_backups"`
	MaxAge     int    `json:"max_age"`
	Compress   bool   `json:"compress"`
}

func DefaultConfig() AppConfig {
	return AppConfig{
		Server: ServerConfig{
			Port:    8080,
			Host:    "0.0.0.0",
			Mode:    "debug",
			Version: "v1",
		},
		Log: LogConfig{
			Level:      "info",
			FilePath:   "./logs/app.log",
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     30,
			Compress:   true,
		},
	}
}
