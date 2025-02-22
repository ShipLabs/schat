package handlers

import (
	"log"
	"net/http"
	"shiplabs/schat/internal/pkg/store"
	"shiplabs/schat/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WSResponse struct {
	StatusCode   int    `json:"status_code"`
	ErrorMessage string `json:"error_msg"`
	Data         string `json:"data"`
}

type wsHandler struct {
	store              store.ConnectionStoreInterface
	privateChatService services.ChatServiceInterface
}

type WsHandlerInterface interface {
	HandlePrivateChat(ctx *gin.Context)
	HandleGroupChat(ctx *gin.Context)
	HandlerGroupCreation(ctx *gin.Context)
}

func NewWebSocketHandler(
	store store.ConnectionStoreInterface,
	pChatService services.ChatServiceInterface,
) WsHandlerInterface {
	return &wsHandler{
		store:              store,
		privateChatService: pChatService,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true //dev purposes only
	},
}

func (w *wsHandler) handleResponse(conn *websocket.Conn, code int, msg, data string) {
	resp := WSResponse{
		StatusCode:   code,
		ErrorMessage: msg,
		Data:         data,
	}
	if err := conn.WriteJSON(&resp); err != nil {
		log.Println(err)
	}
}

func (w *wsHandler) HandlePrivateChat(ctx *gin.Context) {
	userID := uuid.MustParse(ctx.GetString("userID"))
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		w.handleResponse(conn, http.StatusBadRequest, err.Error(), "")
	}

	eConn, err := w.store.GetConn(userID)
	if err != nil {
		w.store.SaveConn(userID, conn)
		eConn = conn
	}

	var message services.PrivateMessageDto
	if err := eConn.ReadJSON(&message); err != nil {
		//check unepected connection closure error
		log.Println(err)
		w.handleResponse(conn, http.StatusBadRequest, err.Error(), "")
	}

	if err := w.privateChatService.SendPrivateMsg(message); err != nil {
		log.Println(err)
		w.handleResponse(conn, http.StatusBadRequest, err.Error(), "")
	}

	rconn, err := w.store.GetConn(message.ReceiverID)
	if err != nil {
		log.Println("receiver is not online")
		return
	}

	w.handleResponse(rconn, http.StatusOK, "", message.Content)
}

func (w *wsHandler) HandleGroupChat(ctx *gin.Context) {
	_, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		//handle
	}
}
func (w *wsHandler) HandlerGroupCreation(ctx *gin.Context) {
	_, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		//handle
	}
}
