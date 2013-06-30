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
	"errors"
	"fmt"
	"github.com/robfig/revel"
	m "smart-kids/models"
	q "smart-kids/query"
	"time"
)

const (
	simpleQueryTpl = "select %s from %s where %s = ?"
)

var (
	// may move query module to build sql
	userByNameSql       = fmt.Sprintf(simpleQueryTpl, m.UserFields, m.USER_TABLE, m.F_USER_NAME)
	bannedUserByNameSql = fmt.Sprintf(simpleQueryTpl, m.BannedUserFields,
		m.BANNED_USER_TABLE, m.F_USER_NAME)
)

type Application struct {
	*GorpController
}

func (c Application) GetClientInfo() (string, string) {
	appKey := c.Request.Header.Get(m.PARAM_CLIENT_ID)
	appSecret := c.Request.Header.Get(m.PARAM_CLIENT_SECRET)
	return appKey, appSecret
}

// Returns User of the specified userId.
func (c Application) findUser(userId int64) *m.User {
	return m.ToUser(c.Txn.Get(m.User{}, userId))
}

// Returns User of the specified userName, or nil if userName not exists.
func (c Application) findUserByName(userName string) *m.User {
	users := m.ToUsers(c.Txn.Select(m.User{}, userByNameSql, userName))
	if len(users) == 0 {
		return nil
	}
	return users[0]
}

func (c Application) findValidUserByName(userName string) (*m.User, error) {
	bUsers := m.ToBannedUsers(c.Txn.Select(m.BannedUser{}, bannedUserByNameSql, userName))
	if len(bUsers) == 0 {
		return c.findUserByName(userName), nil
	}
	timeNow := time.Now()
	for _, bUser := range bUsers {
		if bUser.IsPermanent {
			return nil, c.Message("users.permanentBannedUser", bUser.UserName, bUser.Cause)
		}
		if bUser.UnbanTime.Valid && bUser.UnbanTime.Time.After(timeNow) {
			return nil, errors.New(c.Message("users.timelinessBannedUser",
				bUser.UserName, bUser.BannedTime.Time, bUser.UnbanTime.Time, bUser.Cause))
		}
	}
	return c.findUserByName(userName), nil
}
