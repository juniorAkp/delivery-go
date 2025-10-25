package utils

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Coord struct {
	Longitude float64 `json:"longitude" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
}

type ConnectionHub struct {
	activeCustomers map[string]*UserConnection
	activeCouriers  map[string]*UserConnection
	messageQueue    chan *OutgoingMessage
	registerConn    chan *UserConnection
	removeConn      chan *UserConnection
	mutex           sync.RWMutex
}

type UserConnection struct {
	UserID          string
	AccountType     string
	Socket          *websocket.Conn
	OutgoingMessage chan []byte
	Manager         *ConnectionHub
}

type OutgoingMessage struct {
	MessageType string    `json:"type"`
	SenderID    string    `json:"from"`
	ReceiverID  string    `json:"to,omitempty"`
	Content     any       `json:"data"`
	SentAt      time.Time `json:"timeStamp"`
}
