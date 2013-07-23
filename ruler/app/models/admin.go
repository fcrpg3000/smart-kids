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
package models

import (
	"database/sql"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/go-sql-driver/mysql"
	"reflect"
	"time"
)

// Admin table's field name constants
const (
	ADMIN_TABLE     = "m_admin"
	F_ADMIN_ID      = "admin_id"
	F_ADMIN_NAME    = "admin_name"
	F_HASH_PASSWORD = "hash_password"
	F_SALT          = "pwd_salt"
	F_USER_ID       = "user_id"
	F_USER_NAME     = "user_name"
	F_EMP_NAME      = "emp_name"
	F_EMP_NO        = "emp_no"
	F_IS_ENABLED    = "is_enabled"
	F_LAST_IP       = "last_ip"
)

type Admin struct {
	Id               uint32         `db:"admin_id"`
	AdminName        string         `db:"admin_name"`
	HashPassword     string         `db:"hash_password"`
	Salt             string         `db:"pwd_salt"`
	UserId           uint64         `db:"user_id"`   // relate to User#Id if necessary
	UserName         sql.NullString `db:"user_name"` // relate to User#UserName if necessary
	EmpName          sql.NullString `db:"emp_name"`
	EmpNo            sql.NullString `db:"emp_no"`
	CreatedById      uint32         `db:"created_by_id"`
	CreatedByName    sql.NullString `db:"created_by_name"`
	CreatedTime      mysql.NullTime `db:"created_time"`
	LastModifiedTime mysql.NullTime `db:"last_modified_time"`
	IsEnabled        bool           `db:"is_enabled"`
	LastIp           sql.NullString `db:"last_ip"`

	// Transient
	EmpNameValue string  `db:"-"`
	EmpNoValue   string  `db:"-"`
	Password     string  `db:"-" json:"-"` // used in form
	Roles        []*Role `db:"-" json:",omitempty"`
}

func (a *Admin) UpdateBy(admin *Admin) *Admin {
	if len(admin.AdminName) > 0 {
		a.AdminName = admin.AdminName
	}
	if len(admin.HashPassword) > 0 && len(admin.Salt) > 0 {
		a.HashPassword = admin.HashPassword
		a.Salt = admin.Salt
	}
	if admin.EmpName.Valid {
		a.EmpName = admin.EmpName
	}
	if admin.EmpNo.Valid {
		a.EmpNo = admin.EmpNo
	}
	a.IsEnabled = admin.IsEnabled
	if admin.LastIp.Valid {
		a.LastIp = admin.LastIp
	}
	return a
}

func (a *Admin) CreatedBy(id uint32, name string) *Admin {
	if id > 0 && len(name) > 0 {
		a.CreatedById = id
		a.CreatedByName.Scan(name)
	}
	return a
}

func (a *Admin) PreInsert(_ gorp.SqlExecutor) error {
	timeNow := time.Now()
	a.CreatedTime.Scan(timeNow)
	a.LastModifiedTime.Scan(timeNow)
	return nil
}

func (a *Admin) PreUpdate(_ gorp.SqlExecutor) error {
	a.LastModifiedTime.Scan(time.Now())
	return nil
}

func (a *Admin) String() string {
	return fmt.Sprintf("Admin{Id=%d, AdminName=%s, HashPassword=%s, Salt=%s, "+
		"UserId=%d, UserName=%v, EmpName=%v, EmpNo=%v, CreatedById=%d, "+
		"CreatedByName=%v, CreatedTime=%v, LastModifiedTime=%v, IsEnabled=%v, "+
		"LastIp=%v}", a.Id, a.AdminName, a.HashPassword, a.Salt, a.UserId,
		a.UserName, a.EmpName, a.EmpNo, a.CreatedById, a.CreatedByName,
		a.CreatedTime, a.LastModifiedTime, a.IsEnabled, a.LastIp)
}

// gorp 不能自动处理关联关系，实现这个接口方法，当调用根据Id获取 Admin 对象时，
// 自动获取关联的角色信息列表
func (a *Admin) PostGet(exe gorp.SqlExecutor) error {
	query := BASE_QUERY_ROLE +
		"WHERE role_id IN (SELECT role_id FROM m_admin_role ar JOIN m_admin a " +
		"ON a.admin_id = ar.admin_id WHERE a.admin_id = ?)"
	a.Roles = ToRoles(exe.Select(Role{}, query, a.Id))
	return nil
}

func ToAdmin(i interface{}, err error) *Admin {
	if err != nil {
		panic(err)
	}
	if i == nil || reflect.ValueOf(i).IsNil() {
		return nil
	}
	return i.(*Admin)
}

