package spider

import (
	"config"
)
func DumpData(texts []*config.Text){
	for _, d := range(texts){
		println(d.Data)
	}
}
