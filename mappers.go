package main

import (
	"encoding/json"
	"strings"
)

func MapLongPollData(jsonPayload string) []LongPollEvent {
	var result []LongPollEvent
	if !strings.HasPrefix(jsonPayload, "[") {
		jsonPayload = "[" + jsonPayload + "]"
	}
	json.Unmarshal([]byte(jsonPayload), &result)
	return result
}

func MapGetGameInfoData(jsonPayload string) GameInfo {
	var result GameInfo
	json.Unmarshal([]byte(jsonPayload), &result)
	return result
}
