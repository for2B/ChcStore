package context

import (
	"ChenHC/chc/cache"
	"net/http"
	"sync"
)

// Context struct
type Context struct {
	AppID        string
	AppSecret    string
	VIPAppID     string
	VIPAppSecret string

	Token           string
	EncodingAESKey  string
	PayMchID        string
	PayNotifyURL    string
	PayKey          string
	accessTokenLock *sync.RWMutex
	Cache           *cache.CacheManager

	Writer  http.ResponseWriter
	Request *http.Request
}

// Query returns the keyed url query value if it exists
func (ctx *Context) Query(key string) string {
	value, _ := ctx.GetQuery(key)
	return value

}

// GetQuery is like Query(), it returns the keyed url query value
func (ctx *Context) GetQuery(key string) (string, bool) {
	req := ctx.Request //ctx的request为空指针所以panic
	if values, ok := req.URL.Query()[key]; ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}
