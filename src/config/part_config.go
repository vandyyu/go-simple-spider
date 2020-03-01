package config
import (
	"layer"
	"log"
	"runtime/debug"
)

type DownloadPartPtr func(string, *Link) IDownloaderPart
type ResolverPartPtr func(string, *RawData) IResolverPart

type IDownloader interface{
	GenerateRawData() (*RawData, error)
}
type IResolver interface{
	GenerateLinks() ([]*Link, error)
	GenerateText() (*Text, error)
}
type IDownloaderPart interface{
	layer.IPart
	IDownloader
	SetDownloader(downloader IDownloader)
}
type IResolverPart interface{
	layer.IPart
	IResolver
	SetResolver(resolver IResolver)
}

type PartConfig struct{
	name string
	NewDownloadPartMap map[int]DownloadPartPtr
	NewResolverPartMap map[int]ResolverPartPtr
}
func newPartConfig(name string) *PartConfig{
	c := new(PartConfig)
	c.name = name
	return c
}
func (this *PartConfig) GetName() string{
	return this.name
}
func (this *PartConfig) Get(key string) *ConfigValue{
	switch key{
	case "NewDownloadPartMap":
		return NewConfigValue(this.name, this.NewDownloadPartMap)
	case "NewResolverPartMap":
		return NewConfigValue(this.name, this.NewResolverPartMap)
	default:
		debug.PrintStack()
		log.Fatalf("no value in PartConfig for [%s].\n", key)
	}
	return nil
}
func (this *PartConfig) SetDPPtrMap(dpptrMap map[int]DownloadPartPtr){
	this.NewDownloadPartMap = dpptrMap
}
func (this *PartConfig) SetRPPtrMap(rpptrMap map[int]ResolverPartPtr){

	this.NewResolverPartMap = rpptrMap
}
