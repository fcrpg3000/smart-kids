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
	_ "encoding/json"
	_ "fmt"
	"github.com/robfig/revel"
	"log"
	m "smart-kids/ruler/app/models"
	_ "smart-kids/ruler/app/routes"
	"smart-kids/util"
)

type Administrators struct {
	Application
}

// Returns admin of the specified id.
func (a Administrators) findAdmin(id int) *m.Admin {
	return m.ToAdmin(a.Txn.Get(m.Admin{}, id))
}

// Returns page admin of the pageable.
func (a Administrators) findAllAdmin(pageable *util.Pageable) *util.Page {

	var (
		total   int64
		content []interface{}
		err     error
	)
	total, err = a.Txn.SelectInt(m.CountSql(m.F_ADMIN_NAME, m.ADMIN_TABLE))
	if total == 0 || err != nil {
		return util.NewPage(nil, pageable, total)
	}

	if pageable == nil {
		content, err = a.Txn.Select(m.Admin{}, m.BASE_QUERY_ADMIN)
	} else {
		sql := m.NewSqlBuilder(m.BASE_QUERY_ADMIN).
			PageOrderBy(pageable, util.AscendingSort([]string{m.F_ADMIN_NAME})).
			ToSqlString()
		content, err = a.Txn.Select(m.Admin{}, sql)
	}
	if err != nil {
		panic(err)
	}
	return util.NewPage(content, pageable, total)
}

// Returns all roles in application
func (a Administrators) findRoles() []*m.Role {
	return m.ToRoles(a.Txn.Select(m.Role{}, m.BASE_QUERY_ROLE))
}

// Pagination admin
func (a Administrators) AdminList(p int) revel.Result {
	pageable, err := util.NewPageable(p, 2, util.ASC, []string{m.F_ADMIN_NAME})
	if err != nil { // never heppen
		log.Fatalf("Error for %s", err.Error())
		panic(err)
	}
	pageAdmin := a.findAllAdmin(pageable)
	return a.Render(pageAdmin)
}

func (a Administrators) changeAdminEnabled(admin *m.Admin, isEnabled bool) (int64, error) {
	admin.IsEnabled = isEnabled
	return a.Txn.Update(admin)
}

// Disable admin of the specified id.
func (a Administrators) DisableAdmin(id int) revel.Result {
	return a.updateAdminEnabled(id, false)
}

// Enable admin of the specified id.
func (a Administrators) EnableAdmin(id int) revel.Result {
	return a.updateAdminEnabled(id, true)
}

// Update admin's enabled property
func (a Administrators) updateAdminEnabled(id int, isEnabled bool) revel.Result {
	targetAdmin := a.findAdmin(id)
	if targetAdmin == nil {
		return a.RenderJson(util.FailureResult(a.Message("admin.notFound")))
	}
	// Cannot ban superadmin of the system built
	if targetAdmin.AdminName == "admin" {
		return a.RenderJson(util.ErrorResult(a.Message("admin.cannotUpdateSuper")))
	}
	_, err := a.changeAdminEnabled(targetAdmin, isEnabled)
	if err != nil {
		panic(err)
	}
	return a.RenderJson(util.SuccessResult(a.Message("operation.successFul")))
}

func (a Administrators) AdminDetail(id int) revel.Result {
	var admin *m.Admin
	var roles = make([]*m.Role, 0)
	allRoles := a.findRoles()
	if id <= 0 {
		roles = allRoles
		title := a.Message("AdminDetail.title.creation")
		return a.Render(title, roles)
	}
	admin = a.findAdmin(id)
	if admin == nil {
		roles = allRoles
		title := a.Message("AdminDetail.title.creation")
		return a.Render(title, roles)
	}
	title := a.Message("AdminDetail.title.edit", admin.AdminName,
		admin.EmpName.String, admin.EmpNo.String)
	adminRoles := admin.Roles
	if len(adminRoles) == 0 {
		roles = allRoles
	} else {
	allRoleLabel:
		for _, role := range allRoles {
			for _, adminRole := range adminRoles {
				if role.Id == adminRole.Id {
					continue allRoleLabel
				}
			}
			roles = append(roles, role)
		}
	}
	return a.Render(title, admin, roles)
}

func (a Administrators) SaveAdmin(admin m.Admin) revel.Result {
	var (
		// row     int64
		err     error
		message string
	)
	if admin.Id <= 0 { // Insert
		err = a.Txn.Insert(&admin)
	} else { // Update
		_, err = a.Txn.Update(&admin)
	}

	if err != nil {
		message = a.Message("admin.errorEdit", err.Error())
	} else {
		message = a.Message("admin.successEdit", admin.AdminName)
	}
	return a.RenderJson(util.SuccessResult(message))
}

// Check if admin name is available.
func (a Administrators) CheckAdminName(adminName string) revel.Result {
	if len(adminName) == 0 {
		return a.RenderJson(util.ErrorResult(a.Message("admin.errorName")))
	}
	admin := a.getAdmin(adminName)
	if admin == nil {
		return a.RenderJson(util.SuccessResult(a.Message("admin.nameNotFound")))
	}
	return a.RenderJson(util.FailureResult(a.Message("admin.errorExistName")))
}
