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

	defer w.closeConn(conn, userID)

	for {
		var message services.PrivateMessageDto
		if err := eConn.ReadJSON(&message); err != nil {
			//TOD: check unepected connection closure error
			log.Println(err)
			w.handleResponse(conn, http.StatusBadRequest, err.Error(), "")
		}

		go w.handlerIncomingMsg(&message, conn)
	}
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

func (w *wsHandler) handlerIncomingMsg(message *services.PrivateMessageDto, senderConn *websocket.Conn) {
	if err := w.privateChatService.SendPrivateMsg(*message); err != nil {
		log.Println(err)
		w.handleResponse(senderConn, http.StatusBadRequest, err.Error(), "")
	}

	rconn, err := w.store.GetConn(uuid.MustParse(message.ReceiverID))
	if err != nil {
		log.Println("receiver is not online")
		return
	}

	w.handleResponse(rconn, http.StatusOK, "", message.Content)
}

func (w *wsHandler) handleResponse(conn *websocket.Conn, code int, msg, data string) {
	resp := WSResponse{
		StatusCode:   code,
		ErrorMessage: msg,
		Data:         data,
	}

	if conn == nil {
		log.Println("error sending resp: ", resp)
		return
	}

	if err := conn.WriteJSON(&resp); err != nil {
		log.Println(err)
	}
}

func (w *wsHandler) closeConn(conn *websocket.Conn, userID uuid.UUID) {
	if err := conn.Close(); err != nil {
		log.Println("failed to close connection for user with id", userID)
	}
	w.store.DeleteConn(userID)
	log.Println("connection closed for user with id", userID)
}
