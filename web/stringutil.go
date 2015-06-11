package main

import (
	"encoding/json"
	"strconv"
)

func StringArrayToInt64(array []string) ([]int64, error) {
	var err error
	intArray := make([]int64, len(array))
	for i, value := range array {
		intArray[i], err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return intArray, err
		}
	}
	return intArray, nil
}

func ToJson(data interface{}) []byte {
	output, _ := json.Marshal(data)
	return output
}

func ToErrorJson(resultCode int, errText string) []byte {
	data := map[string]interface{}{
		"r":   resultCode,
		"err": errText,
	}
	return ToJson(data)
}
