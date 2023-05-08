package component

import (
	"go-mail/memcache"
	"go-mail/pubsub"
)

type AppContext interface {
	GetPubsub() pubsub.Pubsub
	GetCaching() memcache.Caching
	GetAppToken() *memcache.AppToken
	SetAppToken(a *memcache.AppToken)
}

type appCtx struct {
	pb       pubsub.Pubsub
	cache    memcache.Caching
	appToken *memcache.AppToken
}

func NewAppContext(pb pubsub.Pubsub, cache memcache.Caching, appToken *memcache.AppToken) *appCtx {
	return &appCtx{pb: pb, cache: cache, appToken: appToken}
}

func (ctx *appCtx) GetPubsub() pubsub.Pubsub         { return ctx.pb }
func (ctx *appCtx) GetCaching() memcache.Caching     { return ctx.cache }
func (ctx *appCtx) SetAppToken(a *memcache.AppToken) { ctx.appToken = a }
func (ctx *appCtx) GetAppToken() *memcache.AppToken  { return ctx.appToken }
