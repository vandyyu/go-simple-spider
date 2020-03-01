package config
import (
	"log"
	"sync"
	"runtime/debug"
	"fmt"
)

var commConfMap map[string]*CommonConfig
var logConfMap map[string]*LogConfig
var httpConfMap map[string]*HttpConfig
var partConfMap map[string]*PartConfig

var muComm sync.Mutex
var muLog sync.Mutex
var muHttp sync.Mutex
var muPart sync.Mutex

type ConfigType int
const (
	COMMON_CONF ConfigType = iota
	LOG_CONF
	HTTP_CONF
	PART_CONF
)

type ConfigValue struct{
	data interface{}
	name string
	a int
	b int64
	c bool
	d string
	e map[int]DownloadPartPtr
	f map[int]ResolverPartPtr
	// g float32
	flag int
}
func NewConfigValue(name string, data interface{}) *ConfigValue{
	c := new(ConfigValue)
	c.data = data
	c.flag = -1
	c.name = name
	return c
}
func (this *ConfigValue) assert(b bool, info string){
	if !b{
		debug.PrintStack()
		log.Fatalf("[%s] [config/init_config.go:convert() failed!]. type of %v is %T, not type %s. \n", this.name, this.data, this.data, info)
	}
}
func (this *ConfigValue) Int() int{
	this.convert()
	this.assert(this.flag == 0, "int")
	return this.a
}
func (this *ConfigValue) Int64() int64{
	this.convert()
	this.assert(this.flag == 1, "int64")
	return this.b
}
func (this *ConfigValue) Bool() bool{
	this.convert()
	this.assert(this.flag == 2, "bool")
	return this.c
}
func (this *ConfigValue) String() string{
	this.convert()
	this.assert(this.flag == 3, "string")
	return this.d
}
func (this *ConfigValue) DPPtr(depth int) DownloadPartPtr{
	this.convert()
	this.assert(this.flag == 4, "DownloadPartPtr")
	if v, ok := this.e[depth];ok{
		return v
	}else{
		log.Printf(fmt.Sprintf("[%s] DownloadPartPtr for layer-%d doesn't setup!", this.name,  depth))
		return nil
	}
}
func (this *ConfigValue) RPPtr(depth int) ResolverPartPtr{
	this.convert()
	this.assert(this.flag == 5, "ResolverPartPtr")
	if v, ok := this.f[depth];ok{
		return v
	}else{
		log.Printf("[%s] ResolverPartPtr for layer-%d doesn't setup!", this.name, depth)
		return nil
	}
}
func (this *ConfigValue) convert(){
	switch v := this.data.(type){
	case int:
		this.a = v
		this.flag = 0
	case int64:
		this.b = v
		this.flag = 1
	case bool:
		this.c = v
		this.flag = 2
	case string:
		this.d = v
		this.flag = 3
	case map[int]DownloadPartPtr:
		this.e = v
		this.flag = 4
	case map[int]ResolverPartPtr:
		this.f = v
		this.flag = 5
	default:
		log.Fatalf("[%s] config/init_config.go: ConfigValue.convert() cannot recognize type: [%T]\n", this.name,  v)
	}
}

func InitConfig(name string, dpptrMap map[int]DownloadPartPtr, rpptrMap map[int]ResolverPartPtr, args ...interface{}){
	muComm.Lock()
	if commConfMap == nil{
		commConfMap = make(map[string]*CommonConfig)
	}
	if _, ok := commConfMap[name];!ok{
		commConfMap[name] = newCommonConfig(name)
	}
	muComm.Unlock()

	muLog.Lock()
	if logConfMap == nil{
		logConfMap = make(map[string]*LogConfig)
	}
	if _, ok := logConfMap[name];!ok{
		logConfMap[name] = newLogConfig(name)
	}
	muLog.Unlock()

	muHttp.Lock()
	if httpConfMap == nil{
		httpConfMap = make(map[string]*HttpConfig)
	}
	if _, ok := httpConfMap[name];!ok{
		httpConfMap[name] = newHttpConfig(name)
	}
	muHttp.Unlock()

	muPart.Lock()
	if partConfMap == nil{
		partConfMap = make(map[string]*PartConfig)
	}
	if _, ok := partConfMap[name];!ok{
		pc := newPartConfig(name)
		if dpptrMap != nil{
			pc.SetDPPtrMap(dpptrMap)
		}else{
			debug.PrintStack()
			log.Fatalln(fmt.Sprintf("[%s] DownloadPartPtr shouldn't be nil.", name))
		}
		if rpptrMap != nil{
			pc.SetRPPtrMap(rpptrMap)
		}else{
			debug.PrintStack()
			log.Fatalln(fmt.Sprintf("[%s] ResolverPartPtr shouldn't be nil.", name))
		}
		partConfMap[name] = pc
	}
	muPart.Unlock()
}
