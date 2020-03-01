package config
import (
//	"reflect"
	"log"
)
func INIT_CONFIG_ERR(){
	log.Fatalln("please call config/init_config.go:InitConfig(name) first!")
}
func Get(name string, confType ConfigType, key string) *ConfigValue{
	var cv *ConfigValue
	switch confType{
		case COMMON_CONF:
			muComm.Lock()
			if commConfMap == nil{
				INIT_CONFIG_ERR()
			}
			cv = commConfMap[name].Get(key)
			muComm.Unlock()
		case LOG_CONF:
			muLog.Lock()
			if logConfMap == nil{
				INIT_CONFIG_ERR()
			}
			cv = logConfMap[name].Get(key)
			muLog.Unlock()
		case HTTP_CONF:
			muHttp.Lock()
			if httpConfMap == nil{
				INIT_CONFIG_ERR()
			}
			cv = httpConfMap[name].Get(key)
			muHttp.Unlock()
		case PART_CONF:
			muPart.Lock()
			if partConfMap == nil{
				INIT_CONFIG_ERR()
			}
			cv = partConfMap[name].Get(key)
			muPart.Unlock()
		default:
			log.Fatalln("[config/config_factory.go:Get()] unknown config type!")
	}
	return cv


/*  // reflect efficiency is too low.
	conf := newConf(name)
	t := reflect.TypeOf(conf).Elem()
	v := reflect.ValueOf(conf).Elem()
	for i:=0;i < t.NumField();i++{
		if t.Field(i).Name == key{
			v.Field(i)
		}
	}
*/
}
