package spider
import (
	"layer"
)

type PartProxy struct{
	name string
	part layer.IPart
	pipeline layer.IPipeline
	unit layer.IUnit
}
func NewPartProxy(part layer.IPart) *PartProxy{
	proxy := new(PartProxy)
	proxy.part = part
	return proxy
}
func (this *PartProxy) Run() (layer.IPart, interface{}){
	return this.part.Run()
}
func (this *PartProxy) SetPipeline(pipeline layer.IPipeline){
	this.pipeline = pipeline
	this.part.SetPipeline(pipeline)
}
func (this *PartProxy) SetUnit(unit layer.IUnit){
	this.unit = unit
	this.part.SetUnit(unit)
}
func (this *PartProxy) GetUnit() layer.IUnit{
	return this.unit
}
func (this *PartProxy) GetName() string{
	return this.part.GetName()
}
func (this *PartProxy) InitPartProxy(){
	// init part proxy
}
func (this *PartProxy) InitPart(){
	this.part.InitPart()
}
func (this *PartProxy) Forward() (layer.IPartProxy, interface{}){
	this.InitPart()
	nextPart, data := this.Run()
	if nextPart != nil{
		return NewPartProxy(nextPart), data
	}else{
		return nil, data
	}
}