func ToAdmins(results []interface{}, err error) []*Admin {
	if err != nil {
		panic(err)
	}
	size := len(results)
	admins := make([]*Admin, size)
	if size == 0 {
		return admins
	}
	for i, r := range results {
		admins[i] = r.(*Admin)
	}
	return admins
}

const (
	ROLE_TABLE  = "m_role"
	F_ROLE_ID   = "role_id"
	F_ROLE_CODE = "role_code"
	F_ROLE_NAME = "role_name"
	F_ROLE_DESC = "role_desc"
)

// mapped table
type Role struct {
	Id               uint32         `db:"role_id"`
	Code             string         `db:"role_code"`
	Name             string         `db:"role_name"`
	Desc             sql.NullString `db:"role_desc"`
	CreatedById      uint32         `db:"created_by_id"`
	CreatedByName    sql.NullString `db:"created_by_name"`
	CreatedTime      mysql.NullTime `db:"created_time"`
	LastModifiedTime mysql.NullTime `db:"last_modified_time"`

	// Transient
	// Relation manage by Role
	Resources []*Resource `db:"-" json:",omitempty"`
}

// gorp Hook pre-Insert
func (r *Role) PreInsert(_ gorp.SqlExecutor) error {
	timeNow := time.Now()
	if !r.Desc.Valid {
		r.Desc = sql.NullString{"<暂无描述>", true}
	}
	r.CreatedTime = mysql.NullTime{timeNow, true}
	r.LastModifiedTime = mysql.NullTime{timeNow, true}
	return nil
}

// gorp Hook pre-Update
func (r *Role) PreUpdate(_ gorp.SqlExecutor) error {
	r.LastModifiedTime = mysql.NullTime{time.Now(), true}
	return nil
}

func ToRole(i interface{}, err error) *Role {
	if err != nil {
		panic(err)
	}
	if i == nil || reflect.ValueOf(i).IsNil() {
		return nil
	}
	return i.(*Role)
}

func ToRoles(results []interface{}, err error) []*Role {
	if err != nil {
		panic(err)
	}
	size := len(results)
	roles := make([]*Role, size)
	if size == 0 {
		return roles
	}
	for i, r := range results {
		roles[i] = r.(*Role)
	}
	return roles
}

const (
	RESOURCE_TABLE = "m_resource"
	F_RES_ID       = "res_id"
	F_RES_NAME     = "res_name"
	F_RES_CODE     = "res_code"
	F_RES_DESC     = "res_desc"
	F_RES_URL      = "res_url"
	F_TOP_ID       = "top_id"
	F_PARENT_ID    = "parent_id"
	F_IS_MENU      = "is_menu"
)

type Resource struct {
	Id               uint32         `db:"res_id"`
	Name             string         `db:"res_name"`
	Code             string         `db:"res_code"`
	Desc             sql.NullString `db:"res_desc"`
	Url              string         `db:"res_url"`
	TopId            uint32         `db:"top_id"`
	ParentId         uint32         `db:"parent_id"`
	IsMenu           bool           `db:"is_menu"`
	CreatedById      uint32         `db:"created_by_id"`
	CreatedByName    sql.NullString `db:"created_by_name"`
	CreatedTime      mysql.NullTime `db:"created_time"`
	LastModifiedTime mysql.NullTime `db:"last_modified_time"`

	// Transient
	Top      *Resource   `db:"-"`
	Parent   *Resource   `db:"-"`
	Children []*Resource `db:"-"`
	// Relation manage by Role
	Roles []*Role `db:"-"`
}

func (r *Resource) PreInsert(_ gorp.SqlExecutor) error {
	timeNow := time.Now()
	if !r.Desc.Valid {
		r.Desc.String = r.Name
	}
	if !r.CreatedByName.Valid {
		r.CreatedByName.String = "System"
	}
	if r.Top != nil {
		r.TopId = r.Top.Id
	}
	if r.Parent != nil {
		r.ParentId = r.Parent.Id
	}
	if r.TopId <= 0 {
		r.TopId = r.ParentId
	}
	r.CreatedTime = mysql.NullTime{timeNow, true}
	r.LastModifiedTime = mysql.NullTime{timeNow, true}
	return nil
}

func (r *Resource) PreUpdate(_ gorp.SqlExecutor) error {
	r.LastModifiedTime = mysql.NullTime{time.Now(), true}
	return nil
}

func ToResource(i interface{}, err error) *Resource {
	if err != nil {
		panic(err)
	}
	if i == nil || reflect.ValueOf(i).IsNil() {
		return nil
	}
	return i.(*Resource)
}

func ToResources(results []interface{}, err error) []*Resource {
	if err != nil {
		panic(err)
	}
	size := len(results)
	resources := make([]*Resource, size)
	if size == 0 {
		return resources
	}
	for i, r := range results {
		resources[i] = r.(*Resource)
	}
	return resources
}
