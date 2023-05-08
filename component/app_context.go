package component

import (
	"go-mail/memcache"
	"go-mail/pubsub"
)

type AppContext interface {
	GetPubsub() pubsub.Pubsub
	GetCaching() memcache.Caching
	GetAppToken() memcache.AppAccess
}

type appCtx struct {
	pb        pubsub.Pubsub
	cache	  memcache.Caching
	appToken memcache.AppAccess
}

func NewAppContext(pb pubsub.Pubsub, cache memcache.Caching, appToken memcache.AppAccess) *appCtx {
	return &appCtx{pb: pb, cache: cache, appToken: appToken}
}

func (ctx *appCtx) GetPubsub() pubsub.Pubsub { return ctx.pb }
func(ctx *appCtx) GetCaching() memcache.Caching {return ctx.cache}
func(ctx *appCtx) GetAppToken() memcache.AppAccess {return ctx.appToken}
