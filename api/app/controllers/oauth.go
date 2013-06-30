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
	"smart-kids/util"
	"strings"
)

var (
	appByKeySql = fmt.Sprintf(simpleQueryTpl, m.AppFields, m.APP_TABLE, m.F_APP_KEY)

	loaderFuncErr = errors.New("Loader function return is nil.")
)

type OAuth struct {
	*Application
}

func (o OAuth) getAppSession(appId int64, loader func(int64) *m.AppSession) *m.AppSession {
	appSession := m.ToAppSession(o.Txn.Get(m.AppSession{}, appId))
	if appSession == nil {
		appSession = loader(appId)
		if appSession == nil {
			panic(loaderFuncErr)
		}
	}
	return appSession
}

// API authoirze
func (o OAuth) Authorize() revel.Result {
	clientId, _ := o.GetClientInfo()
	if len(clientId) == 0 {
		clientId = o.Params.Get(m.PARAM_CLIENT_ID)
	}
	if len(clientId) == 0 {
		return o.RenderJson(m.Err_Invalid_Client)
	}

	app := m.ToApp(o.Txn.Select(m.App{}, appByKeySql, clientId))
	if app == nil {
		return o.RenderJson(m.Err_Invalid_Client)
	}

	redirectUri := o.Params.Get("redirect_uri")
	if len(redirectUri) == 0 {
		return o.RenderJson(m.Err_Redirect_URI_Mismatch)
	}

	display := o.Params.Get("display")
	if len(display) == 0 {
		display = "default"
	}
	state := o.Params.Get("state")
	forceLogin := false
	forceLoginStr := o.Params.Get("forcelogin")
	if len(forceLoginStr) > 0 && "true" == forceLoginStr {
		forceLogin = true
	}

	appSession := o.getAppSession(app.Id, func(appId int64) *m.AppSession {
		target := m.NewAppSession(app)
		row, err := o.Txn.Insert(target)
		if err != nil {
			panic(err)
		}
		return target
	})

	code := appSession.FlushAuthCode()
	redirectUrl := util.AddParamsToUrl(redirectUri, map[string]string{
		"state": state,
		"code":  code,
	})
	return o.Redirect(redirectUrl)
}
