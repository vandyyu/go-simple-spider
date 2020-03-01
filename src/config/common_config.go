package config
import (
	"log"
)


type CommonConfig struct{
	MaxDepth int `json:"maxDepth, int"`

	// a Group composite of {NumUnit} Units.
	NumUnit int `json:"numUnit, int"`

	// storage data in cachedText when len(cachedText) is equal to {NumCachedText}
	NumCachedText int `json:"numCachedText, int"`

	name string
}
func newCommonConfig(name string) *CommonConfig{
	c := new(CommonConfig)
	c.name = name
	c.MaxDepth = 0
	c.NumUnit = 2
	c.NumCachedText = 5
	loadConfFile(c, name, "common.conf")
	return c
}
func (this *CommonConfig) GetName() string{
	return this.name
}
// maybe should use reflect to make a common Get() method, but not good for efficience.
func (this *CommonConfig) Get(key string) *ConfigValue{
	switch key{
	case "MaxDepth":
		return NewConfigValue(this.name, this.MaxDepth)
	case "NumUnit":
		return NewConfigValue(this.name, this.NumUnit)
	case "NumCachedText":
		return NewConfigValue(this.name, this.NumCachedText)
	default:
		log.Fatalf("no value in CommonConfig for [%s].\n", key)
	}
	return nil
}
