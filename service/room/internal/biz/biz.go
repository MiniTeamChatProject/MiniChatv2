package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewRoomUsecase, NewRoomMemberUsecase, NewMessageUsecase)
