package main

import (
	"github.com/julienschmidt/httprouter"
	log "github.com/log4go"
	"net/http"
	"strconv"
)

func uploaddid(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	token := r.FormValue("token")
	user := r.FormValue("user")
	did := r.FormValue("did")

	cv, err := strconv.ParseUint(r.FormValue("cv"), 10, 64) //201506031648
	if err != nil || cv < Conf.MinClientVersion {
		log.Error("param cv %s err", r.FormValue("cv"))
	}

	log.Debug("uploaddid: %s %s %s", token, user, did)

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

	if tokenData.User != user {
		w.Write(ToErrorJson(TOKEN_INVALID, "token invalid"))
		return
	}

	err = saveUserDid(user, did)
	if err == nil {
		w.Write(ToJson(map[string]interface{}{"r": SUCCESS}))
	} else {
		w.Write(ToErrorJson(SERVER_EXCEPTION, err.Error()))
	}
}

func saveUserDid(user, did string) error {
	return nil
}
