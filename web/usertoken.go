package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

const TOKEN_SALT = "@L-k3CO2,eJLpCgBuy3^iR8P"

type TokenData struct {
	User     string `json:"u"`
	UserType int8   `json:"ut"`
	Did      string `json:"did"`
}

func AuthenticateUser(user, password, did string) (string, *TokenData, error) {
	// do some check
	userType := int8(1)
	token := MakeUserToken(user, did)
	return token, &TokenData{User: user, UserType: userType, Did: did}, nil
}

func MakeUserToken(key, did string) string {
	seed := time.Now().Unix() / (60 * 10) //保证10分钟之内生成的token相同
	return MakeMd5Str([]byte(fmt.Sprintf("%s%s%d%s", key, did, seed, TOKEN_SALT)))
}

func GetTokenData(token string) (*TokenData, error) {
	data, err := Global.Redis.Get(token)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	tokenData := &TokenData{}
	err = json.Unmarshal(data, tokenData)
	if err != nil {
		return nil, err
	}
	return tokenData, nil
}

// 一类消息存放在redis中的key
func MakeTypeKey(key string, typ int16) string {
	return MakeMd5Str([]byte(fmt.Sprintf("%s%d", key, typ)))
}

func MakeMd5Str(data []byte) string {
	return hex.EncodeToString(MakeMd5(data))
}

func MakeMd5(data []byte) []byte {
	h := md5.New()
	h.Write(data)
	return h.Sum(nil)
}
