package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	log "github.com/log4go"
	"net/http"
	"strconv"
)

func auth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := r.FormValue("user")
	password := r.FormValue("password")
	did := r.FormValue("did")
	cv, err := strconv.ParseUint(r.FormValue("cv"), 10, 64) //201506031648
	if err != nil || cv < Conf.MinClientVersion {
		log.Error("param cv %s err", r.FormValue("cv"))
	}

	log.Debug("auth: %s %s %s %d", user, password, did, cv)

	token, tokenData, err := AuthenticateUser(user, password, did)
	if err != nil {
		log.Error("AuthenticateUser err user=%s, password=%s, err=%s", user, password, err.Error())
		w.Write(ToErrorJson(AUTHENTICATION_EXCEPTION, ""))
		return
	}

	tokenDataBytes, err := json.Marshal(tokenData)
	if err != nil {
		log.Error("json.Marshal tokenData err err=%s", err.Error())
		w.Write(ToErrorJson(SERVER_EXCEPTION, ""))
		return
	}

	log.Debug("Redis.Setex token data token=%s, data=%s", token, string(tokenDataBytes))
	err = Global.Redis.Setex(token, Conf.TokenTimeout, string(tokenDataBytes))
	if err != nil {
		log.Error("Redis.Setex err err=%s", err.Error())
		w.Write(ToErrorJson(SERVER_EXCEPTION, ""))
		return
	}

	result := map[string]interface{}{"r": SUCCESS, "token": token, "node": FindNodeHost(user)}
	w.Write(ToJson(result))
}
