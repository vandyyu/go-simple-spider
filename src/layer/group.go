package layer

type IGroup interface{
	GetName() string
	SetPipeline(pipeline IPipeline)
	AddUnit(u IUnit)
	InitGroup()
	Run()
	SetLayer(layer ILayer)
	GetLayer() ILayer
}

type Group struct{
	name string
	units []IUnit
	pipeline IPipeline
	layer ILayer
}
func NewGroup(name string) *Group{
	g := new(Group)
	g.name = name
	g.units = make([]IUnit, 0, 128)
	return g
}
func (this *Group) SetPipeline(pipeline IPipeline){
	this.pipeline = pipeline
}
func (this *Group) AddUnit(u IUnit){
	// TODO: exception to handle
	this.units = append(this.units, u)
}
func (this *Group) InitGroup(){
	for _, u := range(this.units){
		u.SetPipeline(this.pipeline)
		u.SetGroup(this)
		u.InitUnit()
	}
}
func (this *Group) Run(){
	for _, u := range(this.units){
		u.Run()
	}
}
func (this *Group) SetLayer(layer ILayer){
	this.layer = layer
}
func (this *Group) GetLayer() ILayer{
	return this.layer
}
func (this *Group) GetName() string{
	return this.name
}
