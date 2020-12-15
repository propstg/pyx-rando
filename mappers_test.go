package main

import "testing"
import "github.com/stretchr/testify/assert"

func TestMapLongPollData_multiple(t *testing.T) {
	jsonPayload := `[{
		"h": [
			{"cid": 1, "wi": false},
			{"cid": 457, "wi": false},
			{"cid": 160, "wi": false},
			{"cid": 416, "wi": false},
			{"cid": 144, "wi": false},
			{"cid": 322, "wi": false},
			{"cid": 13, "wi": false},
			{"cid": 178, "wi": false},
			{"cid": 66, "wi": false},
			{ "cid": 410, "wi": false }
		],
		"gid": 0,
		"E": "hd"
	},{
		"Pt": 2147483647,
		"bc": {
			"W": "US",
			"cid": 94,
			"D": 0,
			"PK": 1,
			"T": ""
		},
		"E": "gsc",
		"ts": 1607886037548,
		"gid": 0,
		"gs": "p"
	},{
		"E":"gsc",
		"wc":[
			[{ "cid":159,"wi":false }],
			[{ "cid":143, "wi":false }]
		],
		"gs":"j"
	}]`

	result := MapLongPollData(jsonPayload)

	assert.Equal(t, 3, len(result))

	assert.Equal(t, "hd", result[0].EventType)
	assert.Equal(t, 10, len(result[0].Hand))
	assert.Equal(t, 1, result[0].Hand[0].CardId)
	assert.Equal(t, false, result[0].Hand[0].WriteIn)
	assert.Equal(t, 410, result[0].Hand[9].CardId)
	assert.Equal(t, false, result[0].Hand[9].WriteIn)

	assert.Equal(t, "gsc", result[1].EventType)
	assert.Equal(t, 1, result[1].BlackCard.NumberOfCardsToPick)
	assert.Equal(t, "p", result[1].GameState)

	assert.Equal(t, "gsc", result[2].EventType)
	assert.Equal(t, "j", result[2].GameState)
	assert.Equal(t, 2, len(result[2].WhiteCards))
	assert.Equal(t, 1, len(result[2].WhiteCards[0]))
	assert.Equal(t, 1, len(result[2].WhiteCards[1]))
	assert.Equal(t, 159, result[2].WhiteCards[0][0].CardId)
	assert.Equal(t, 143, result[2].WhiteCards[1][0].CardId)
}

func TestMapLongPollData_single(t *testing.T) {
	jsonPayload := `{"ts":1607887600741,"E":"_"}`

	result := MapLongPollData(jsonPayload)

	assert.Equal(t, 1, len(result))
	assert.Equal(t, "_", result[0].EventType)
}

func TestMapGetGameInfoData(t *testing.T) {
	jsonPayload := `{
		"pi": [
			{
				"st": "sp",
				"N": "Player 1"
			},
			{
				"st": "sp",
				"N": "Player 2"
			},
			{
				"st": "sj",
				"N": "Rando719"
			}
		]
	}`

	result := MapGetGameInfoData(jsonPayload)

	assert.Equal(t, 3, len(result.PlayerInfo))

	assert.Equal(t, "sp", result.PlayerInfo[0].Status)
	assert.Equal(t, "Player 1", result.PlayerInfo[0].Name)

	assert.Equal(t, "sp", result.PlayerInfo[1].Status)
	assert.Equal(t, "Player 2", result.PlayerInfo[1].Name)

	assert.Equal(t, "sj", result.PlayerInfo[2].Status)
	assert.Equal(t, "Rando719", result.PlayerInfo[2].Name)
}