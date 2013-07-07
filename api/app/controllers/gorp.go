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
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/revel"
	"github.com/robfig/revel/modules/db/app"
	"smart-kids/models"
)

var (
	Dbm *gorp.DbMap
)

// Application Initialize
func Init() {
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	initUsers()
	initApp()
	Dbm.TraceOn("[gorp]", revel.INFO)
}

func setColumnSizes(t *gorp.TableMap, colSizes map[string]int) {
	for col, size := range colSizes {
		t.ColMap(col).MaxSize = size
	}
}

func initUsers() {
	// Register User model
	t := Dbm.AddTableWithName(models.User{}, models.USER_TABLE).SetKeys(true, "UserId")
	setColumnSizes(t, map[string]int{
		"UserName":      50,
		"HashPassword":  100,
		"PasswordSalt":  100,
		"Email":         50,
		"CommonlyEmail": 50,
	})
	t.ColMap("UserName").SetUnique(true)

	// Register UserDigital model
	t = Dbm.AddTableWithName(models.UserDigital{}, models.USER_DIGITAL_TABLE).SetKeys(false, "UserId")
	setColumnSizes(t, map[string]int{"UserName": 50})

	// Register UserInfo model
	t = Dbm.AddTableWithName(models.UserInfo{}, models.USER_INFO_TABLE).SetKeys(false, "UserId")
	setColumnSizes(t, map[string]int{
		"UserName":       50,
		"Nickname":       50,
		"DateOfBirthStr": 10,
		"OtherState":     100,
	})

	// Register UserAvatar model
	t = Dbm.AddTableWithName(models.UserAvatar{}, models.USER_AVATAR_TABLE).SetKeys(false, "Id")
	setColumnSizes(t, map[string]int{
		"UserName":        50,
		"ImageDomain":     100,
		"AvatarPath":      255,
		"SrcAvatarPath":   255,
		"SmallAvatarPath": 255,
		"ThumbAvatarPath": 255,
		"AvatarName":      250,
	})

	// Register BannedUser model
	t = Dbm.AddTableWithName(models.BannedUser, models.BANNED_USER_TABLE).SetKeys(true, "Id")
	setColumnSizes(t, map[string]int{
		"UserName":           50,
		"OperatorName":       50,
		"Cause":              2000,
		"LastModifiedByName": 50,
	})
}

func initApp() {
	// Register Developer model
	t := Dbm.AddTableWithName(models.Developer{}, models.DEVELOPER_TABLE).SetKeys(false, "UserId")
	setColumnSizes(t, map[string]int{
		"UserName": 50,
		"DevName":  200,
		"Email":    50,
		"Phone":    20,
		"DevIm":    50,
		"DevSite":  255,
	})

	t = Dbm.AddTableWithName(models.App{}, models.APP_TABLE).SetKeys(true, "Id")
	setColumnSizes(t, map[string]int{
		"Name":        100,
		"Url":         255,
		"Summary":     100,
		"Description": 3000,
		"UserName":    50,
		"AppKey":      100,
		"AppSecret":   100,
	})
	t.ColMap("Name").SetUnique(true)
	t.ColMap("Url").SetUnique(true)
	t.ColMap("AppKey").SetUnique(true)
	t.ColMap("AppSecret").SetUnique(true)

	t = Dbm.AddTableWithName(models.AppSession{}, models.APP_SESSION_TABLE).SetKeys(false, "AppId")
	setColumnSizes(t, map[string]int{
		"AppName":     100,
		"AppAuthCode": 50,
		"AccessToken": 50,
		"AppKey":      100,
		"AppSecret":   100,
	})
	t.ColMap("AccessToken").SetUnique(true)
}

type GorpController struct {
	*revel.Controller
	Txn *gorp.Transaction
}

func (c *GorpController) Begin() revel.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() revel.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() revel.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}
