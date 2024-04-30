package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Move struct {
	partOne string
	partTwo string
}

type Server struct {
	conns    map[string]*websocket.Conn
	upgrader websocket.Upgrader
	nextMove Move
}

func newServer() *Server {
	return &Server{
		conns: make(map[string]*websocket.Conn),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Accepting all requests
			},
		},

		nextMove: Move{
			partOne: "",
			partTwo: "",
		},
	}
}

// called whenever a new client connects
func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	// upgade http connection to websocket connection
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Connection upgrade error: ", err)
		return
	}

	fmt.Println("Client connected: ", conn.LocalAddr())
	// read the identification msg
	_, msg, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("read error", err)
	}
	// add new connection to server.conns map and print all connections
	s.conns[string(msg)] = conn
	fmt.Println("connections:")
	for dev, _ := range s.conns {
		fmt.Println(dev)
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte("ack"))
	if err != nil {
		fmt.Println("Write err", err)
	}

	s.readLoop(conn)
}

func (s *Server) msgDistributor(msg string, conn *websocket.Conn) {

	for len(s.getMissingConnections()) > 0 {
		time.Sleep(time.Second * 2)
	}

  dev := keysByValue(s.conns, conn)
	switch dev {

	case "engine":
		s.handleEngineMsg(msg, conn)

	case "psoc":
		s.handlePsocMsg(msg, conn)

	case "input":
		s.handleVoiceMsg(msg, conn)
	}
}

func (s *Server) handleVoiceMsg(msg string, conn *websocket.Conn) {

	if s.nextMove.partOne == "" {
		// get legal squares from engine
		s.nextMove.partOne = msg

		msg := "light " + s.nextMove.partOne
		err := s.conns["engine"].WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			fmt.Println("write error:", err)
		}
		conn.WriteMessage(websocket.TextMessage, []byte("input msg was sent to the engine"))

	} else {
		// get robotMove
		s.nextMove.partTwo = msg
		newMsg := "move " + s.nextMove.partOne + s.nextMove.partTwo
		err := s.conns["engine"].WriteMessage(websocket.TextMessage, []byte(newMsg))
		if err != nil {
			fmt.Println("write error:", err)
		}
		conn.WriteMessage(websocket.TextMessage, []byte("input msg was sent to the engine"))
	}

}

func (s *Server) handleEngineMsg(msg string, conn *websocket.Conn) {
  switch msg[:1]{

  case "l":
    s.nextMove.partOne = msg[1:]
    fmt.Println("got move part 1 from input device: ", s.nextMove.partOne)
    msg = "light " + s.nextMove.partOne
    err:= s.conns["psoc"].WriteMessage(websocket.TextMessage, []byte(msg))
    if err != nil {
			fmt.Println("write error:", err)
    }

  case "m":
    s.nextMove.partTwo = msg[1:]
    fmt.Println("got move part 2 from input device: ", s.nextMove.partTwo)
    msg = "move" + s.nextMove.partOne + s.nextMove.partTwo
    err:= s.conns["psoc"].WriteMessage(websocket.TextMessage, []byte(msg))
    if err != nil {
			fmt.Println("write error:", err)
    }
  }

}

func (s *Server) handlePsocMsg(msg string, conn *websocket.Conn) {
	fmt.Println("psoc message confirmed")
}

func (s *Server) readLoop(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read Error: ", err)
		}
		s.msgDistributor(string(msg), conn)
	}
}
func keysByValue(m map[string]*websocket.Conn, value *websocket.Conn) string {
    for k, v := range m {
        if value == v {
          return k
        }
    }
  return ""
  // add error return value
}


func (s *Server) connectionExists(key string) bool {
	if _, exists := s.conns[key]; exists {
		return true
	}
	return false
}

func (s *Server) getMissingConnections() []string {
	var devices [3]string
	devices[0] = "engine"
	devices[1] = "psoc"
	devices[2] = "input"

	var notConnected []string

	for _, dev := range devices {
		if s.connectionExists(dev) == false {
			notConnected = append(notConnected, dev)
		}
	}

	for _, dev := range notConnected {
		fmt.Println(dev, " not connected")
	}
	return notConnected
}








func main() {
	server := newServer()
	http.HandleFunc("/ws", server.wsHandler)
	fmt.Println("Starting server on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
