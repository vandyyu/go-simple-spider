package config
import (
	"log"
)

type HttpConfig struct{
	name string
}
func newHttpConfig(name string) *HttpConfig{
	c := new(HttpConfig)
	c.name = name
	return c
}
func (this *HttpConfig) GetName() string{
	return this.name
}
func (this *HttpConfig) Get(key string) *ConfigValue{
	switch key{
	case "":
	//	return NewConfigValue(this.name, this.LogRootDir)
	default:
		log.Fatalf("no value in LogConfig for [%s].\n", key)
	}
	return nil
}
