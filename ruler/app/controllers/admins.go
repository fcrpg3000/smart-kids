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
	_ "encoding/json"
	"fmt"
	"github.com/robfig/revel"
	"log"
	m "smart-kids/ruler/app/models"
	_ "smart-kids/ruler/app/routes"
	"smart-kids/util"
	"strings"
	_ "time"
)

type Administrators struct {
	Application
}

func hashPassword(srcPwd, salt string) (string, string) {
	if len(salt) == 0 {
		salt = util.RandomAlphanumeric(8)
	}
	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(fmt.Sprintf("%s{%s}", srcPwd, salt)))
	hashPwd := hex.EncodeToString(sha1Hash.Sum(nil))
	return hashPwd, salt
}

// Returns admin of the specified id.
func (a Administrators) findAdmin(id uint32) *m.Admin {
	return m.ToAdmin(a.Txn.Get(m.Admin{}, id))
}

func (a Administrators) findAdminByName(name string) (*m.Admin, bool) {
	if len(name) == 0 {
		return nil, false
	}
	admins := m.ToAdmins(a.Txn.Select(m.Admin{}, m.QUERY_ADMIN_BY_NAME, name))
	if len(admins) == 0 {
		return nil, false
	}
	return admins[0], true
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

func (a Administrators) validationAdmin(admin *m.Admin) {
	nameLength := a.Message("admin.v.nameLength", 4, 16)
	a.Validation.Range(len(admin.AdminName), 4, 16).
		Key("admin.AdminName").Message(nameLength)
	if len(admin.EmpNameValue) > 0 {
		empNameLength := a.Message("admin.v.empNameLength", 2, 10)
		a.Validation.Range(len(admin.EmpNameValue), 2, 10).
			Key("admin.EmpNameValue").Message(empNameLength)
		admin.EmpName.Scan(admin.EmpNameValue)
	}
	if len(admin.EmpNoValue) > 0 {
		empNoLength := a.Message("admin.v.empNoLength", 2, 16)
		a.Validation.Range(len(admin.EmpNoValue), 2, 16).
			Key("admin.EmpNoValue").Message(empNoLength)
		admin.EmpNo.Scan(admin.EmpNoValue)
	}
}

func (a Administrators) SaveAdmin(admin m.Admin) revel.Result {
	var (
		updated, exists *m.Admin
		found           bool
		err             error
	)
	if exists, found = a.findAdminByName(admin.AdminName); found {
		if admin.Id <= 0 || (admin.Id > 0 && exists.Id != admin.Id) {
			return a.RenderJson(util.FailureResult(a.Message("admin.errorExistName")))
		}
	}
	if admin.Id > 0 { // update admin
		updated = a.findAdmin(admin.Id)
		if updated == nil {
			return a.RenderJson(util.ErrorResult(a.Message("admin.notFound")))
		}
	}

	a.validationAdmin(&admin)
	if a.Validation.HasErrors() {
		result := util.FailureResult("保存失败")
		for k, v := range a.Validation.ErrorMap() {
			if v != nil {
				result.AddValue(k, v.Message)
			}
		}
		return a.RenderJson(result)
	}
	admin.LastIp.Scan(a.clientIp())
	if admin.Id <= 0 { // create new admin
		admin.HashPassword, admin.Salt = hashPassword("123456", "")
		if currentAdmin := a.connected(); currentAdmin != nil {
			admin.CreatedBy(currentAdmin.Id, currentAdmin.AdminName)
		}
		err = a.Txn.Insert(&admin)
	} else {
		_, err = a.Txn.Update(updated.UpdateBy(&admin))
	}
	if err != nil {
		return a.RenderJson(util.ErrorResult(a.Message("admin.errorEdit", err.Error())))
	}
	return a.RenderJson(util.SuccessResult(a.Message("admin.s.saved", admin.AdminName)))
}

// Pagination admin
func (a Administrators) AdminList(p int) revel.Result {
	pageable, err := util.NewPageable(p, DEFAULT_PAGE_SIZE, util.ASC, []string{m.F_ADMIN_NAME})
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
func (a Administrators) DisableAdmin(id uint32) revel.Result {
	return a.updateAdminEnabled(id, false)
}

// Enable admin of the specified id.
func (a Administrators) EnableAdmin(id uint32) revel.Result {
	return a.updateAdminEnabled(id, true)
}

// Update admin's enabled property
func (a Administrators) updateAdminEnabled(id uint32, isEnabled bool) revel.Result {
	targetAdmin := a.findAdmin(id)
	if targetAdmin == nil {
		return a.RenderJson(util.FailureResult(a.Message("admin.notFound")))
	}
	// Cannot ban superadmin of the system built
	if targetAdmin.AdminName == "admin" {
		return a.RenderJson(util.ErrorResult(a.Message("admin.cannotUpdateSuper")))
	}
	targetAdmin.LastIp.Scan(a.clientIp())
	_, err := a.changeAdminEnabled(targetAdmin, isEnabled)
	if err != nil {
		panic(err)
	}
	return a.RenderJson(util.SuccessResult(a.OperOkMessage()))
}

func (a Administrators) AdminDetail(id uint32) revel.Result {
	var admin *m.Admin
	var roles []*m.Role
	allRoles := a.findRoles()
	if id <= 0 {
		roles = allRoles
		title := a.Message("admin.title.creation")
		return a.Render(title, roles)
	}
	admin = a.findAdmin(id)
	if admin == nil {
		roles = allRoles
		title := a.Message("admin.title.creation")
		return a.Render(title, roles)
	}
	title := a.Message("admin.title.detail", admin.AdminName,
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

// To change hash password page.
func (a Administrators) Chpwd() revel.Result {
	currentAdmin := a.connected()
	title := a.Message("admin.title.chpwd")
	return a.Render(title, currentAdmin)
}

// Returns (nil, true) if validate passed, otherwise (result, false).
func (a Administrators) validatePasswd(oldpwd, newpwd1, newpwd2 string) (revel.Result, bool) {
	oldpwdLen, newpwd1Len, newpwd2Len := len(oldpwd), len(newpwd1), len(newpwd2)

	if oldpwdLen == 0 {
		return a.RenderJson(util.FailureResult(a.Message("admin.v.oldpwdRequired"))), false
	}
	if newpwd1Len == 0 && newpwd2Len == 0 {
		return a.RenderJson(util.FailureResult(a.Message("admin.v.newpwdRequired"))), false
	} else if newpwd1Len != newpwd2Len {
		return a.RenderJson(util.FailureResult(a.Message("admin.v.newpwdNotEquals"))), false
	} else if newpwd1Len < 6 && newpwd1Len > 16 {
		return a.RenderJson(util.FailureResult(a.Message("admin.v.passwordLength", 6, 16))), false
	}
	return nil, true
}

// Post request to modify current login admin's password.
func (a Administrators) Passwd(oldpwd, newpwd1, newpwd2 string) revel.Result {
	oldpwd = strings.TrimSpace(oldpwd)
	newpwd1 = strings.TrimSpace(newpwd1)
	newpwd2 = strings.TrimSpace(newpwd2)
	if result, passed := a.validatePasswd(oldpwd, newpwd1, newpwd2); !passed {
		return result
	}
	currentAdmin := a.connected()
	hashOldpwd, _ := hashPassword(oldpwd, currentAdmin.Salt)
	if currentAdmin.HashPassword != hashOldpwd {
		return a.RenderJson(util.ErrorResult(a.Message("admin.v.oldpwdError")))
	}
	newHashPassword, newSalt := hashPassword(newpwd1, "")
	currentAdmin.HashPassword = newHashPassword
	currentAdmin.Salt = newSalt
	_, err := a.Txn.Update(currentAdmin)
	if err != nil {
		return a.RenderJson(util.ErrorResult(err.Error()))
	}
	return a.RenderJson(util.SuccessResult(a.Message("admin.s.modifyPassword")))
}

// Check if admin name is available.
func (a Administrators) CheckAdminName(adminName string) revel.Result {
	if len(adminName) == 0 {
		return a.RenderJson(util.ErrorResult(a.Message("admin.v.nameFormat")))
	}
	_, found := a.findAdminByName(adminName)
	if !found {
		return a.RenderJson(util.SuccessResult(a.Message("admin.nameNotFound")))
	}
	return a.RenderJson(util.FailureResult(a.Message("admin.errorExistName")))
}
