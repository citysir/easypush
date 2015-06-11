package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	log "github.com/log4go"
	"net/http"
	"strconv"
)

func syncid(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	token := r.FormValue("token")
	cv, err := strconv.ParseUint(r.FormValue("cv"), 10, 64) //201506031648
	if err != nil || cv < Conf.MinClientVersion {
		log.Error("param cv %s err", r.FormValue("cv"))
	}

	typ, err := strconv.Atoi(r.FormValue("type"))
	if err != nil {
		w.Write(ToErrorJson(PARAM_INVALID, fmt.Sprintf("type %s invalid", r.FormValue("type"))))
		return
	}
	_, err = strconv.ParseInt(r.FormValue("lastid"), 10, 64)
	if err != nil {
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
	remCount, err4 := Global.Redis.ZRemRangeByScore(typeKey, "-inf", r.FormValue("lastid"))
	if err4 != nil {
		w.Write(ToErrorJson(SERVER_EXCEPTION, err4.Error()))
		return
	}

	result := map[string]interface{}{"r": SUCCESS, "count": remCount}
	w.Write(ToJson(result))
}
