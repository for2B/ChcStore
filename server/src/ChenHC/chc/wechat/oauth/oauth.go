package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"ChenHC/chc/wechat/context"
	"ChenHC/chc/wechat/util"
)

const (
	redirectOauthURL      = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&agentid=%s&state=%s#wechat_redirect"
	accessTokenURL        = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	getUserIdURL          = "https://qyapi.weixin.qq.com/cgi-bin/user/getuserinfo?access_token=%s&code=%s"
	refreshAccessTokenURL = "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s"
	userInfoURL           = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
	checkAccessTokenURL   = "https://api.weixin.qq.com/sns/auth?access_token=%s&openid=%s"
)

//Oauth 保存用户授权信息
type Oauth struct {
	*context.Context
}

//NewOauth 实例化授权信息
func NewOauth(context *context.Context, w http.ResponseWriter, r *http.Request) *Oauth {
	auth := new(Oauth)
	auth.Context = context
	auth.Context.Request = r
	auth.Context.Writer = w
	return auth
}

//GetRedirectURL 获取跳转的url地址
func (oauth *Oauth) GetRedirectURL(redirectURI, scope, agentid, state string) (string, error) {
	//url encode
	urlStr := url.QueryEscape(redirectURI)
	return fmt.Sprintf(redirectOauthURL, oauth.AppID, urlStr, scope, agentid, state), nil
}

//Redirect 跳转到网页授权
func (oauth *Oauth) Redirect(writer http.ResponseWriter, req *http.Request, redirectURI, scope, agentid, state string) error {
	location, err := oauth.GetRedirectURL(redirectURI, scope, agentid, state)
	if err != nil {
		return err
	}
	http.Redirect(writer, req, location, 302)
	return nil
}

// ResAccessToken 获取用户授权access_token的返回结果
type ResAccessToken struct {
	util.CommonError

	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// GetUserAccessToken 通过网页授权的code 换取access_token(区别于context中的access_token)
func (oauth *Oauth) GetUserAccessToken() (result ResAccessToken, err error) {
	urlStr := fmt.Sprintf(accessTokenURL, oauth.AppID, oauth.AppSecret)
	var response []byte
	response, err = util.HTTPGet(urlStr)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("GetUserAccessToken error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}

// ResAccessToken 获取用户授权access_token的返回结果
type UserInfo struct {
	util.CommonError
	UserId   string `json:"UserId"`
	OpenId   string `json:"OpenId"`
	DeviceId int64  `json:"DeviceId"`
}

// getUserid
func (oauth *Oauth) GetUserID(acceseeToken, code string) (result UserInfo, err error) {
	urlStr := fmt.Sprintf(getUserIdURL, acceseeToken, code)
	var response []byte
	response, err = util.HTTPGet(urlStr)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("GetUserID error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}

//RefreshAccessToken 刷新access_token
func (oauth *Oauth) RefreshAccessToken(refreshToken string) (result ResAccessToken, err error) {
	urlStr := fmt.Sprintf(refreshAccessTokenURL, oauth.AppID, refreshToken)
	var response []byte
	response, err = util.HTTPGet(urlStr)
	if err != nil {
		return
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("GetUserAccessToken error : errcode=%v , errmsg=%v", result.ErrCode, result.ErrMsg)
		return
	}
	return
}

//CheckAccessToken 检验access_token是否有效
func (oauth *Oauth) CheckAccessToken(accessToken, openID string) (b bool, err error) {
	urlStr := fmt.Sprintf(checkAccessTokenURL, accessToken, openID)
	var response []byte
	response, err = util.HTTPGet(urlStr)
	if err != nil {
		return
	}
	var result util.CommonError
	err = json.Unmarshal(response, &result)
	if err != nil {
		return
	}
	if result.ErrCode != 0 {
		b = false
		return
	}
	b = true
	return
}
