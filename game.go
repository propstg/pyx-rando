package main

import (
	"log"
	"sort"
	"time"
)

const (
	EventHandDeal = "hd"
	EventGameStateChange = "gsc"
	GameStatePlaying = "p"
	GameStateJudging = "j"
	GameStateLobby = "l"
	PlayerStatusJudge = "sj"
	PlayerStatusJudging = "sjj"
)

func startGame(networkService INetworkService, name string, secondsToWaitBeforeVoting int) {
	var cardsInHand []WhiteCard
	areWeCzar := false

	networkService.Login(name)
	networkService.JoinGame()

	for {
		events := networkService.WaitForLongPoll()
		log.Printf("Received %d events...\n", len(events))
		sortEventsSoHandDealIsBeforeGameStateChange(events)

		for _, event := range events {
			log.Println(event.EventType)
			switch event.EventType {
			case EventHandDeal:
				for _, card := range event.Hand {
					cardsInHand = append(cardsInHand, card)
				}
				log.Printf("Added %d cards to hand.\n", len(event.Hand))
			case EventGameStateChange:
				gameInfo := networkService.GetGameInfo()
				areWeCzar = areWeTheCzar(gameInfo, name)
				log.Printf("New game state: %s\n", event.GameState)
				log.Printf("Are we the czar? %t\n", areWeCzar)
				if event.GameState == GameStatePlaying && !areWeCzar {
					for i := 0; i < event.BlackCard.NumberOfCardsToPick; i++ {
						var chosenCard WhiteCard
						chosenCard, cardsInHand = drawCard(cardsInHand)
						log.Printf("Attempting to play card: %d\n", chosenCard.CardId)
						networkService.PlayCard(chosenCard)
					}
				} else if event.GameState == GameStateJudging && areWeCzar {
					time.Sleep(time.Duration(secondsToWaitBeforeVoting) * time.Second)
					log.Printf("Attempting to select winning card: %d\n", event.WhiteCards[0][0].CardId)
					networkService.SelectWinningCard(event.WhiteCards[0][0])
				} else if event.GameState == GameStateLobby {
					cardsInHand = []WhiteCard{}
				}
			}
		}
	}
}

func sortEventsSoHandDealIsBeforeGameStateChange(events []LongPollEvent) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].GameState < events[j].GameState
	})
}

func areWeTheCzar(gameInfo GameInfo, name string) bool {
	for _, info := range gameInfo.PlayerInfo {
		if name == info.Name {
			return info.Status == PlayerStatusJudge || info.Status == PlayerStatusJudging
		}
	}
	return false
}

func drawCard(cards []WhiteCard) (WhiteCard, []WhiteCard) {
	return cards[0], cards[1:]
}
