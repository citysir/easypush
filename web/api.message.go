package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	log "github.com/log4go"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func message(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	token := r.FormValue("token")
	cv, err := strconv.ParseUint(r.FormValue("cv"), 10, 64) //201506031648
	if err != nil || cv < Conf.MinClientVersion {
		log.Error("param cv %s err", r.FormValue("cv"))
	}

	tokenData, err := GetTokenData(token)
	if err != nil {
		log.Error("GetTokenData err token=%s, err=%s", token, err.Error())
		w.Write(ToErrorJson(SERVER_EXCEPTION, ""))
		return
	}

	if tokenData == nil {
		w.Write(ToErrorJson(TOKEN_INVALID, "token invalid or timeout"))
		return
	}

	ids := strings.Split(r.FormValue("id"), ",")
	messageBytesArray, err := getMessageBytesArrayParallel(ids)
	if err != nil {
		w.Write(ToErrorJson(SERVER_EXCEPTION, err.Error()))
		return
	}

	if len(messageBytesArray) == 0 {
		w.Write(ToJson(map[string]interface{}{"r": DATA_NOT_EXISTS}))
		return
	}

	messages := make([]*Message, len(messageBytesArray))
	for i, messageBytes := range messageBytesArray {
		message := &Message{}
		err2 := json.Unmarshal(messageBytes, message)
		if err2 != nil {
			log.Error("json.Unmarshal err, message=%s", string(messageBytes))
			w.Write(ToErrorJson(SERVER_EXCEPTION, ""))
			return
		}
		messages[i] = message
	}
	w.Write(ToJson(map[string]interface{}{"r": SUCCESS, "msgs": messages}))
}

type ParallelGetMessageJob struct {
	index     int
	messageId string
}

//并行访问redis，提高速度
func getMessageBytesArrayParallel(messageIds []string) (messageBytesArray [][]byte, finalError error) {
	messageBytesArray = make([][]byte, len(messageIds))
	jobChan := make(chan *ParallelGetMessageJob, runtime.NumCPU())
	var wg sync.WaitGroup
	for i, messageId := range messageIds {
		jobChan <- &ParallelGetMessageJob{i, messageId}
		wg.Add(1)
		go func() {
			job := <-jobChan
			messageBytes, err := Global.Redis.Get(job.messageId)
			if err != nil {
				finalError = err
			}
			messageBytesArray[job.index] = messageBytes
			wg.Done()
		}()
	}
	wg.Wait()
	return messageBytesArray, finalError
}
