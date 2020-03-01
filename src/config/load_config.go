package config

import (
	"os"
	"path"
	"tool"
	"encoding/json"
	"log"
	"io/ioutil"
//	"fmt"
)

func loadConfFile(c interface{}, dirname string, confName string){
	var (
		err error
		f *os.File
		bytes []byte
	)
	rootdir := path.Join(os.Getenv("GOPATH"), "configs", dirname)
	if !tool.ExistDir(rootdir){
		os.MkdirAll(rootdir, 0777)
	}
	fpath := path.Join(rootdir, confName)
	if !tool.ExistFile(fpath){
		log.Printf("config file [%s] doesn't exists! generating default config file, and using default config!\n", fpath)
		bytes, err = json.MarshalIndent(c, "", "\t")
		if err != nil{
			log.Fatalf("%v\n", err)
		}
		f, err = os.Create(fpath)
		defer f.Close()
		if err == nil{
			f.Write(bytes)
		}else{
			log.Fatalf("[%s.load():AAA] create file failed! file: %s. [%v]\n", confName, fpath, err)
		}
	}else{
		f, err = os.OpenFile(fpath, os.O_RDONLY, 0777)
		if err == nil{
			bytes, err = ioutil.ReadAll(f)
			if err == nil{
				err = json.Unmarshal(bytes, c)
				if err != nil{
					log.Fatalf("[%s.load():BBB] failed to resolve config file: %s. [%v]\n", confName, fpath, err)
				}
			}else{
				log.Fatalf("[%s.load():CCC] failed to read config file: %s. [%v]\n", confName, fpath, err)
			}
		}else{
			log.Fatalf("[%s.load():DDD] failed to open file: %s. [%v]\n", confName, fpath, err)
		}
	}
//	fmt.Println("%v\n", c)
}
