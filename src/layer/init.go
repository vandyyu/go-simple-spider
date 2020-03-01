package layer
import (
	"fmt"
)

/*
numNode: param of how many units in this layer.
numUnit: param of how many units make up a group in this layer.
*/
func InitLayer(pipeline IPipeline, name string, numNode int, numUnit int, curDepth int) ILayer{
	var g IGroup
	var k = 0
	var m = 0
	curLayer := NewLayer(name, curDepth)
	for i:=0;i < numNode;i++{
		if i % numUnit == 0{
			k = 0
			g = NewGroup(fmt.Sprintf("Group-%d", m))
			curLayer.AddGroup(g)
			m += 1
		}
		u := NewUnit(fmt.Sprintf("Unit-%d", k))
		u.SetFirstPartProxy(pipeline.GetFirstPartProxy())
		g.AddUnit(u)
		k += 1
	}
	curLayer.SetPipeline(pipeline)
	curLayer.InitLayer()
	return curLayer
}
