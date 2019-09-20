package controllers

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/ESCah/go-against-humanity/app/game"
	"github.com/ESCah/go-against-humanity/app/models"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
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
	Name                      string
	NewCard                   models.Card
	Totals                    []Total
	Duration                  int
	Expires                   int64
	State                     models.MatchState
	InitialBlackCard          models.BlackCard
	SecondsUntilFinishPicking int
	WinnerUsername            string
	WinnerText                string
	Leaderboard               []models.Player
	Jury                      []models.Juror
}

type SocketServer struct {
	sync.Mutex
	mm    *game.MatchManager
	rooms map[int][]*websocket.Conn
	e     *echo.Echo
}

var upgrader = websocket.Upgrader{}

func MakeSocketServer(e *echo.Echo, mm *game.MatchManager) SocketServer {
	return SocketServer{sync.Mutex{}, mm, map[int][]*websocket.Conn{}, e}
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
				Name:             "join_successful",
				InitialBlackCard: card,
				State:            m.State,
				Expires:          expires,
			})
			s.Unlock()
		} else {
			s.Lock()
			conn.WriteJSON(Event{
				Name:  "join_successful",
				State: m.State,
			})
			s.Unlock()
		}
		s.rooms[matchID] = append(s.rooms[matchID], conn)
		return
	}
}

func (s *SocketServer) wsHandler(c echo.Context) error {
	matchId, err := strconv.Atoi(c.QueryParam("match"))
	if err != nil {
		log.Println("Failed to open ws: invalid match param")
		return c.String(http.StatusBadRequest, "Invalid 'match' parameter")
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("Failed to open ws: %s\n", err)
		return c.String(http.StatusBadRequest, "Could not open websocket connection")
	}
	// defer ws.Close()
	go s.onConnect(ws, matchId)

	return nil
}

func (s *SocketServer) Start() {
	s.e.GET("/ws", s.wsHandler)
}
