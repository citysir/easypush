package main

import (
	"github.com/citysir/easypush/perf"
	"github.com/julienschmidt/httprouter"
	log "github.com/log4go"
	"net/http"
	"os"
)

func StartWeb() {
	log.Info("Start web listen addr %s", Conf.WebBind)
	router := httprouter.New()
	router.POST("/v1/auth", auth)
	router.POST("/v1/syncid", syncid)
	router.GET("/v1/messages", messages)
	router.GET("/v1/message", message)
	router.POST("/v1/message", message)
	router.POST("/v1/uploaddid", uploaddid)
	if err := http.ListenAndServe(Conf.WebBind, router); err != nil {
		log.Error("http.ListenAdServe(\"%s\") error(%v)", Conf.WebBind, err)
		panic(err)
	}
}

func StartPush() {
	log.Info("Start push listen addr: %s", Conf.PushBind)
	router := httprouter.New()
	router.POST("/v1/push", push)
	if err := http.ListenAndServe(Conf.PushBind, router); err != nil {
		log.Error("http.ListenAdServe(\"%s\") error(%v)", Conf.PushBind, err)
		panic(err)
	}
}

func main() {
	InitGlobal()

	log.Info("Start perf listen addr %s", Conf.PerfBind)
	perf.BindAddr(Conf.PerfBind)

	go StartWeb()
	go StartPush()
	go StartRpc()

	SignalWatchRegister(quit, os.Kill, os.Interrupt) //, syscall.SIGTERM)
	SignalWatchRun()

	<-make(chan int)
}

func quit() {
	log.Debug("quit")
	CloseRedisClient(Global.Redis)
}
