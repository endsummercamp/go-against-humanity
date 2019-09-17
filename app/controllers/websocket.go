package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ESCah/go-against-humanity/app/game"
	"github.com/ESCah/go-against-humanity/app/models"
	"github.com/gorilla/websocket"
	"sync"
)

type Card struct {
	ID   int
	Text string
}

type Total struct {
	ID    int
	Votes int
}

type Event struct {
	Name    string
	NewCard models.Card
	Totals  []Total
	Duration int
	Expires int64
	State models.MatchState
	InitialBlackCard models.BlackCard
	SecondsUntilFinishPicking int
	WinnerUsername string
	WinnerText string
}

type SocketServer struct {
	sync.Mutex
	mm    *game.MatchManager
	rooms map[int][]*websocket.Conn
}

func MakeSocketServer(mm *game.MatchManager) SocketServer {
	return SocketServer{sync.Mutex{}, mm, map[int][]*websocket.Conn{}}
}

func (s *SocketServer) BroadcastToRoom(room int, msg interface{}) {
	s.Lock()
	for _, conn := range s.rooms[room] {
		conn.WriteJSON(msg)
	}
	s.Unlock()
}

func (s *SocketServer) onConnect(conn *websocket.Conn, matchID int) {
	log.Printf("MatchID: %d\n", matchID)
	if !s.mm.IsJoinable(matchID) {
		s.Lock()
		conn.WriteJSON(Event{
			Name: "cannot_join",
		})
		s.Unlock()
		return
	} else {
		m := s.mm.GetMatchByID(matchID)
		round := m.GetRound()
		if round != nil {
			card := *round.BlackCard
			expires := round.Expires
			s.Lock()
			conn.WriteJSON(Event{
				Name: "join_successful",
				InitialBlackCard: card,
				State: m.State,
				Expires: expires,
			})
			s.Unlock()
		} else {
			s.Lock()
			conn.WriteJSON(Event{
				Name: "join_successful",
				State: m.State,
			})
			s.Unlock()
		}
		s.rooms[matchID] = append(s.rooms[matchID], conn)
		return
	}
}

func (s *SocketServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	matchID := r.URL.Query().Get("match")
	if matchID == "" {
		http.Error(w, "The 'match' parameter is required", http.StatusBadRequest)
		return
	}

	matchIDInt, err := strconv.Atoi(matchID)
	if err != nil {
		http.Error(w, "Invalid 'match' parameter", http.StatusBadRequest)
		return
	}

	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}
	go s.onConnect(conn, matchIDInt)
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
	return 0
}
