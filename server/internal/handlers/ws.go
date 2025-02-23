package handlers

import (
	"errors"
	"log"
	"net/http"
	"shiplabs/schat/internal/pkg/store"
	"shiplabs/schat/internal/services"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	ErrHandShakeFail = errors.New("failed handshake, connection not established")
)

type WSResponse struct {
	StatusCode   int    `json:"status_code"`
	ErrorMessage string `json:"error_msg"`
	Data         string `json:"data"`
}

type wsHandler struct {
	store              store.ConnectionStoreInterface
	privateChatService services.ChatServiceInterface
	groupService       services.GroupServiceInterface
}

type WsHandlerInterface interface {
	Connect(userID uuid.UUID, ctx *gin.Context) (*websocket.Conn, error)
	HandlePrivateChat(ctx *gin.Context)
	GroupCreationHandler(ctx *gin.Context)
	// HandleGroupChat(ctx *gin.Context)
	// HandlerGroupCreation(ctx *gin.Context)
}

func NewWebSocketHandler(
	store store.ConnectionStoreInterface,
	pChatService services.ChatServiceInterface,
	groupService services.GroupServiceInterface,
) WsHandlerInterface {
	return &wsHandler{
		store:              store,
		privateChatService: pChatService,
		groupService:       groupService,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true //dev purposes only
	},
}

func (w *wsHandler) Connect(userID uuid.UUID, ctx *gin.Context) (*websocket.Conn, error) {
	//TOD0: how do I manage connections better (at scale or not???)
	eConn, err := w.store.GetConn(userID)
	if err == nil {
		return eConn, nil
	}
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return nil, ErrHandShakeFail
	}

	w.store.SaveConn(userID, conn)

	return conn, nil
}

func (w *wsHandler) HandlePrivateChat(ctx *gin.Context) {
	userID := uuid.MustParse(ctx.GetString("userID"))
	conn, err := w.Connect(userID, ctx)
	if err != nil {
		return //how do I better handle error here?? send back http json seems best since ws conn has not been established
	}
	defer w.closeConn(conn, userID) //should I close the connection here???

	for {
		var message services.PrivateMessageDto
		if err := conn.ReadJSON(&message); err != nil {
			//TODO: check unepected connection closure error
			log.Println(err)
			w.handleResponse(conn, http.StatusBadRequest, err.Error(), "")
		}

		go w.handlerIncomingPrivateMsg(&message, conn)
	}
}

func (w *wsHandler) HandleGroupChat(ctx *gin.Context) {
	_, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		//handle
	}
}

func (w *wsHandler) GroupCreationHandler(ctx *gin.Context) {
	userID := uuid.MustParse(ctx.GetString("userID"))
	conn, err := w.Connect(userID, ctx)
	if err != nil {
		return //how do I better handle error here?? send back http json seems best since ws conn has not been established
	}
	defer w.closeConn(conn, userID)

	for {
		var message services.CreateGroupDto
		if err := conn.ReadJSON(&message); err != nil {
			//TOD: check unepected connection closure error
			log.Println(err)
			w.handleResponse(conn, http.StatusBadRequest, err.Error(), "")
		}

		go w.handleGroupCreation(userID, &message, conn)
	}
}

func (w *wsHandler) handleGroupCreation(userID uuid.UUID, data *services.CreateGroupDto, createrConn *websocket.Conn) {
	if err := w.groupService.CreateGroup(userID, *data); err != nil {
		log.Println(err)
		w.handleResponse(createrConn, http.StatusBadRequest, err.Error(), "")
	}

	var wg sync.WaitGroup
	for _, memberID := range data.Members {
		wg.Add(1)
		go w.notifyOnlineUser(memberID, data, &wg)
	}

	wg.Wait()
}

func (w *wsHandler) notifyOnlineUser(userID uuid.UUID, data *services.CreateGroupDto, wg *sync.WaitGroup) {
	conn, err := w.store.GetConn(userID)
	if err != nil {
		log.Println("not online")
		wg.Done()
		return
	}

	msg := "you have been added to "
	w.handleResponse(conn, http.StatusOK, "", msg+data.GroupName)
	wg.Done()
}

func (w *wsHandler) handlerIncomingPrivateMsg(message *services.PrivateMessageDto, senderConn *websocket.Conn) {
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
