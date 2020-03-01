package sina
import (
	"config"
	"snet"
	"spider"
)
/********************** layer-0 download strategy *************************/
type DP0 struct{
	spider.DownloadPart
}
func NewDP0(name string, link *config.Link) config.IDownloaderPart{
	d := new(DP0)
	d.InitObject(name, link)
	return d
}
func (this *DP0) GenerateRawData() (*config.RawData, error){
	data := snet.RequestNative(this.GetLogName(), this.GetLink())
	rawData := config.NewRawData(this.GetLink(), data)
	return rawData, nil
}

/********************** layer-1 download strategy *************************/
type DP1 struct{
	spider.DownloadPart
}
func NewDP1(name string, link *config.Link) config.IDownloaderPart{
	d := new(DP1)
	d.InitObject(name, link)
	return d
}
func (this *DP1) GenerateRawData() (*config.RawData, error){
	data := snet.RequestSelenium(this.GetLogName(), this.GetLink())
	rawData := config.NewRawData(this.GetLink(), data)
	return rawData, nil
}

/********************** layer-2 download strategy *************************/
type DP2 struct{
	spider.DownloadPart
}
func NewDP2(name string, link *config.Link) config.IDownloaderPart{
	d := new(DP2)
	d.InitObject(name, link)
	return d
}
func (this *DP2) GenerateRawData() (*config.RawData, error){
	data := snet.RequestNative(this.GetLogName(), this.GetLink())
	rawData := config.NewRawData(this.GetLink(), data)
	return rawData, nil
}
