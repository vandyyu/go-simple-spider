package spider
import (
	"layer"
	"config"
	"slog"
	"os"
	"sync"
	"path"
	"log"
	"io"
	"fmt"
)
type IManagerRunner interface{
	Run()
	GetName() string
}
var runners = make([]IManagerRunner, 0, 32)
var wg sync.WaitGroup

func AppendManager(r IManagerRunner){
	runners = append(runners, r)
}
func RunAll(){
	wg.Add(len(runners))
	for _, r := range(runners){
		go func(r_p IManagerRunner){
			r_p.Run()
			wg.Done()
			log.Printf("[%s] spider has Done.\n", r_p.GetName())
		}(r)
	}
	wg.Wait()
}

type SpiderManager struct{
	pipeline *SpiderPipeline
	initLayer layer.ILayer
	dpptrMap map[int]config.DownloadPartPtr
	rpptrMap map[int]config.ResolverPartPtr
	seeds []*config.Link
	name string
	maxDepth int
	start bool;
}
func NewSpiderManager(name string, seeds ...*config.Link) *SpiderManager{
	sm := new(SpiderManager)
	sm.name = name
	sm.start = true;
	sm.seeds = seeds
	InitGlobalLog(name)

	sm.rpptrMap = make(map[int]config.ResolverPartPtr)
	sm.dpptrMap = make(map[int]config.DownloadPartPtr)

	config.InitConfig(name, sm.dpptrMap, sm.rpptrMap)
	InitSpiderLog(name)
	sm.pipeline = NewSpiderPipeline(name)
	sm.pipeline.AddSeeds(seeds...)

	sm.maxDepth = config.Get(name, config.COMMON_CONF, "MaxDepth").Int()
	sm.check(sm.maxDepth <= 0, fmt.Sprintf("[%s] maxDepth shoudn't be zero! please check your config file.\n", sm.name))
	return sm
}
func (this *SpiderManager) AddLayerStrategy(layer int, dpptr config.DownloadPartPtr, rpptr config.ResolverPartPtr){
	if _, ok := this.dpptrMap[layer];!ok{
		this.dpptrMap[layer] = dpptr
	}
	if _, ok := this.rpptrMap[layer];!ok{
		this.rpptrMap[layer] = rpptr
	}
}
func (this *SpiderManager) checkDownloadPartPtr(){
	for i:=0;i < this.maxDepth;i++{
		if v := config.Get(this.name, config.PART_CONF, "NewDownloadPartMap").DPPtr(i);v == nil{
			log.Printf("[%s] use DefaultDownloadPart in layer-%d\n", this.name, i)
			this.dpptrMap[i] = NewDefaultDownloadPart
		}
	}
}
func (this *SpiderManager) checkResolverPartPtr(){
	for i:=0;i < this.maxDepth;i++{
		if v := config.Get(this.name, config.PART_CONF, "NewResolverPartMap").RPPtr(i);v == nil{
			log.Printf("[%s] use DefaultResolverPart in layer-%d\n", this.name, i)
			this.rpptrMap[i] = NewDefaultResolverPart
		}
	}
}
/*
The logs imported from orinal "log" are global log.
generally, when there are some framework level infos or serious errors, it will use global log.
*/
func InitGlobalLog(name string){
	rootdir := os.Getenv("GOPATH")
	fpath := path.Join(rootdir, "global.log")
	file, e := os.OpenFile(fpath, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
	if e != nil{
		log.Fatalf("[%s] [Fail to open log file: %s] %s\n", fpath, e)
	}
	log.SetOutput(io.MultiWriter(file, os.Stdout))

}
// The logs imported from "slog" are spider log.
func InitSpiderLog(name string){
	slog.Init(name)
}
func (this *SpiderManager) Run(){
	if !this.start{
		return
	}
	this.checkDownloadPartPtr()
	this.checkResolverPartPtr()

	numUnit := config.Get(this.name, config.COMMON_CONF, "NumUnit").Int()
	this.initLayer = layer.InitLayer(this.pipeline, fmt.Sprintf("%s-%s-%d", this.name, "Layer", 0), len(this.seeds), numUnit, 0)

	layer.Start(this.initLayer, this.pipeline, this.maxDepth)
	slog.Close(this.name)
}
func (this *SpiderManager) GetName() string{
	return this.name
}

func (this *SpiderManager) check(b bool, info string){
	if b{
		log.Println(info)
		this.start = false
	}
}
