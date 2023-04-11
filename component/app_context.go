package component

import (
	"go-mail/pubsub"
)

type AppContext interface {
	GetPubsub() pubsub.Pubsub
}

type appCtx struct {
	pb        pubsub.Pubsub
}

func NewAppContext(pb pubsub.Pubsub) *appCtx {
	return &appCtx{pb: pb}
}

func (ctx *appCtx) GetPubsub() pubsub.Pubsub { return ctx.pb }
