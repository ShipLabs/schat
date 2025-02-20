package base

import "shiplabs/schat/internal/handlers"

func (b *base) WithAuthController() handlers.AuthHandlerInterface {
	return handlers.NewAuthHandler(b.WithAuthService())
}
