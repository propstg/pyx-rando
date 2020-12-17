package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// I've made a huge mistake.

func TestGame(t *testing.T) {
	networkService := FakeNetworkService{
		WaitForLongPollReturnValues: [][]LongPollEvent{
			// Hand deal to prime cardsInHand to test lobby clearing cards
			{
				{EventType: "hd",
					Hand: []WhiteCard{
						{CardId: 99},
					}},
			},
			// Switch to lobby to test lobby clearing cards
			{
				{EventType: "gsc", GameState: "l"},
			},
			// Deal cards (gsc first) and play first card from this hand (1) (proving lobby cleared cards)
			{
				{
					EventType: "gsc",
					BlackCard: BlackCard{NumberOfCardsToPick: 1},
					GameState: "p",
				},
				{EventType: "hd",
					Hand: []WhiteCard{
						{CardId: 1, WriteIn: false},
						{CardId: 2, WriteIn: false},
					},
				},
			},
			// Judging for first round
			{
				{EventType: "gsc", GameState: "j"},
			},
			// Deal cards (gsc last) and play first two cards from hand (2, 3)
			{
				{
					EventType: "hd",
					Hand: []WhiteCard{
						{CardId: 3, WriteIn: false},
					},
				}, {
					EventType: "gsc",
					BlackCard: BlackCard{NumberOfCardsToPick: 2},
					GameState: "p",
				},
			},
			// Judging for second round
			{
				{EventType: "gsc", GameState: "j"},
			},
			// Deal cards for third round
			{
				{
					EventType: "hd",
					Hand: []WhiteCard{
						{CardId: 4, WriteIn: false},
					},
				}, {
					EventType: "gsc",
					BlackCard: BlackCard{NumberOfCardsToPick: 1},
					GameState: "p",
				},
			},
			// Judging for third round, should pick first white card
			{
				{
					EventType: "gsc",
					GameState: "j",
					WhiteCards: [][]WhiteCard{
						{WhiteCard{CardId: 10}},
						{WhiteCard{CardId: 11}},
						{WhiteCard{CardId: 12}},
					},
				},
			},
			// Kicked event to exit
			{
				{EventType: "k"},
			},
		}, GetGameInfoReturnValues: []GameInfo{
			{PlayerInfo: []PlayerInfo{{Name: "fake name", Status: "sp"}}},  // Lobby gsc
			{PlayerInfo: []PlayerInfo{{Name: "fake name", Status: "sp"}}},  // First playing round gsc
			{PlayerInfo: []PlayerInfo{{Name: "fake name", Status: "sp"}}},  // First judging round gsc
			{PlayerInfo: []PlayerInfo{{Name: "fake name", Status: "sp"}}},  // Second playing round gsc
			{PlayerInfo: []PlayerInfo{{Name: "fake name", Status: "sp"}}},  // Second judging round gsc
			{PlayerInfo: []PlayerInfo{{Name: "fake name", Status: "sj"}}},  // Third playing round gsc
			{PlayerInfo: []PlayerInfo{{Name: "fake name", Status: "sjj"}}}, // Third judging round
		},
	}

	startGame(&networkService, "fake name", 0)

	assert.Equal(t, 1, networkService.LoginCalledTimes)
	assert.Equal(t, "fake name", networkService.LoginCalledWith[0])

	assert.Equal(t, 1, networkService.JoinCalledTimes)
	assert.Equal(t, 7, networkService.GetGameInfoCalledTimes)
	assert.Equal(t, 9, networkService.WaitForLongPollCalledTimes)

	assert.Equal(t, 3, networkService.PlayCardCalledTimes)
	assert.Equal(t, WhiteCard{CardId: 1, WriteIn: false}, networkService.PlayCardCalledWith[0])
	assert.Equal(t, WhiteCard{CardId: 2, WriteIn: false}, networkService.PlayCardCalledWith[1])
	assert.Equal(t, WhiteCard{CardId: 3, WriteIn: false}, networkService.PlayCardCalledWith[2])

	assert.Equal(t, 1, networkService.SelectWinningCardCalledTimes)
	assert.Equal(t, WhiteCard{CardId: 10}, networkService.SelectWinningCardCalledWith[0])
}

type FakeNetworkService struct {
	LoginCalledTimes int
	LoginCalledWith  []string

	JoinCalledTimes int

	WaitForLongPollCalledTimes  int
	WaitForLongPollReturnValues [][]LongPollEvent

	PlayCardCalledTimes int
	PlayCardCalledWith  []WhiteCard

	SelectWinningCardCalledTimes int
	SelectWinningCardCalledWith  []WhiteCard

	GetGameInfoCalledTimes  int
	GetGameInfoReturnValues []GameInfo
}

func (n *FakeNetworkService) Login(name string) {
	n.LoginCalledTimes++
	n.LoginCalledWith = append(n.LoginCalledWith, name)
}

func (n *FakeNetworkService) JoinGame() {
	n.JoinCalledTimes++
}

func (n *FakeNetworkService) WaitForLongPoll() []LongPollEvent {
	n.WaitForLongPollCalledTimes++
	return n.WaitForLongPollReturnValues[n.WaitForLongPollCalledTimes-1]
}

func (n *FakeNetworkService) PlayCard(card WhiteCard) {
	n.PlayCardCalledTimes++
	n.PlayCardCalledWith = append(n.PlayCardCalledWith, card)
}

func (n *FakeNetworkService) SelectWinningCard(card WhiteCard) {
	n.SelectWinningCardCalledTimes++
	n.SelectWinningCardCalledWith = append(n.SelectWinningCardCalledWith, card)
}

func (n *FakeNetworkService) GetGameInfo() GameInfo {
	n.GetGameInfoCalledTimes++
	return n.GetGameInfoReturnValues[n.GetGameInfoCalledTimes-1]
}
