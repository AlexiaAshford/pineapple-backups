package https

import (
	"fmt"
	"github.com/VeronicaAlexia/pineapple-backups/config"
)

type Context struct {
	url      string
	Params   string
	GetQuery []map[string]interface{}
}

func (c *Context) Query(key string, value interface{}) {
	c.GetQuery = append(c.GetQuery, map[string]interface{}{key: value})
}
func (c *Context) AddCatToken() {
	c.Params = ""
	if config.Vars.AppType == "cat" {
		c.Query("login_token", config.Apps.Cat.Params.LoginToken)
		c.Query("account", config.Apps.Cat.Params.Account)
		c.Query("app_version", config.Apps.Cat.Params.AppVersion)
		c.Query("device_token", config.Apps.Cat.Params.DeviceToken)
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
	return c.url + "?" + c.Params
}
