package context

import (
	"ChenHC/internal/httpapi"
	"ChenHC/chc/constant"
	"encoding/json"
	"fmt"

	"ChenHC/chc/wechat/util"
	"sync"
	"time"
)

const (
	//AccessTokenURL 获取access_token的接口
	AccessTokenURL = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
)

//ResAccessToken struct
type ResAccessToken struct {
	util.CommonError

	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	ExpireTime  time.Time
}

//SetAccessTokenLock 设置读写锁（一个appID一个读写锁）
func (ctx *Context) SetAccessTokenLock(l *sync.RWMutex) {
	ctx.accessTokenLock = l
}

//GetAccessToken 获取access_token(先cache-->没有就去wechat server获取)
func (ctx *Context) GetAccessToken(isvip int) (accessToken string, err error) {
	var resAccessToken ResAccessToken

	ctx.accessTokenLock.Lock()
	defer ctx.accessTokenLock.Unlock()
	appID := ""
	if isvip == 1 {
		appID = ctx.VIPAppID + "VIP"
	}
	appID = ctx.AppID
	accessTokenCacheKey := fmt.Sprintf("access_token_%s", appID)

	err = ctx.Cache.Get("accsstoken", accessTokenCacheKey, &resAccessToken)
	if err != nil {
		if err.Error() == "not found" {
			//从微信服务器获取
			resAccessToken, err = ctx.GetAccessTokenFromServer(isvip)
			if err != nil {
				return
			}
			accessToken = resAccessToken.AccessToken
			return
		} else {
			return "", httpapi.NewErr(constant.GLOBAL_SYS_ERR, "get the key accesstoken fail", err)

		}

	}

	//过期删除
	if resAccessToken.ExpireTime.After(time.Now()) {
		ctx.Cache.Del("accsstoken", accessTokenCacheKey)
		return "", httpapi.Err{Code: constant.GLOBAL_SYS_ERR, Msg: "session is expired,please auth again"}
	}
	// 更新ACESSTOKEN的过期时间 ，定位7000ms
	resAccessToken.ExpireTime = time.Now().Add(time.Duration(7000) * time.Second)
	// fmt.Println("resAccessToken.ExpireTime=", resAccessToken.ExpireTime)

	return resAccessToken.AccessToken, nil
}

//GetAccessTokenFromServer 从微信服务器获取token
func (ctx *Context) GetAccessTokenFromServer(isvip int) (resAccessToken ResAccessToken, err error) {

	appID := ""
	secret := ""

	// 根据是否vip获取不同的at
	if isvip == 1 {
		appID = ctx.VIPAppID
		secret = ctx.VIPAppSecret
	}
	appID = ctx.AppID
	secret = ctx.AppSecret

	url := fmt.Sprintf(AccessTokenURL, appID, secret)

	var body []byte
	body, err = util.HTTPGet(url)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resAccessToken)
	if err != nil {
		return
	}

	if 0 != resAccessToken.CommonError.ErrCode {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v", resAccessToken.ErrCode, resAccessToken.ErrMsg)
		return
	}

	// 根据是否vip获取不同的at
	cache_key := appID
	if isvip == 1 {
		cache_key = appID + "VIP"
	}
	accessTokenCacheKey := fmt.Sprintf("access_token_%s", cache_key)

	// 设置过期时间
	resAccessToken.ExpireTime = time.Now().Add(time.Duration(7000) * time.Second)
	err = ctx.Cache.Set("accesstoken", accessTokenCacheKey, resAccessToken)
	if err != nil {
		return resAccessToken, httpapi.NewErr(constant.GLOBAL_SYS_ERR, "set accesstoken key in cache failed", err)
	}

	return
}
