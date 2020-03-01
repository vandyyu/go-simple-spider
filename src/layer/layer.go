package layer
import (
	"sync"
)

type ILayer interface{
	InitLayer()
	AddGroup(g IGroup)
	GetName() string
	Run()
	GetDepth() int
	SetPipeline(c IPipeline)
}

type Layer struct{
	grps []IGroup
	name string
	depth int
	wg sync.WaitGroup
	pipeline IPipeline
}
func (this *Layer) SetPipeline(pipeline IPipeline){
	this.pipeline = pipeline
}
func (this *Layer) InitLayer(){
	for _, g := range(this.grps){
		g.SetPipeline(this.pipeline)
		g.SetLayer(this)
		g.InitGroup()
	}
}
func (this *Layer) GetDepth() int{
	return this.depth
}
func NewLayer(name string, depth int) *Layer{
	l := new(Layer)
	l.name = name
	l.depth = depth
	l.grps = make([]IGroup, 0, 128)
	return l
}
func (this *Layer) AddGroup(g IGroup) {
	this.grps = append(this.grps, g)
}
func (this *Layer) GetName() string{
	return this.name
}
func (this *Layer) Run(){
	this.wg.Add(len(this.grps))
	for _, g := range(this.grps){
		go func(p_g IGroup){
			p_g.Run()
			this.wg.Done()
		}(g)
	}
	this.wg.Wait()
}
