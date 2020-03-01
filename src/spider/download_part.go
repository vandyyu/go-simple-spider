package spider
import (
	"config"
	"slog"
	"fmt"
	"layer"
)

type DownloadPart struct{
	link *config.Link
	name string
	rawData *config.RawData
	pipeline layer.IPipeline
	unit layer.IUnit
	downloader config.IDownloader
}
func (this *DownloadPart) InitObject(name string, link *config.Link){
	this.name = name
	this.link = link
}
/*
func NewDownloadPart(name string, link *config.Link) layer.IPart{
	d := new(DownloadPart)
	d.InitObject(name, link)
	return d
}
*/
func (this *DownloadPart) GetName() string{
	return this.name
}
func (this *DownloadPart) InitPart(){
}
func (this *DownloadPart) SetDownloader(downloader config.IDownloader){
	this.downloader = downloader
}
func (this *DownloadPart) Run() (layer.IPart, interface{}){
	if this.link == nil{
		msg := fmt.Sprintf("link is nil, it's not reasonable.")
		slog.Error(this.GetLogName(), msg)
		return nil, nil
	}
	if !this.link.Available{
		msg := fmt.Sprintf("link %s is not an available url during Downloading.", this.link.UnavailableAddr)
		slog.Error(this.GetLogName(), msg)
		return nil, nil
	}
	var err error
	this.rawData, err = this.downloader.GenerateRawData()
	if err != nil{
		msg := fmt.Sprintf("Dowlading url %s failed.", this.link.URL.String())
		slog.Error(this.GetLogName(), msg, err)
		return nil, nil
	}
	if this.rawData == nil{
		slog.Warning(this.GetLogName(), "The rawData is nil after downloading.Maybe forget to setup value of this.rawData in GenerateRawData() function.")
		return nil, this.rawData
	}
	depth := this.unit.GetGroup().GetLayer().GetDepth()
	NewRP := config.Get(this.GetLogName(), config.PART_CONF, "NewResolverPartMap").RPPtr(depth)
	if NewRP == nil{
		slog.Error(this.GetLogName(), "Failed to initialize ResolverPartMap.")
		return nil, nil
	}
	rp := NewRP("ResolverPart", this.rawData)
	rp.SetResolver(rp)
	return rp, this.rawData
}
func (this *DownloadPart) SetPipeline(pipeline layer.IPipeline){
	this.pipeline = pipeline
}
func (this *DownloadPart) SetUnit(unit layer.IUnit){
	this.unit = unit
}
func (this *DownloadPart) GetUnit() layer.IUnit{
	return this.unit
}
func (this *DownloadPart) GetPipeline() layer.IPipeline{
	return this.pipeline
}
func (this *DownloadPart) SetRawData(rawData *config.RawData){
	this.rawData = rawData
}
func (this *DownloadPart) GetLink() *config.Link{
	return this.link
}
func (this *DownloadPart) GetLogName() string{
	return this.pipeline.GetName()
}
