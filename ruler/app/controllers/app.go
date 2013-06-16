package controllers

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/robfig/revel"
	"smart-kids/ruler/app/models"
	"smart-kids/ruler/app/routes"
)

type Application struct {
	GorpController
}

func (c Application) Index() revel.Result {
	if c.connected() == nil {
		c.Flash.Error("请先登录系统再操作！")
		return c.Redirect(routes.Application.Login())
	}
	return c.Render()
}

// Login page
func (c Application) Login() revel.Result {
	if c.connected() != nil {
		return c.Redirect(routes.Application.Index())
	}
	return c.Render()
}

func (c Application) DoLogin(adminName, password string) revel.Result {
	if len(adminName) == 0 || len(password) == 0 {

	}
	admin := c.getAdmin(adminName)
	if admin != nil {
		sha1Hash := sha1.New()
		var srcPassword string
		if len(admin.Salt) == 0 {
			srcPassword = fmt.Sprintf("%s{%s}", password, admin.AdminName)
		} else {
			srcPassword = fmt.Sprintf("%s{%s}", password, admin.Salt)
		}
		sha1Hash.Write([]byte(srcPassword))
		if admin.HashPassword == hex.EncodeToString(sha1Hash.Sum(nil)) {
			c.Session["AdminName"] = adminName
			return c.Redirect(routes.Application.Index())
		}
	}
	c.Flash.Out["adminName"] = adminName
	c.Flash.Error("用户名或密码错误！")
	return c.Redirect(routes.Application.Login())
}

func (c Application) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.Application.Login())
}

func (c Application) AddAdmin() revel.Result {
	if admin := c.connected(); admin != nil {
		c.RenderArgs["admin"] = admin
	}
	return nil
}

func (c Application) connected() *models.Admin {
	if c.RenderArgs["admin"] != nil {
		return c.RenderArgs["admin"].(*models.Admin)
	}
	if adminName, ok := c.Session["AdminName"]; ok {
		return c.getAdmin(adminName)
	}
	return nil
}

func (c Application) getAdmin(adminName string) *models.Admin {
	admins, err := c.Txn.Select(models.Admin{}, models.QUERY_ADMIN_BY_NAME, adminName)
	if err != nil {
		panic(err)
	}
	if len(admins) == 0 {
		return nil
	}
	return admins[0].(*models.Admin)
}
