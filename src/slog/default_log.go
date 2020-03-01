package slog
import (
	"log"
	"io"
	"time"
	"os"
	"runtime/debug"
	"sync"
	"path"
	"config"
	"fmt"
	"tool"
	"strings"
)
/*
1. diffent website obtain different slog entity.
2. it will auto create a new log file when logfile size exceed {LogFileSize}.
*/
var logMap = make(map[string]*SpiderLog)
var muInitLog sync.Mutex

type SpiderLog struct{
	file *os.File
	info *log.Logger
	warning *log.Logger
	err *log.Logger
	counterMap map[string]int64

//	muInfo sync.Mutex
//	muWarn sync.Mutex
//	muErr sync.Mutex
	muFinish sync.Mutex
	finished bool
	heartBeat chan struct{}
	heartBeatFeedback chan struct{}
	closeBlocking chan struct{}
}
func newSpiderLog(fpath string) *SpiderLog{
	sl := new(SpiderLog)
	sl.counterMap = make(map[string]int64)
	sl.updateLogFile(fpath)
	sl.finished = false
	sl.closeBlocking = make(chan struct{})
	sl.heartBeat = make(chan struct{})
	sl.heartBeatFeedback = make(chan struct{})
	return sl
}
func (this *SpiderLog) updateLogFile(fpath string){
	this.file = createNewFile(fpath)
	this.info = log.New(io.MultiWriter(this.file, os.Stdout), "[INFO]: ", log.Ldate | log.Ltime)
	this.warning = log.New(io.MultiWriter(this.file, os.Stdout), "[WARNING]: ", log.Ldate | log.Ltime)
	this.err = log.New(io.MultiWriter(this.file, os.Stderr), "[ERROR]: ", log.Ldate | log.Ltime)
}
func createNewFile(fpath string) *os.File{
	file, e := os.OpenFile(fpath, os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
	if e != nil{
		log.Fatalln("[Fail to open log file: ", fpath, "] ", e)
	}
	return file
}
func getFilePath(rootdir string) string{
	timeStr := time.Now().Format("2006-01-02_15-04-05.log")
	return path.Join(rootdir, timeStr)
}
func Init(subdir string){
	initLogFile(subdir, false)
	updateLogFile(subdir)
}
func initLogFile(subdir string, update bool){
	muInitLog.Lock()
	rootdir := path.Join(config.Get(subdir, config.LOG_CONF, "LogRootDir").String(), subdir)
	if !tool.ExistDir(rootdir){
		os.MkdirAll(rootdir, 0777)
	}
	fpath := getFilePath(rootdir)
	if update{
		logMap[subdir].file.Close()
		if tool.ExistFile(fpath){
			logMap[subdir].counterMap[fpath] += 1
			fpath = fmt.Sprintf("%s.%d", fpath, logMap[subdir].counterMap[fpath])
		}
		logMap[subdir].updateLogFile(fpath)
	}else{
		if _, ok := logMap[subdir];!ok{
			logMap[subdir] = newSpiderLog(fpath)
		}
	}
	muInitLog.Unlock()
}
func getSpiderLog(subdir string) *SpiderLog{
	var spiderLog *SpiderLog
	muInitLog.Lock()
	if _, ok := logMap[subdir];!ok{
		log.Fatalln("log for ", subdir, " doesn't initialize!")
	}else{
		spiderLog = logMap[subdir]
	}
	muInitLog.Unlock()
	return spiderLog
}
func check(info string, v ...interface{}){
	if len(v) == 0{
		log.Fatalln(info)
	}
}
func combineParams(v ...interface{}) string{
	tmp := make([]string, 0, 8)
	for _, s := range(v){
		s := fmt.Sprintf("%v", s)
		if len(s) > 2{
			if s[:1] == "[" && s[len(s)-1:] == "]"{
				s = s[1 : (len(s) - 1)]
			}
		}
		tmp = append(tmp, s)
		println(fmt.Sprintf("%s", s))
	}
	return strings.Join(tmp, "\n")
}
func Info(pipelineName string, v ...interface{}){
	check("slog.Info(pipelineName string, v...interface{}), you must give pipelineName and info message!", v)
	spiderLog := getSpiderLog(pipelineName)
	if spiderLog != nil{
		params := combineParams(v)
		spiderLog.info.Println(fmt.Sprintf("<%s>\n%s", pipelineName, params))
		spiderLog.file.Sync()
		<-spiderLog.heartBeat
		<-spiderLog.heartBeatFeedback
	}
}
func Warning(pipelineName string, v ...interface{}){
	check("slog.Warning(pipelineName string, v...interface{}), you must give pipelineName and warning message!", v)
	spiderLog := getSpiderLog(pipelineName)
	if spiderLog != nil{
		params := combineParams(v)
		spiderLog.warning.Println(fmt.Sprintf("<%s>\n%s", pipelineName, params))
		spiderLog.file.Sync()
		<-spiderLog.heartBeat
		<-spiderLog.heartBeatFeedback
	}
}
func Error(pipelineName string, v ...interface{}){
	check("slog.Error(pipelineName string, v...interface{}), you must give pipelineName and error message!", v)
	spiderLog := getSpiderLog(pipelineName)
	if spiderLog != nil{
		params := combineParams(v)
		spiderLog.err.Println(fmt.Sprintf("<%s>\n%s\n%s", pipelineName, params, debug.Stack()[230:]))
		spiderLog.file.Sync()
		<-spiderLog.heartBeat
		<-spiderLog.heartBeatFeedback
	}
}
func Close(pipelineName string){
	spiderLog := getSpiderLog(pipelineName)
	spiderLog.muFinish.Lock()
	spiderLog.finished = true
	spiderLog.muFinish.Unlock()
	<-spiderLog.heartBeat
	<-spiderLog.heartBeatFeedback
	<-spiderLog.closeBlocking
}
func updateLogFile(subdir string){
	spiderLog := getSpiderLog(subdir)
	go func(){
		for {
			fi, err := spiderLog.file.Stat()
			if err != nil{
				log.Fatalln("log file appears error when update log file.", err)
			}
			spiderLog.heartBeat <- struct{}{}
			if fi.Size() > config.Get(subdir, config.LOG_CONF, "LogFileSize").Int64() {
				initLogFile(subdir, true)
			}
			spiderLog.muFinish.Lock()
			if spiderLog.finished{
				break
			}
			spiderLog.muFinish.Unlock()

			spiderLog.heartBeatFeedback <- struct{}{}
		}
		spiderLog.file.Close()
		close(spiderLog.heartBeat)
		close(spiderLog.heartBeatFeedback)

		tool.ClearChanAfterClosed(spiderLog.heartBeat)
		tool.ClearChanAfterClosed(spiderLog.heartBeatFeedback)

		spiderLog.muFinish.Unlock()
		spiderLog.closeBlocking <- struct{}{}
	}()
}
