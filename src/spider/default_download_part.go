package spider
import (
	"config"
	"snet"
	"fmt"
	"errors"
)
type DefaultDownloadPart struct{
	DownloadPart
}
func NewDefaultDownloadPart(name string, link *config.Link) config.IDownloaderPart{
	d := new(DefaultDownloadPart)
	d.InitObject(name, link)
	return d
}
func (this *DefaultDownloadPart) GenerateRawData() (*config.RawData, error){
	data := snet.RequestNative(this.pipeline.GetName(), this.link)
	if data == ""{
		return nil, errors.New(fmt.Sprintf("Failed to download data for url %s.", this.link.URL.String()))
	}
	rawData := config.NewRawData(this.link, data)
	return rawData, nil
}
