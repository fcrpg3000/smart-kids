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
	"time"
)

type Admin struct {
	Id               int            `db:"admin_id"`
	AdminName        string         `db:"admin_name"`
	HashPassword     string         `db:"hash_password"`
	Salt             string         `db:"pwd_salt"`
	UserId           int64          `db:"user_id"`   // relate to User#Id if necessary
	UserName         sql.NullString `db:"user_name"` // relate to User#UserName if necessary
	EmpName          sql.NullString `db:"emp_name"`
	EmpNo            sql.NullString `db:"emp_no"`
	CreatedById      int            `db:"created_by_id"`
	CreatedByName    sql.NullString `db:"created_by_name"`
	CreatedTime      mysql.NullTime `db:"created_time"`
	LastModifiedTime mysql.NullTime `db:"last_modified_time"`
	LastIp           sql.NullString `db:"last_ip"`

	// Transient
	Roles []*Role `db:"-"`
}

// gorp 不能自动处理关联关系，实现这个接口方法，当调用根据Id获取 Admin 对象时，
// 自动获取关联的角色信息列表
func (a *Admin) PostGet(exe gorp.SqlExecutor) error {
	query := "SELECT role_id, role_name, role_code, role_desc, created_by_id, " +
		"created_by_name, created_time, last_modified_time " +
		"FROM m_role WHERE role_id IN " +
		"(SELECT role_id FROM m_admin_role ar JOIN m_admin a " +
		"ON a.admin_id = ar.admin_id " +
		"WHERE a.admin_id = ?)"
	objs, err := exe.Select(Role{}, query, a.Id)
	if err != nil {
		return fmt.Errorf("Error loading admin's(%d) roles : %s", a.Id, err)
	}
	if len(objs) > 0 {
		a.Roles = make([]*Role, 0)
		for _, obj := range objs {
			a.Roles = append(a.Roles, obj.(*Role))
		}
	}
	return nil
}

// mapped table
type Role struct {
	Id               int            `db:"role_id"`
	Code             string         `db:"role_code"`
	Name             string         `db:"role_name"`
	Desc             sql.NullString `db:"role_desc"`
	CreatedById      int            `db:"created_by_id"`
	CreatedByName    sql.NullString `db:"created_by_name"`
	CreatedTime      mysql.NullTime `db:"created_time"`
	LastModifiedTime mysql.NullTime `db:"last_modified_time"`

	// Transient
	// Relation manage by Role
	Resources []*Resource `db:"-"`
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

type Resource struct {
	Id               int            `db:"res_id"`
	Name             string         `db:"res_name"`
	Code             string         `db:"res_code"`
	Desc             string         `db:"res_desc"`
	Url              string         `db:"res_url"`
	ParentId         int            `db:"parent_id"`
	IsMenu           bool           `db:"is_menu"`
	CreatedById      int            `db:"created_by_id"`
	CreatedByName    string         `db:"created_by_name"`
	CreatedTime      mysql.NullTime `db:"created_time"`
	LastModifiedTime mysql.NullTime `db:"last_modified_time"`

	// Transient
	Parent   *Resource   `db:"-"`
	Children []*Resource `db:"-"`
	// Relation manage by Role
	Roles []*Role `db:"-"`
}

func (r *Resource) PreInsert(_ gorp.SqlExecutor) error {
	timeNow := time.Now()
	r.CreatedTime = mysql.NullTime{timeNow, true}
	r.LastModifiedTime = mysql.NullTime{timeNow, true}
	return nil
}

func (r *Resource) PreUpdate(_ gorp.SqlExecutor) error {
	r.LastModifiedTime = mysql.NullTime{time.Now(), true}
	return nil
}
