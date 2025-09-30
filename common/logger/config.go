package logger

type Config struct {
	Level         string `yaml:"level"`         // 日志级别: debug, info, warn, error
	OutputPath    string `yaml:"outputPath"`    // 日志输出路径, 为空则不写文件
	MaxSize       int    `yaml:"maxSize"`       // 每个日志文件的最大大小（MB）
	MaxBackups    int    `yaml:"maxBackups"`    // 保留的旧日志文件的最大个数
	MaxAge        int    `yaml:"maxAge"`        // 保留旧日志文件的最大天数
	Compress      bool   `yaml:"compress"`      // 是否压缩旧日志文件
	EnableConsole bool   `yaml:"enableConsole"` // 是否同时输出到控制台
}

func DefaultConfig() Config {
	return Config{
		Level:         "info",
		OutputPath:    "logs/server.log",
		MaxSize:       100,
		MaxBackups:    7,
		MaxAge:        30,
		Compress:      true,
		EnableConsole: true,
	}
}
