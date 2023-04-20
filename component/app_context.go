package component

import (
	"go-mail/memcache"
	"go-mail/pubsub"
)

type AppContext interface {
	GetPubsub() pubsub.Pubsub
	GetCaching() memcache.Caching
}

type appCtx struct {
	pb        pubsub.Pubsub
	cache	  memcache.Caching
}

func NewAppContext(pb pubsub.Pubsub, cache memcache.Caching) *appCtx {
	return &appCtx{pb: pb, cache: cache}
}

func (ctx *appCtx) GetPubsub() pubsub.Pubsub { return ctx.pb }
func(ctx *appCtx) GetCaching() memcache.Caching {return ctx.cache}
