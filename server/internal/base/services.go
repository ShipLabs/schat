package base

import (
	repos "shiplabs/schat/internal/repositories"
	"shiplabs/schat/internal/services"
)

func (b *base) WithAuthService(userRepo repos.UserRepoInterface) services.AuthServiceInterface {
	return services.NewAuthService(b.WithUserRepo())
}

func (b *base) WithPrivateChatService(
	userRepo repos.UserRepoInterface,
	privateChatRepo repos.PrivateChatRepoInterface,
	groupRepo repos.GroupRepoInterface,
	groupMsgRepo repos.GroupMessageRepoInterface,
	privateMsgRepo repos.PrivateMessageRepoInterface,
) services.ChatServiceInterface {
	return services.NewChatService(
		b.WithUserRepo(),
		b.WithPrivateChatRepo(),
		b.WithGroupRepo(),
		b.WithGroupMsgRepo(),
		b.WithPrivateMsgRepo(),
	)
}
