package main

import (
	"log"
	"spider"
	"config"
	"parts/baidu"
	"parts/sina"
)


/*
parts/baidu and parts/sina are just demos.
not really implement download and resolve in every web layer.
you can use this sample and copy them to create your own download and
resolve logic in every layer of the specify website.
*/

func main(){
	// baidu layer-3 using default downloader and resolver.
	m1 := spider.NewSpiderManager("baidu", config.NewLink("https://www.baidu.com"))
	m1.AddLayerStrategy(0, baidu.NewDP0, baidu.NewRP0)
	m1.AddLayerStrategy(1, baidu.NewDP1, baidu.NewRP1)

	// all layers in sina spider are default downloader and resolver except layer-1.
	m2 := spider.NewSpiderManager("sina", config.NewLink("www.sina.com"))
	m2.AddLayerStrategy(1, sina.NewDP1, sina.NewRP1)

	spider.AppendManager(m1)
	spider.AppendManager(m2)
	spider.RunAll()

	log.Println("all spider Done.")
}

