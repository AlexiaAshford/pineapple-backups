package request

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
)

type Context struct {
	Url      string
	Params   string
	GetQuery []map[string]interface{}
	Structs  interface{}
}

//func (c *Context) Query(key string, value interface{}) {
//	c.GetQuery = append(c.GetQuery, map[string]interface{}{key: value})
//}

func (c *Context) Query(key string, value string) *Context {
	if key != "" && value != "" {
		c.GetQuery = append(c.GetQuery, map[string]interface{}{key: value})
	}
	return c
}
func (c *Context) Init(url string) *Context {
	c.Url = url
	c.Params = ""
	return c
}

func (c *Context) AddCatToken() {
	if config.Vars.AppType == "cat" {
		c.Query("login_token", config.Apps.Cat.Params.LoginToken).
			Query("account", config.Apps.Cat.Params.Account).
			Query("app_version", config.Apps.Cat.Params.AppVersion).
			Query("device_token", config.Apps.Cat.Params.DeviceToken)
	}
}
func (c *Context) QueryToString() string {
	c.AddCatToken()
	for _, queryMap := range c.GetQuery {
		for k, v := range queryMap {
			c.Params += fmt.Sprintf("&%s=%v", k, v)
		}
	}
	c.GetQuery = nil
	return c.Url + "?" + c.Params[1:]
}
