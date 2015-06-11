package main

import (
	"github.com/citysir/easypush/perf"
	log "github.com/log4go"
	"net/http"
)

func nodeListen(bind string) {
	httpServeMux := http.NewServeMux()
	httpServeMux.HandleFunc("/v1/ws", serveWs)
	if err := http.ListenAndServe(bind, httpServeMux); err != nil {
		log.Error("http.ListenAdServe(\"%s\") error(%v)", bind, err)
		panic(err)
	}
}

func StartNode() {
	log.Info("Start node listen addr: %s", Conf.NodeBind)
	go nodeListen(Conf.NodeBind)
}

func main() {
	go hub.run()

	log.Info("Start perf listen addr: %s", Conf.PerfBind)
	perf.BindAddr(Conf.PerfBind)

	StartStat()
	StartNode()
	StartRpc()

	<-make(chan int)
}
