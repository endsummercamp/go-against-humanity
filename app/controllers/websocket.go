package controllers

import (
	"fmt"
	"github.com/ESCah/go-against-humanity/app/models"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Card struct {
	ID int
	Text string
}

type Total struct {
	ID int
	Votes int
}

type Event struct {
	Name string
	NewCard models.Card
	Totals []Total
}

type SocketServer struct {

}

func (s *SocketServer) onConnect(conn *websocket.Conn, matchID string) {
	m := Event{
		Name: "new_game",
	}
	// time.Sleep(time.Second)
	conn.WriteJSON(m)

	m = Event{
		Name: "new_black",
		NewCard: models.BlackCard{
			Deck: "",
			Icon: "",
			Text: "Lorem ipsum?",
			Pick: 1,
			Id: 1,
		},
	}
	conn.WriteJSON(m)

	var totals []Total

	for i := 0; i < 5; i++ {
		m = Event{
			Name: "new_white",
			NewCard: models.WhiteCard{
				Deck: "",
				Icon: "",
				Text: "Answer 1",
				Id: i,
			},
		}
		totals = append(totals, Total{
			ID: i,
			Votes: 0,
		})
		time.Sleep(time.Second)
		conn.WriteJSON(m)
	}

	for i := 0; i < 20; i++ {
		totals[rand.Intn(5)].Votes++
		time.Sleep(time.Second / 2)
		conn.WriteJSON(Event{
			Name: "totals",
			Totals: totals,
		})
	}
}

func (s *SocketServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	matchID := r.URL.Query().Get("match")
	if matchID == "" {
		http.Error(w, "The 'match' parameter is required", http.StatusBadRequest)
		return
	}
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	go s.onConnect(conn, matchID)
}

func (s *SocketServer) Start() int {
	http.HandleFunc("/ws", s.wsHandler)
	fmt.Println("Websocket server listening on :8080.")

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			panic(err)
		}
	}()
	return 0 // So that it can be used as "var _ = s.Start()"
}