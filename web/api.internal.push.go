package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/citysir/easypush/hash"
	"github.com/julienschmidt/httprouter"
	log "github.com/log4go"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func push(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	inputVc, err := strconv.ParseUint(r.FormValue("vc"), 10, 32)
	if err != nil {
		w.Write(ToErrorJson(PARAM_INVALID, fmt.Sprintf("invalid vc %s", inputVc)))
		return
	}

	accessKey := r.FormValue("ac")
	keyString := r.FormValue("key")
	messageString := r.FormValue("message")

	secretKey, exists := Conf.AccessKeys[accessKey]
	if !exists {
		w.Write(ToErrorJson(AUTHENTICATION_EXCEPTION, fmt.Sprintf("invalid ac %s", accessKey)))
		return
	}

	vcData := bytes.NewBufferString(keyString)
	vcData.WriteString(messageString)
	vcData.WriteString(secretKey)
	calcVc := hash.HashCrc32(vcData.Bytes())

	if uint32(inputVc) != calcVc {
		w.Write(ToErrorJson(PARAM_INVALID, fmt.Sprintf("error vc %d", inputVc)))
		return
	}

	message := &Message{}
	err = json.Unmarshal([]byte(messageString), message)
	if err != nil {
		w.Write(ToErrorJson(PARAM_INVALID, fmt.Sprintf("message %s invalid", r.FormValue("message"))))
		return
	}

	log.Debug("push: %s %s", keyString, messageString)

	keys := strings.Split(keyString, ",")
	err = saveMessage(keys, messageString, message)
	if err != nil {
		w.Write(ToErrorJson(SERVER_EXCEPTION, fmt.Sprintf("saveMessage failed err=%s", err.Error())))
		return
	}

	nodeKeys := map[string][]string{}
	for _, key := range keys {
		host := FindNodeHost(key)
		if _, exists := nodeKeys[host]; !exists {
			nodeKeys[host] = []string{}
		}
		nodeKeys[host] = append(nodeKeys[host], key)
	}

	rpcClientKeys := map[*NodeRpcClient][]string{}
	for node, keys := range nodeKeys {
		rpcAddr := fmt.Sprintf("%s:%d", node, Conf.NodeRpcPort)
		rpcClient, err := NewNodeRpcClient(rpcAddr)
		if err != nil {
			w.Write(ToErrorJson(SERVER_EXCEPTION, fmt.Sprintf("error get rpc client %s", rpcAddr)))
			return
		}
		rpcClientKeys[rpcClient] = keys
	}

	message.Data = nil //只推送通知，不推送内容
	allFailedKeys := []string{}
	for rpcClient, keys := range rpcClientKeys {
		args := &PushMessageArgs{Keys: keys, Message: message}
		failedKeys, err := rpcClient.CallPushMessages(args)
		if err != nil {
			log.Error("failed CallPushMessages %s", strings.Join(keys, ","))
			allFailedKeys = append(allFailedKeys, keys...)
		} else {
			allFailedKeys = append(allFailedKeys, failedKeys...)
		}
	}

	result := map[string]interface{}{"r": SUCCESS, "fails": allFailedKeys}
	w.Write(ToJson(result))
}

func saveMessage(keys []string, messageString string, message *Message) error {
	return saveMessageParallel(keys, messageString, message)
}

type ParallelSaveMessageJob struct {
	key           string
	messageString string
	message       *Message
}

//并行访问redis，提高速度
func saveMessageParallel(keys []string, messageString string, message *Message) (finalError error) {
	jobChan := make(chan *ParallelSaveMessageJob, runtime.NumCPU())
	var wg sync.WaitGroup
	for _, key := range keys {
		jobChan <- &ParallelSaveMessageJob{key, messageString, message}
		wg.Add(1)
		go func() {
			job := <-jobChan
			err := saveMessageToRedis(job)
			if err != nil {
				finalError = err
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return finalError
}

func saveMessageToRedis(job *ParallelSaveMessageJob) error {
	messageTypeKey := MakeTypeKey(job.key, job.message.Type)
	log.Debug("saveMessageToRedis key=%s, type=%d, typeKey=%s", job.key, job.message.Type, messageTypeKey)
	err := Global.Redis.Setex(strconv.FormatInt(job.message.Id, 10), Conf.MessageTimeout, job.messageString)
	if err != nil {
		return err
	}
	_, err = Global.Redis.ZAdd(messageTypeKey, map[string]float64{strconv.FormatInt(job.message.Id, 10): float64(job.message.Id)})
	if err != nil {
		return err
	}
	return nil
}
