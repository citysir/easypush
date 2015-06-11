package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	log "github.com/log4go"
	"net/http"
	"strconv"
)

func messages(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	token := r.FormValue("token")
	cv, err := strconv.ParseUint(r.FormValue("cv"), 10, 64) //201506031648
	if err != nil || cv < Conf.MinClientVersion {
		log.Error("param cv %s err", r.FormValue("cv"))
	}

	typ, err := strconv.Atoi(r.FormValue("type"))
	if err != nil {
		log.Error("param type %s err", r.FormValue("type"))
		w.Write(ToErrorJson(PARAM_INVALID, fmt.Sprintf("type %s invalid", r.FormValue("type"))))
		return
	}
	lastMessageId, err := strconv.ParseInt(r.FormValue("lastid"), 10, 64)
	if err != nil {
		log.Error("param lastid %s err", r.FormValue("lastid"))
		w.Write(ToErrorJson(PARAM_INVALID, fmt.Sprintf("lastid %s invalid", r.FormValue("lastid"))))
		return
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

	typeKey := MakeTypeKey(tokenData.User, int16(typ))
	log.Debug("ZRangeByScore token=%s, user=%s, type=%d, typeKey=%s", token, tokenData.User, int16(typ), typeKey)
	messageIdStrings, err := Global.Redis.ZRangeByScore(typeKey, strconv.FormatInt(lastMessageId+1, 10), "+inf", false, false, 0, 0)
	if err != nil {
		w.Write(ToErrorJson(SERVER_EXCEPTION, err.Error()))
		return
	}

	messageIds := make([]int64, len(messageIdStrings))
	for i, messageIdString := range messageIdStrings {
		messageIds[i], err = strconv.ParseInt(messageIdString, 10, 64)
		if err != nil {
			w.Write(ToErrorJson(SERVER_EXCEPTION, err.Error()))
			return
		}
	}

	result := map[string]interface{}{"r": SUCCESS, "ids": messageIds}
	w.Write(ToJson(result))
}
