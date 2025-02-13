package base

import repos "shiplabs/schat/internal/repositories"

func (b *base) WithUserRepo() repos.UserRepoInterface {
	return repos.NewUserRepo(*b.db)
}

func (b *base) WithGroupMsgRepo() repos.GroupMessageRepoInterface {
	return repos.NewGroupMessageRepo(*b.db)
}

func (b *base) WithPrivateMsgRepo() repos.PrivateMessageRepoInterface {
	return repos.NewPrivateMessageRepo(*b.db)
}

func (b *base) WithGroupRepo() repos.GroupRepoInterface {
	return repos.NewGroupRepo(*b.db)
}

func (b *base) WithPrivateChatRepo() repos.PrivateChatRepoInterface {
	return repos.NewPrivateChatRepo(*b.db)
}
