package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type INetworkService interface {
	Login(name string)
	JoinGame()
	WaitForLongPoll() []LongPollEvent
	PlayCard(card WhiteCard)
	SelectWinningCard(card WhiteCard)
	GetGameInfo() GameInfo
}

type NetworkService struct {
	domain string
	gameId string
	cookie string
}

func (n *NetworkService) Login(name string) {
	_, res := n.doGenericAjaxCall(fmt.Sprintf("o=r&n=%s&s=1", name))
	n.cookie = res.Header.Get("Set-Cookie")
}

func (n *NetworkService) JoinGame() {
	n.doGenericAjaxCall(fmt.Sprintf("o=jg&gid=%s&pw=&s=12", n.gameId))
}

func (n *NetworkService) WaitForLongPoll() []LongPollEvent {
	body, _ := n.doLongPollCall()
	return MapLongPollData(string(body))
}

func (n *NetworkService) PlayCard(card WhiteCard) {
	n.doGenericAjaxCall(fmt.Sprintf("o=pc&gid=%s&cid=%d&s=22", n.gameId, card.CardId))
}

func (n* NetworkService) SelectWinningCard(card WhiteCard) {
	n.doGenericAjaxCall(fmt.Sprintf("o=js&gid=%s&cid=%d&s=42", n.gameId, card.CardId))
}

func (n* NetworkService) GetGameInfo() GameInfo {
	body, _ := n.doGenericAjaxCall(fmt.Sprintf("o=ggi&gid=%s&s=63", n.gameId))
	return MapGetGameInfoData(string(body))
}

func (n* NetworkService) doGenericAjaxCall(payload string) ([]byte, *http.Response) {
	return n.doGenericServletCall(payload, "/AjaxServlet")
}

func (n* NetworkService) doLongPollCall() ([]byte, *http.Response) {
	return n.doGenericServletCall("", "/LongPollServlet")
}

func (n* NetworkService) doGenericServletCall(payload string, endpoint string) ([]byte, *http.Response) {
	fmt.Printf("Sending %s to %s...\n", payload, endpoint)
	method := "POST"
	payloadWrapper := strings.NewReader(payload)

	client := &http.Client {
	}
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", n.domain, endpoint), payloadWrapper)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Sec-Fetch-Site", "same-origin")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Referer", fmt.Sprintf("%s/game.jsp", n.domain))
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Cookie", n.cookie)

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))

	return body, res
}
