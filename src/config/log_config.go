package config
import (
	"path"
	"os"
	"log"
)

type LogConfig struct{
	name string
	LogRootDir string `json:"logRootDir, string"`
	LogFileSize int64 `json:"logFileSize, int64"`
}
func newLogConfig(name string) *LogConfig{
	c := new(LogConfig)
	c.name = name
	c.LogRootDir = path.Join(os.Getenv("GOPATH"), "logs")

	// max size of single log file is about 15M.
	c.LogFileSize = 1024 * 1024 * 15
	loadConfFile(c, name, "log.conf")
	return c
}
func (this *LogConfig) GetName() string{
	return this.name
}
func (this *LogConfig) Get(key string) *ConfigValue{
	switch key{
	case "LogRootDir":
		return NewConfigValue(this.name, this.LogRootDir)
	case "LogFileSize":
		return NewConfigValue(this.name, this.LogFileSize)
	default:
		log.Fatalf("no value in LogConfig for [%s].\n", key)
	}
	return nil
}
