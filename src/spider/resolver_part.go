package spider
import (
	"layer"
	"slog"
	"config"
	"fmt"
)
type ResolverPart struct{
	name string
	rawData *config.RawData      // source
	text *config.Text            // one choice of destination
	links []*config.Link           // another choice of dest
	pipeline layer.IPipeline
	unit layer.IUnit
	resolver config.IResolver
}
func (this *ResolverPart) InitObject(name string, rawData *config.RawData){
	this.name = name
	this.rawData = rawData
}
/*
func NewResolverPart(name string, rawData *config.RawData) layer.IPart{
	r := new(ResolverPart)
	r.InitObject(name, rawData)
	return r
}
*/
func (this *ResolverPart) SetPipeline(pipeline layer.IPipeline){
	this.pipeline = pipeline
}
func (this *ResolverPart) SetUnit(unit layer.IUnit){
	this.unit = unit
}
func (this *ResolverPart) GetUnit() layer.IUnit{
	return this.unit
}
func (this *ResolverPart) SetResolver(resolver config.IResolver){
	this.resolver = resolver
}

func (this *ResolverPart) GetName() string{
	return this.name
}
func (this *ResolverPart) InitPart(){
	// TODO: init resolver part
}
func (this *ResolverPart) Run() (layer.IPart, interface{}){
	if this.rawData == nil{
		slog.Error(this.GetLogName(), "rawData is nil, it's not reasonable.")
		return nil, nil
	}
	if !this.rawData.Available{
		if this.rawData.LINK != nil{
			msg := fmt.Sprintf("The this.rawData of %s is not available during Resolving links.", this.rawData.LINK.UnavailableAddr)
			slog.Warning(this.GetLogName(), msg)
			return nil, nil
		}else{
			msg := "The this.rawData is Unavailable and LINK is nil, nothing to Resolve."
			slog.Warning(this.GetLogName(), msg)
			return nil, nil
		}

	}
	if this.rawData.Data == ""{
		msg := fmt.Sprintf("No data after downloading url %s. Cannot resolve to generate links.", this.rawData.LINK.URL.String())
		slog.Warning(this.GetLogName(), msg)
		return nil, nil
	}
	var err error
	this.links, err = this.resolver.GenerateLinks()
	if err != nil{
		msg := fmt.Sprintf("Resolving url %s failed!", this.rawData.LINK.URL.String())
		slog.Error(this.GetLogName(), msg, err)
		return nil, nil
	}
	if this.links == nil{
		msg := fmt.Sprintf("The this.links is nil after Resolving. Maybe forget to setup this.links value in GenerateLinks() function.")
		slog.Warning(this.GetLogName(), msg)
		return nil, nil
	}
	depth := this.unit.GetGroup().GetLayer().GetDepth()
	maxDepth := config.Get(this.GetLogName(), config.COMMON_CONF, "MaxDepth").Int()
	flag := depth < maxDepth - 1 && len(this.links) != 0
	if flag{
		return nil, this.links
	}
	this.text, err = this.resolver.GenerateText()
	if err != nil{
		msg := fmt.Sprintf("Generating Text of %s failed.", this.rawData.LINK.URL.String())
		slog.Error(this.GetLogName(), msg, err)
		return nil, nil
	}
	if this.text == nil{
		msg := fmt.Sprintf("The this.text is nil after Resolving. Maybe forget to setup this.text value in GenerateText() function.")
		slog.Warning(this.GetLogName(), msg)
		return nil, nil
	}
	if !this.text.Available{
		msg := fmt.Sprintf("The this.text of %s is not available during Resolving text.", this.rawData.LINK.UnavailableAddr)
		slog.Warning(this.GetLogName(), msg)
		return nil, nil
	}
	return nil, this.text
}
func (this *ResolverPart) GetPipeline() layer.IPipeline{
	return this.pipeline
}
func (this *ResolverPart) SetText(text *config.Text){
	this.text = text
}
func (this *ResolverPart) Setlinks(links []*config.Link){
	this.links = links
}
func (this *ResolverPart) GetRawData() *config.RawData{
	return this.rawData
}
func (this *ResolverPart) GetLogName() string{
	return this.pipeline.GetName()
}
