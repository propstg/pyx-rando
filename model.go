package main

type LongPollEvent struct {
	EventType  string        `json:"E"`
	Hand       []WhiteCard   `json:"h"`
	BlackCard  BlackCard     `json:"bc"`
	GameState  string        `json:"gs"`
	WhiteCards [][]WhiteCard `json:"wc"`
}

type BlackCard struct {
	NumberOfCardsToPick int `json:"PK"`
}

type WhiteCard struct {
	CardId  int  `json:"cid"`
	WriteIn bool `json:"wi"`
}

type GameInfo struct {
	PlayerInfo []struct {
		Status string `json:"st"`
		Name   string `json:"N"`
		Score  int    `json:"sc"`
	} `json:"pi"`
}
