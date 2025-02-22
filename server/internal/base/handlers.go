package base

import "shiplabs/schat/internal/handlers"

func (b *base) WithAuthController() handlers.AuthHandlerInterface {
	return handlers.NewAuthHandler(b.WithAuthService())
}

func (b *base) WithChatController() handlers.WsHandlerInterface {
	return handlers.NewWebSocketHandler(b.wsStore, b.WithPrivateChatService())
}
