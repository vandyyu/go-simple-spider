package layer
import (
	"log"
)
type IUnit interface{
	GetName() string
	SetPipeline(pipeline IPipeline)
	SetFirstPartProxy(p IPartProxy)
	Run()
	InitUnit()
	SetGroup(g IGroup)
	GetGroup() IGroup

	// back to previous PartProxy, return int data shows how much ProxyPart is left after backwarding.
	// note: Backward() will delete current PartProxy and back to previous PartProxy
	Backward() int
}

type Unit struct{
	name string
	firstPP IPartProxy
	pipeline IPipeline
	grp IGroup
	pps []IPartProxy
	i_pps int
	backward bool
}
func NewUnit(name string) IUnit{
	u := new(Unit)
	u.name = name
	u.pps = make([]IPartProxy, 0, 128)
	return u
}
func (this *Unit) SetPipeline(pipeline IPipeline){
	this.pipeline = pipeline
}
func (this *Unit) GetName() string{
	return this.name
}
func (this *Unit) SetGroup(g IGroup){
	this.grp = g
}
func (this *Unit) GetGroup() IGroup{
	return this.grp
}
func (this *Unit) InitUnit(){
	if this.firstPP == nil{
		log.Fatalln("layer/unit.go: InitUnit() firstPartProxy shoudn't be nil.")
	}
	this.init(this.firstPP)
}
func (this *Unit) init(pp IPartProxy){
	pp.SetPipeline(this.pipeline)
	pp.SetUnit(this)
	pp.InitPartProxy()
	this.pps = append(this.pps, pp)
}
func (this *Unit) SetFirstPartProxy(p IPartProxy){
	this.firstPP = p
}
func (this *Unit) Backward() int{
	if len(this.pps) > 0{
		this.pps = this.pps[:(len(this.pps)-1)]
	}
	this.backward = true
	return len(this.pps)
}
// traversing all PartProxy in Units.
func (this *Unit) Run(){
	for {
		if len(this.pps) == 0{
			log.Println("layer/unit.go: Run() Skip unit runing, no parts in units [", this.name, "]")
			break
		}
		ret, data := this.pps[len(this.pps)-1].Forward()
		if this.backward {
			this.backward = false
			continue
		}
		this.pipeline.Write(data)   // save result of current part
		if ret == nil{         // no next part
			break
		}else{
			this.init(ret)
		}
	}
}

