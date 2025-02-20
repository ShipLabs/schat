package base

import (
	"shiplabs/schat/internal/services"
)

func (b *base) WithAuthService() services.AuthServiceInterface {
	return services.NewAuthService(b.WithUserRepo())
}

func (b *base) WithPrivateChatService() services.ChatServiceInterface {
	return services.NewChatService(
		b.WithUserRepo(),
		b.WithPrivateChatRepo(),
		b.WithGroupRepo(),
		b.WithGroupMsgRepo(),
		b.WithPrivateMsgRepo(),
	)
}
