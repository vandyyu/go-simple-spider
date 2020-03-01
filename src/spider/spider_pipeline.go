package spider
import (
	"layer"
	"sync"
	"slog"
	"config"
	"fmt"
)

type SpiderPipeline struct{
	layer.Dump

	seeds []*config.Link
	i_sed int

	cachedLink []*config.Link
	cachedRawData []*config.RawData
	cachedText []*config.Text
	dumpedText []*config.Text

	muLink sync.Mutex
	muRawData sync.Mutex
	muText sync.Mutex

	name string
	depth int  // layer depth
}

func NewSpiderPipeline(name string) *SpiderPipeline{
	pipeline := new(SpiderPipeline)
	pipeline.name = name
	pipeline.depth = 0
	pipeline.seeds = make([]*config.Link, 0, 128)
	pipeline.cachedLink = make([]*config.Link, 0, 128)
	pipeline.cachedRawData = make([]*config.RawData, 0, 128)
	pipeline.cachedText = make([]*config.Text, 0, 128)
	return pipeline
}

/*
1. flush some cached Text data to local file or database, etc.
2. the number of flushing cached Text can setup by {config.NumCachedText}.
3. framework will auto call Flush() as long as the number of cached Text reaches {config.NumCachedText}.
4. run in a goroutine, do not share data with outside unless you known how to do sync.
there are many goroutine to do this function, but {this.dumpedText} has it own memory
and truely has necessary value data in every goroutine because they're copied from 
outside with thread-safe. 
*/
func (this *SpiderPipeline) Flush(){
	DumpData(this.dumpedText)
}
func (this *SpiderPipeline) GetNextLayer() layer.ILayer{
	this.depth += 1
	numUnit := config.Get(this.name, config.COMMON_CONF, "NumUnit").Int()
	ret := layer.InitLayer(this, this.name, len(this.seeds), numUnit, this.depth)
	return ret
}

func (this *SpiderPipeline) AddSeeds(seeds ...*config.Link){
	this.seeds = append(this.seeds, seeds...)
}

// all threads has done, this function only executed in single-thread, thread-safe
// update and turn to next layer, all about this layer will be cleared.
func (this *SpiderPipeline) Update(){
	this.seeds = this.cachedLink
	this.i_sed = 0
	this.cachedLink = make([]*config.Link, 0, 128)
	this.cachedRawData = make([]*config.RawData, 0, 128)
}

// many parallel tasks call Write(), must do sync.
func (this *SpiderPipeline) Write(d interface{}){
	if d == nil{
		msg := fmt.Sprintf("skip written data of %s in %d layer, data is nil", this.name, this.depth)
		slog.Warning(this.name, msg)
		return
	}
	switch d := d.(type){
		case *config.Link:
			this.muLink.Lock()
			this.cachedLink = append(this.cachedLink, d)
			this.muLink.Unlock()
		case []*config.Link:
			this.muLink.Lock()
			this.cachedLink = append(this.cachedLink, d...)
			this.muLink.Unlock()
		case *config.RawData:
			this.muRawData.Lock()
			// TODO: maybe the downloaded rawData in current layer can be setup as cached to exceed spider speed. Not implemented.
			this.cachedRawData = append(this.cachedRawData, d)
			this.muRawData.Unlock()
		case *config.Text:
			this.muText.Lock()
			this.cachedText = append(this.cachedText, d)
			numCachedText := config.Get(this.name, config.COMMON_CONF, "NumCachedText").Int()
			if len(this.cachedText) >= numCachedText{
				this.dumpedText = make([]*config.Text, len(this.cachedText), cap(this.cachedText))
				copy(this.dumpedText, this.cachedText)
				this.Notify() // notify to flush data to persistant
				this.cachedText = this.cachedText[:0]
			}
			this.muText.Unlock()
		default:
			slog.Error(this.name, fmt.Sprintf("cannot recognize type %T, write failed!", d))
	}
}

// only in initializing processing, single-thread, and thread-safe
func (this *SpiderPipeline) GetFirstPartProxy() layer.IPartProxy{
	var partProxy layer.IPartProxy
	if this.i_sed < len(this.seeds){
		NewDP := config.Get(this.name, config.PART_CONF, "NewDownloadPartMap").DPPtr(this.depth)
		if NewDP == nil{
			return nil
		}
		part := NewDP("DownloadPart", this.seeds[this.i_sed])
		part.SetPipeline(this)
		part.SetDownloader(part)
		partProxy = NewPartProxy(part)
		this.i_sed += 1
	}
	return partProxy
}
func (this *SpiderPipeline) GetName() string{
	return this.name
}


