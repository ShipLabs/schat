package handlers

import (
	"net/http"
	"shiplabs/schat/internal/pkg/store"

	"github.com/gorilla/websocket"
)

//route - middleware - ws handler

//handler upgrades
// is connection existing? check store (how do I manage connections in say a chat app)
//any ws headers to pay attention to??
//get message (data, should be json)
//call appropriate handler depending on endpoint called
//how do I handle error???
//broadcast message to receipient, how do I send feedback to sender? like tick for delivered message??

type wsHandler struct {
	store store.ConnectionStoreInterface
}

type wsHandlerInterface interface {
}

func NewWebSocketHandler(store store.ConnectionStoreInterface) *wsHandler {
	return &wsHandler{
		store: store,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true //dev purposes onl
	},
}

func (w *wsHandler) handlePrivateChat()    {}
func (w *wsHandler) handleGroupChat()      {}
func (w *wsHandler) handlerGroupCreation() {}
