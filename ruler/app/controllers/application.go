// Copyright (C) 2012-2013 king4go authors All rights reserved.
//
// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License. You may obtain a copy of the License at
//
//           http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package controllers

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/robfig/revel"
	_ "log"
	"smart-kids/ruler/app/models"
	"smart-kids/ruler/app/routes"
	"smart-kids/util"
)

const (
	DEFAULT_PAGE_SIZE = 10
)

var (
	sqlMainMenus = fmt.Sprintf("%s WHERE parent_id = 0 and top_id = 0",
		models.BASE_QUERY_RESOURCE)
	sqlSubMenus = fmt.Sprintf("%s WHERE is_menu = 1 and parent_id > 0 ORDER BY res_url ASC",
		models.BASE_QUERY_RESOURCE)
)

type Application struct {
	*GorpController
}

func (c Application) NotFoundMessage(entity string) string {
	return c.Message("resource.notFound", entity)
}

func (c Application) Index() revel.Result {
	if admin := c.connected(); admin == nil {
		return c.Redirect(routes.Application.Login())
	}
	return c.Render()
}

// Login page
func (c Application) Login() revel.Result {
	if c.connected() != nil {
		return c.Redirect(routes.Application.Index())
	}
	redirectUrl := c.Params.Get("redirectUrl")
	return c.Render(redirectUrl)
}

func (c Application) DoLogin(adminName, password string) revel.Result {
	if len(adminName) == 0 || len(password) == 0 {

	}
	redirectUrl := c.Params.Get("redirectUrl")
	if len(redirectUrl) == 0 {
		redirectUrl = routes.Application.Index()
	}
	admin := c.getAdmin(adminName)
	if admin != nil {
		if !admin.IsEnabled {
			c.Flash.Error("用户(%s)已被禁用", adminName)
		} else {
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
				return c.Redirect(redirectUrl)
			}
		}
	} else {
		c.Flash.Error("用户名或密码错误！")
	}
	c.Flash.Out["adminName"] = adminName

	if len(redirectUrl) > 0 && redirectUrl != routes.Application.Index() {
		redirectUrl = util.AddParamsToUrl(routes.Application.Login(), map[string]string{
			"redirectUrl": redirectUrl,
		})
	} else {
		redirectUrl = routes.Application.Login()
	}
	return c.Redirect(redirectUrl)
}

func (c Application) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.Application.Login())
}

func (c Application) AddAdmin() revel.Result {
	if admin := c.connected(); admin != nil {
		c.RenderArgs["adminUser"] = admin
	}
	return nil
}

// Intercepter method for every request
func (c Application) checkLogin() revel.Result {
	requestURI := c.Request.URL.Path
	if requestURI == "/login" || requestURI == "/do_login" || requestURI == "/logout" {
		return nil
	}
	if admin := c.connected(); admin == nil {
		redirectUrl := routes.Application.Login()
		if len(requestURI) > 0 && requestURI != "/" {
			redirectUrl = util.AddParamsToUrl(routes.Application.Login(), map[string]string{
				"redirectUrl": c.Request.URL.String(),
			})
		}
		return c.Redirect(redirectUrl)
	} else {
		c.RenderArgs["adminUser"] = admin
	}
	return nil
}

func (c Application) connected() *models.Admin {
	if c.RenderArgs["adminUser"] != nil {
		return c.RenderArgs["adminUser"].(*models.Admin)
	}
	if adminName, ok := c.Session["AdminName"]; ok {
		return c.getAdmin(adminName)
	}
	return nil
}

func (c Application) getAdmin(adminName string) *models.Admin {
	admins := models.ToAdmins(c.Txn.Select(models.Admin{},
		models.QUERY_ADMIN_BY_NAME, adminName))
	if len(admins) == 0 {
		return nil
	}
	return admins[0]
}

func (c Application) AddMenus() revel.Result {
	if c.RenderArgs["adminUser"] == nil {
		return nil
	}
	mainMenus := models.ToResources(c.Txn.Select(models.Resource{}, sqlMainMenus))
	subMenus := models.ToResources(c.Txn.Select(models.Resource{}, sqlSubMenus))

	var idIdxMap = make(map[uint]int)
	for i, mainMenu := range mainMenus {
		idIdxMap[mainMenu.Id] = i
	}

	for _, subMenu := range subMenus {
		if idx, ok := idIdxMap[subMenu.ParentId]; ok {
			children := mainMenus[idx].Children
			children = append(children, subMenu)
			mainMenus[idx].Children = children
		}
	}
	c.RenderArgs["mainMenus"] = mainMenus
	return nil
}
