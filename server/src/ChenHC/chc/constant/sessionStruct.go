package constant

import "time"

type Session struct {
	Openid           string
	Userid           string
	Level            string
	RedirectURL      string
	Name             string
	Identity         string
	ExpireTime       time.Time
	Department_array []string
}
