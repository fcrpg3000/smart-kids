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
	_ "fmt"
	"github.com/robfig/revel"
	"smart-kids/ruler/app/models"
	_ "smart-kids/ruler/app/routes"
	"smart-kids/util"
)

type Administrators struct {
	Application
}

func (a Administrators) findAdmin(id int) *models.Admin {
	obj, err := a.Txn.Get(models.Admin{}, id)
	if err != nil {
		panic(err)
	}
	if obj == nil {
		return nil
	}
	return obj.(*models.Admin)
}

func (a Administrators) AdminList(p int) revel.Result {

	return a.Render()
}

// Disable admin
func (a Administrators) disableAdminInternal(admin *models.Admin) (int64, error) {
	admin.IsEnabled = false
	return a.Txn.Update(admin)
}

// Disabled admin of the specified id.
func (a Administrators) DisableAdmin(id int) revel.Result {
	targetAdmin := a.findAdmin(id)
	if targetAdmin == nil {
		return a.RenderJson(util.FailureResult(a.Message("admin.notFound")))
	}
	// Cannot ban superadmin of the system built
	if targetAdmin.AdminName == "admin" {
		return a.RenderJson(util.ErrorResult(a.Message("admin.cannotBanSuper")))
	}
	_, err := a.disableAdminInternal(targetAdmin)
	if err != nil {
		panic(err)
	}
	return a.RenderJson(util.SuccessResult(a.Message("admin.successBanned")))
}

func (a Administrators) AdminDetail(id int) revel.Result {
	admin := a.findAdmin(id)
	if admin == nil {
		//return a.NotFound(a.Message("admin.notFound"))
		return a.Render()
	} else {
		title := a.Message("AdminDetail.title", admin.AdminName,
			admin.EmpName, admin.EmpNo)
		return a.Render(title, admin)
	}
}
