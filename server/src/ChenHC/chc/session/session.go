package session

import (
	"ChenHC/internal/httpapi"
	"ChenHC/chc/cache"
	"ChenHC/chc/constant"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"
	"time"
	"fmt"
)

var (
	cookieName  string
	maxLifeTime int
	path        string
	httpOnly    bool
	maxAge      int
)

type SessionManager struct {
	SessionID  string
	data       map[string]*constant.Session  //session的数据 [随机十六位] 对应的session
	cache      *cache.CacheManager //使用cache来存储session
	Build_type string
}

func (s *SessionManager) Init(SessionKey string, MaxLifeTime int, Path string, HTTPOnly bool, MaxAge int, c *cache.CacheManager, BUILD_TYPE string) {
	cookieName = SessionKey
	maxLifeTime = MaxLifeTime
	path = Path
	httpOnly = HTTPOnly
	maxAge = MaxAge
	s.cache = c
	s.Build_type = BUILD_TYPE
	fmt.Println("init SessionManager success!")
}

/*
*将session存入cache，w,r,已经设置好值得session
 */
func (s *SessionManager) SetSession(w http.ResponseWriter, r *http.Request, session *constant.Session) error {
	var sid string
	cookie, err := r.Cookie(cookieName)  //判断r中是否带有cookie
	if err != nil && err != http.ErrNoCookie {
		return httpapi.NewErr(constant.GLOBAL_SYS_ERR, "get cookie failed", err)
	}
	if err == http.ErrNoCookie || cookie.Value == "" {
		sid = s.randomSID()
		// fmt.Println("sid=", sid)
		newCookie := http.Cookie{
			Name:     cookieName,
			Value:    url.QueryEscape(sid),
			Path:     path,
			HttpOnly: httpOnly,
			MaxAge:   maxAge,
		}
		http.SetCookie(w, &newCookie)  //设置到http随响应返回前端
	} else {
		// cookie不为空
		sid, err = url.QueryUnescape(cookie.Value)
		// fmt.Println("sid=", sid)
		if err != nil {
			return httpapi.NewErr(constant.GLOBAL_SYS_ERR, "Unescape the client's cookie failed", err)
		}
	}
	// 设置过期时间
	// fmt.Println("time.Duration(maxLifeTime) * time.Second", time.Duration(maxLifeTime)*time.Second)
	session.ExpireTime = time.Now().Add(time.Duration(maxLifeTime) * time.Second)
	// fmt.Println("session.ExpireTime=", session.ExpireTime)
	err = s.cache.Set("session", constant.CACHE_SESSION+sid, session)
	if err != nil {
		return httpapi.NewErr(constant.GLOBAL_SYS_ERR, "set session in cache failed", err)
	}
	return nil
}

/*
*从cache读取session
 */
func (s *SessionManager) GetSession(w http.ResponseWriter, r *http.Request) (sess *constant.Session, err error) {
	// 取出请求带来的cookie
	// level_cookie, err := r.Cookie("level")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("level_cookie", level_cookie)

	cookie, err := r.Cookie(cookieName)
	// fmt.Println("cookie=", cookie)

	if err != nil {
		if err == http.ErrNoCookie {
			return nil, httpapi.NewErr(constant.COOKIE_NULL, "client request without cookie", err)
		}
		return nil, httpapi.NewErr(constant.GLOBAL_SYS_ERR, "r.cookie(key) func failed", err)
	}
	if cookie.Value == "" {
		return nil, httpapi.NewErr(constant.COOKIE_EMPTY, "the content of cookie is empty", nil)
	}

	//获得sid
	sid, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return nil, httpapi.NewErr(constant.GLOBAL_SYS_ERR, "url.QueryUnescape(cookie.Value) func fail", err)
	}
	sess = &constant.Session{}
	err = s.cache.Get("session", constant.CACHE_SESSION+sid, sess)
	if err != nil {
		if err.Error() == "not found" {
			Cookie := http.Cookie{
				Name:     cookieName,
				Value:    url.QueryEscape(""),
				Path:     "/",
				HttpOnly: false,
				MaxAge:   -1,
			}
			http.SetCookie(w, &Cookie)

			return nil, httpapi.NewErr(constant.COOKIE_CACHE_NOFOUND, "have no this cookie's key in cache", err)
		}
		return nil, httpapi.NewErr(constant.GLOBAL_SYS_ERR, "get the key sid in cache failed", err)
	}

	//过期删除
	// fmt.Println("time.Now()", time.Now())
	if sess.ExpireTime.Before(time.Now()) {
		s.DestroySession(r)
		return nil, httpapi.NewErr(constant.SESSION_EXPIRED, "session is expired,please auth again", nil)
	}
	// 更新过期时间
	sess.ExpireTime = time.Now().Add(time.Duration(maxLifeTime) * time.Second)
	// fmt.Println("sess.expiretime=", sess.ExpireTime)

	return
}
//从cache中删除session
func (s *SessionManager) DestroySession(r *http.Request) error {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return nil
		}
		return err
	}
	sid, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return err
	}
	return s.cache.Del("session", constant.CACHE_SESSION+sid)
}
//获取16位进制随机编号字符
func (s *SessionManager) randomSID() string {
	sid := make([]byte, 16)
	rand.Read(sid)
	return hex.EncodeToString(sid)
}
