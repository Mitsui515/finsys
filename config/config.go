package config

type AppConfig struct {
	Server ServerConfig `json:"server"`
	Log    LogConfig    `json:"log"`
	LLM    LLMConfig    `json:"llm"`
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

type LLMConfig struct {
	APIKey string `json:"api_key"`
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
		LLM: LLMConfig{
			APIKey: "sk-0709296703ab41f280007b5f45b8f9a0",
		},
	}
}
