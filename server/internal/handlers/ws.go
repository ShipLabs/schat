package handlers

import (
	"errors"
	"log"
	"net/http"
	"shiplabs/schat/internal/pkg/store"
	"shiplabs/schat/internal/services"
	"shiplabs/schat/pkg/shared"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	ErrHandShakeFail   = errors.New("failed handshake, connection not established")
	ErrInvalidGroup    = errors.New("invalid group")
	ErrInvalidMemberID = errors.New("invalid member id")
)

type WSResponse struct {
	StatusCode   int    `json:"status_code"`
	ErrorMessage string `json:"error_msg"`
	Data         string `json:"data"`
}

type wsHandler struct {
	store        store.ConnectionStoreInterface
	chatService  services.ChatServiceInterface
	groupService services.GroupServiceInterface
}

type WsHandlerInterface interface {
	Connect(userID uuid.UUID, ctx *gin.Context) (*websocket.Conn, error)
	HandlePrivateChat(ctx *gin.Context)
	GroupCreationHandler(ctx *gin.Context)
	HandleGroupChat(ctx *gin.Context)
	HandleMembership(ctx *gin.Context)
}

func NewWebSocketHandler(
	store store.ConnectionStoreInterface,
	pChatService services.ChatServiceInterface,
	groupService services.GroupServiceInterface,
) WsHandlerInterface {
	return &wsHandler{
		store:        store,
		chatService:  pChatService,
		groupService: groupService,
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
		shared.ErrorResponse(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	defer w.closeConn(conn, userID)

	for {
		var message services.PrivateMessageDto
		if err := conn.ReadJSON(&message); err != nil {
			//TODO: check unepected connection closure error
			log.Println(err)
			w.handleResponse(conn, http.StatusBadRequest, err, "")
		}

		go w.handlerIncomingPrivateMsg(&message, conn)
	}
}

func (w *wsHandler) HandleGroupChat(ctx *gin.Context) {
	userID := uuid.MustParse(ctx.GetString("userID"))
	conn, err := w.Connect(userID, ctx)
	if err != nil {
		shared.ErrorResponse(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	defer w.closeConn(conn, userID)

	for {
		var message services.GroupMessageDto
		if err := conn.ReadJSON(&message); err != nil {
			//TOD: check unepected connection closure error
			log.Println(err)
			w.handleResponse(conn, http.StatusBadRequest, err, "")
		}

		go w.handleGroupMessage(&message, conn)
	}
}

func (w *wsHandler) GroupCreationHandler(ctx *gin.Context) {
	userID := uuid.MustParse(ctx.GetString("userID"))
	conn, err := w.Connect(userID, ctx)
	if err != nil {
		shared.ErrorResponse(ctx, http.StatusUnprocessableEntity, err.Error())
		return
	}
	defer w.closeConn(conn, userID)

	for {
		var message services.CreateGroupDto
		if err := conn.ReadJSON(&message); err != nil {
			//TOD: check unepected connection closure error
			log.Println(err)
			w.handleResponse(conn, http.StatusBadRequest, err, "")
		}

		go w.handleGroupCreation(userID, &message, conn)
	}
}

func (w *wsHandler) HandleMembership(ctx *gin.Context) {
	userID := uuid.MustParse(ctx.GetString("userID"))
	groupID, err := uuid.Parse(ctx.Param("group_id"))
	if err != nil {
		shared.ErrorResponse(ctx, http.StatusUnprocessableEntity, "invalid group id")
		return
	}
	conn, err := w.Connect(userID, ctx)
	if err != nil {
		return
	}
	defer w.closeConn(conn, userID)

	for {
		var message services.GroupMembershipDto
		if err := conn.ReadJSON(&message); err != nil {
			//TOD: check unepected connection closure error
			log.Println(err)
			w.handleResponse(conn, http.StatusBadRequest, err, "")
		}

		go w.handleGroupMemberShip(userID, groupID, &message, conn)
	}
}
func (w *wsHandler) handleGroupMemberShip(userID, groupID uuid.UUID, data *services.GroupMembershipDto, createrConn *websocket.Conn) {
	memberID, err := uuid.Parse(data.MemberID)
	if err != nil {
		log.Println(err)
		w.handleResponse(createrConn, http.StatusBadRequest, ErrInvalidMemberID, "")
		return
	}
	if err := w.groupService.HandleMembership(groupID, userID, memberID, data.Action); err != nil {
		log.Println(err)
		w.handleResponse(createrConn, http.StatusBadRequest, err, "")
		return
	}

	members, err := w.groupService.GetGroupMembers(groupID)
	if err != nil {
		log.Println(err)
		//group has been created but members could not be fetched for notification, how to best handle???
		w.handleResponse(createrConn, http.StatusBadRequest, err, "")
		return
	}

	var wg sync.WaitGroup
	for _, member := range members {
		wg.Add(1)
		go w.groupMembershipNotification(member.UserID, data)
	}
}

func (w *wsHandler) handleGroupCreation(userID uuid.UUID, data *services.CreateGroupDto, createrConn *websocket.Conn) {
	if err := w.groupService.CreateGroup(userID, *data); err != nil {
		log.Println(err)
		w.handleResponse(createrConn, http.StatusBadRequest, err, "")
	}

	for _, memberID := range data.Members {
		//handle this better - find a way
		memberUUID, err := uuid.Parse(memberID)
		if err != nil {
			log.Println("invalid member UUID:", memberID)
			continue
		}
		go w.groupCreationNotification(memberUUID, data)
	}
}

func (w *wsHandler) handleGroupMessage(data *services.GroupMessageDto, senderConn *websocket.Conn) {
	groupUUID, err := uuid.Parse(data.GroupID)
	if err != nil {
		log.Println(err)
		w.handleResponse(senderConn, http.StatusBadRequest, err, "")
		return
	}
	if err := w.chatService.SendMsgToGroup(*data); err != nil {
		log.Println(err)
		w.handleResponse(senderConn, http.StatusBadRequest, err, "")
		return
	}
	members, err := w.groupService.GetGroupMembers(groupUUID)
	if err != nil {
		log.Println(err)
		w.handleResponse(senderConn, http.StatusBadRequest, err, "")
		return
	}

	for _, member := range members {
		go w.groupMessageNotification(member.UserID, data)
	}
}

func (w *wsHandler) groupMessageNotification(userID uuid.UUID, data *services.GroupMessageDto) {
	conn, err := w.store.GetConn(userID)
	if err != nil {
		log.Println("not online")
		return
	}
	w.handleResponse(conn, http.StatusOK, nil, data.Content)
}

func (w *wsHandler) groupMembershipNotification(userID uuid.UUID, data *services.GroupMembershipDto) {
	conn, err := w.store.GetConn(userID)
	if err != nil {
		log.Println("not online")

		return
	}
	var action string
	if data.Action == services.Add {
		action += "ed"
	} else {
		action += "d"
	}

	w.handleResponse(conn, http.StatusOK, nil, data.MemberID+" has been "+action)
}

func (w *wsHandler) groupCreationNotification(userID uuid.UUID, data *services.CreateGroupDto) {
	conn, err := w.store.GetConn(userID)
	if err != nil {
		log.Println("not online")
		return
	}
	msg := "you have been added to "
	w.handleResponse(conn, http.StatusOK, nil, msg+data.GroupName)
}

func (w *wsHandler) handlerIncomingPrivateMsg(message *services.PrivateMessageDto, senderConn *websocket.Conn) {
	if err := w.chatService.SendPrivateMsg(*message); err != nil {
		w.handleResponse(senderConn, http.StatusBadRequest, err, "")
	}

	rconn, err := w.store.GetConn(uuid.MustParse(message.ReceiverID))
	if err != nil {
		log.Println("receiver is not online")
		return
	}

	w.handleResponse(rconn, http.StatusOK, nil, message.Content)
}

func (w *wsHandler) handleResponse(conn *websocket.Conn, code int, err error, data string) {
	if err != nil {
		log.Println(err.Error())
	}
	resp := WSResponse{
		StatusCode:   code,
		ErrorMessage: err.Error(),
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
