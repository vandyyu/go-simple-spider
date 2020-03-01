package layer
import (
	"log"
	"sync"
)
/*
layer write data to pipeline, pipeline will cached data first,
when current layer has done, cached data in pipeline will update
to next layer. if cached datas for persistance reach MAX, they will flush to somewhere specified.
operating for persistance will do recurrently as long as IsFull() return true.

Note: 
maybe you should consider more about the concurrency problem, because there're many concurrent read and write during the Pipeline working.
*/
type IPipeline interface{
	Write(d interface{})
	Update()
	GetNextLayer() ILayer
	GetFirstPartProxy() IPartProxy
	GetName() string
	IDump
}
type IDump interface{
	Flush()
	Notify()
	Wait()
	InitDumpSignal()
	closeDump()
}
type Dump struct{
	signal chan struct{}
}
func (this *Dump) InitDumpSignal(){
	this.signal = make(chan struct{})
}
func (this *Dump) Flush(){
}
func (this *Dump) Notify(){
	<-this.signal
}
func (this *Dump) Wait(){
	this.signal <- struct{}{}
}
func (this *Dump) closeDump(){
	close(this.signal)
	for {
		if _, ok := <-this.signal;!ok{
			break
		}
	}
}

func Start(initLayer ILayer, pipeline IPipeline, maxDepth int){
	var curLayer = initLayer
	var wg sync.WaitGroup
	var stopDump = false
	var watchDone = make(chan struct{})

	pipeline.InitDumpSignal()

	if curLayer != nil{
		// persistance
		go func(){
			watching(pipeline, &wg, &stopDump)
			watchDone <- struct{}{}
		}()
	}
	for i:=0;i < maxDepth;i++{
		if curLayer == nil{
			log.Println("layer-", i, " is nil, stop layer runing.")
			break
		}
		curLayer.Run()
		if i == maxDepth - 1{
			break
		}
		pipeline.Update()
		curLayer = pipeline.GetNextLayer()
	}
	pipeline.Flush()
	wg.Wait()
	stopDump = true
	pipeline.Notify()
	<-watchDone
}
func watching(pipeline IPipeline, wg *sync.WaitGroup, stopDump *bool){
	var muFlush sync.Mutex
	for{
		muFlush.Lock()
		pipeline.Wait()  // wait dump
		if *stopDump{
			pipeline.closeDump()
			break
		}
		wg.Add(1)
		go func(){
			pipeline.Flush()
			wg.Done()
		}()
		muFlush.Unlock()
	}
	muFlush.Unlock()
}
