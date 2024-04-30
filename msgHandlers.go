package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)


func (s *Server) msgDistributor(msg string, conn *websocket.Conn) {

	switch msg[:1] {

	case "e":
		s.handleEngineMsg(msg[1:], conn)

	case "p":
		s.handlePsocMsg(msg[1:], conn)

  case "i":
    s.handleVoiceInput(msg[1:], conn)
	}

}

func (s *Server) handleVoiceInput(msg string, conn *websocket.Conn) {

	if s.nextMove.partOne == "" {
		// get legal squares from engine
		s.nextMove.partOne = msg
		msg := "l" + s.nextMove.partOne
		err := s.conns["engine"].WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			fmt.Println("write error:", err)
		}
		conn.WriteMessage(websocket.TextMessage, []byte("msg passed on to engine"))
	} else {
		// get robotMove
		s.nextMove.partTwo = msg
		newMsg := "m" + s.nextMove.partOne + s.nextMove.partTwo
		err := s.conns["engine"].WriteMessage(websocket.TextMessage, []byte(newMsg))
		if err != nil {
			fmt.Println("write error:", err)
		}
		conn.WriteMessage(websocket.TextMessage, []byte("msg passed on to engine"))
		s.nextMove.partOne = ""
		s.nextMove.partTwo = ""
	}

}

func (s *Server) handleEngineMsg(msg string, conn *websocket.Conn) {
	fmt.Println("handled Engine Message: ", msg)
	confirmMsg := []byte("engine msg handled by server")
	err := conn.WriteMessage(websocket.TextMessage, confirmMsg)
	if err != nil {
		fmt.Println("Write error: ", err)
	}
}

func (s *Server) handlePsocMsg(msg string, conn *websocket.Conn) {
	fmt.Println("handled PSoC Message:", msg)
	confirmMsg := []byte("PSoC msg handled by server")
	err := conn.WriteMessage(websocket.TextMessage, confirmMsg)
	if err != nil {
		fmt.Println("Write error: ", err)
	}
}
