package base

import repos "shiplabs/schat/internal/repositories"

func (b *base) WithUserRepo() repos.UserRepoInterface {
	return repos.NewUserRepo(*b.db)
}
