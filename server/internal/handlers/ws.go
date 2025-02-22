package handlers

import (
	"log"
	"net/http"
	"shiplabs/schat/internal/pkg/store"
	"shiplabs/schat/internal/services"
	"shiplabs/schat/pkg/shared"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	store              store.ConnectionStoreInterface
	privateChatService services.ChatServiceInterface
}

type wsHandlerInterface interface {
	handlePrivateChat(ctx *gin.Context)
	handleGroupChat(ctx *gin.Context)
	handlerGroupCreation(ctx *gin.Context)
}

func NewWebSocketHandler(
	store store.ConnectionStoreInterface,
	pChatService services.ChatServiceInterface,
) wsHandlerInterface {
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

func (w *wsHandler) handlePrivateChat(ctx *gin.Context) {
	userID := uuid.MustParse(ctx.GetString("userID"))
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		shared.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
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
		shared.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if err := w.privateChatService.SendPrivateMsg(message); err != nil {
		log.Println(err)
		shared.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	rconn, err := w.store.GetConn(message.ReceiverID)
	if err != nil {
		log.Println(err)
		shared.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	//assume text message is default
	if err := rconn.WriteMessage(websocket.TextMessage, []byte(message.Content)); err != nil {
		log.Println(err)
		shared.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
}

func (w *wsHandler) handleGroupChat(ctx *gin.Context) {
	_, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		//handle
	}
}
func (w *wsHandler) handlerGroupCreation(ctx *gin.Context) {
	_, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		//handle
	}
}
