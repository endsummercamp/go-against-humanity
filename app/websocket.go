package app

import (
	"fmt"
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
	NewCard Card
	Totals []Total
}

func echo(conn *websocket.Conn) {
	m := Event{
		Name: "new_game",
	}
	// time.Sleep(time.Second)
	conn.WriteJSON(m)

	m = Event{
		Name: "new_black",
		NewCard: Card{
			Text: "Lorem ipsum?",
		},
	}
	time.Sleep(time.Second)
	conn.WriteJSON(m)

	var totals []Total

	for i := 0; i < 5; i++ {
		m = Event{
			Name: "new_white",
			NewCard: Card{
				ID: i,
				Text: "Lorem ipsum.",
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

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	go echo(conn)
}

func wsMain() {
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("Websocket server listening on :8080.")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}