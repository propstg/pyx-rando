package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	domain := os.Getenv("DOMAIN")
	secondsToWaitBeforeVoting, _ := strconv.Atoi(os.Getenv("SECONDS_TO_WAIT_BEFORE_VOTING"))
	gameId := os.Getenv("GAME_ID")

	networkService := &NetworkService{domain: domain, gameId: gameId}
	startGame(networkService, createRandomName(), secondsToWaitBeforeVoting)
}

func createRandomName() string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprintf("Rando%d", rand.Intn(1000))
}
