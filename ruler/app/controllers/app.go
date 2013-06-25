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
)

const (
	DEFAULT_PAGE_SIZE = 10
)

var (
	sqlMainMenus = fmt.Sprintf("%s WHERE parent_id = -1 and top_id = -1",
		models.BASE_QUERY_RESOURCE)
	sqlSubMenus = fmt.Sprintf("%s WHERE is_menu = 1 and parent_id > 0 ORDER BY res_url ASC",
		models.BASE_QUERY_RESOURCE)
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
	admins, err := c.Txn.Select(models.Admin{}, models.QUERY_ADMIN_BY_NAME, adminName)
	if err != nil {
		panic(err)
	}
	if len(admins) == 0 {
		return nil
	}
	return admins[0].(*models.Admin)
}

func (c Application) AddMenus() revel.Result {
	mainMenus := loadResources(c.Txn.Select(models.Resource{}, sqlMainMenus))
	subMenus := loadResources(c.Txn.Select(models.Resource{}, sqlSubMenus))

	var idIdxMap = make(map[int]int)
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
