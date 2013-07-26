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
	"log"
	_ "reflect"
	m "smart-kids/ruler/app/models"
	"smart-kids/util"
	"strconv"
)

type Privileges struct {
	Application
}

func (c Privileges) findAllResource(pageable *util.Pageable) *util.Page {
	var (
		total   int64
		content []interface{}
		err     error
	)
	total, err = c.Txn.SelectInt(m.CountSql(m.F_RES_URL, m.RESOURCE_TABLE))
	if err != nil {
		log.Printf("Count for resource error: %s", err.Error())
		panic(err)
	}
	if total == 0 {
		return util.NewPage(nil, pageable, total)
	}
	if pageable == nil {
		content, err = c.Txn.Select(m.Resource{}, m.BASE_QUERY_RESOURCE)
	} else {
		sql := m.NewSqlBuilder(m.BASE_QUERY_RESOURCE).
			PageOrderBy(pageable, util.AscendingSort([]string{m.F_RES_URL})).
			ToSqlString()
		content, err = c.Txn.Select(m.Resource{}, sql)
	}
	if err != nil {
		log.Printf("Select resources error: %s", err.Error())
		panic(err)
	}
	return util.NewPage(content, pageable, total)
}

func (p Privileges) findTopResources() []*m.Resource {
	return m.ToResources(p.Txn.Select(m.Resource{}, m.BASE_QUERY_RESOURCE+
		" WHERE parent_id = 0"))
}

func (p Privileges) addResource(res *m.Resource) error {
	if res.Id > 0 {
		return errors.New("The resource already exists.")
	}
	admin := p.connected()
	res.CreatedById = admin.Id
	res.CreatedByName.Scan(admin.AdminName)
	return p.Txn.Insert(res)
}

func (p Privileges) updateResource(res *m.Resource) (int64, error) {
	if res.Id <= 0 {
		return 0, errors.New(p.NotFoundMessage("权限资源信息"))
	}
	existsRes := p.loadResource(res.Id)
	if existsRes == nil {
		return 0, errors.New(p.NotFoundMessage("权限资源信息"))
	}
	existsRes.Name = res.Name
	existsRes.Code = res.Code
	existsRes.Url = res.Url
	existsRes.Desc = res.Desc
	existsRes.ParentId = res.ParentId
	existsRes.IsMenu = res.IsMenu
	return p.Txn.Update(existsRes)
}

func (p Privileges) loadResource(id uint32) *m.Resource {
	return m.ToResource(p.Txn.Get(m.Resource{}, id))
}

func (c Privileges) findAllRoles(pageable *util.Pageable) *util.Page {
	var (
		total   int64
		content []interface{}
		err     error
	)
	total, err = c.Txn.SelectInt(m.CountSql(m.F_ROLE_NAME, m.ROLE_TABLE))
	if err != nil {
		log.Printf("Count for role error: %s", err.Error())
		panic(err)
	}
	if total == 0 {
		return util.NewPage(nil, pageable, total)
	}
	if pageable == nil {
		content, err = c.Txn.Select(m.Role{}, m.BASE_QUERY_ROLE)
	} else {
		sql := m.NewSqlBuilder(m.BASE_QUERY_ROLE).
			PageOrderBy(pageable, util.AscendingSort([]string{m.F_ROLE_CODE})).
			ToSqlString()
		content, err = c.Txn.Select(m.Role{}, sql)
	}
	if err != nil {
		log.Printf("Select roles error: %s", err.Error())
		panic(err)
	}
	return util.NewPage(content, pageable, total)
}

// Pagination resources
func (c Privileges) ResourceList(p, ps int) revel.Result {
	if ps <= 0 {
		ps = DEFAULT_PAGE_SIZE
	}
	sort := util.AscendingSort([]string{m.F_RES_URL, m.F_RES_NAME})
	pageable, err := util.NewPageable0(p, ps, sort)
	if err != nil {
		panic(err)
	}
	title := c.Message("resource.title.list")
	pageResource := c.findAllResource(pageable)
	return c.Render(title, pageResource)
}

// Resource edit page
func (p Privileges) ResourceEdit(id uint32) revel.Result {
	topResources := p.findTopResources()
	title := p.Message("resource.title.creation")
	if id <= 0 {
		return p.Render(title, topResources)
	}
	res := p.loadResource(id)
	if res == nil {
		return p.Render(title, topResources)
	}
	if res.ParentId > 0 {
		res.Parent = p.loadResource(res.ParentId)
	}
	title = p.Message("resource.title.edit", res.Name)
	return p.Render(title, res, topResources)
}

// Insert or Update Resource (maybe ajax post request, so return json data)
func (p Privileges) SaveResource(res m.Resource) revel.Result {
	var (
		err    error
		row    int64
		result *util.ResponseResult
		update bool
	)
	if res.Id > 0 { // Update
		update = true
		row, err = p.updateResource(&res)
	} else {
		update = false
		err = p.addResource(&res)
	}
	if err != nil {
		result = util.ErrorResult(p.Message("resource.error.saved", err.Error()))
	} else {
		if row <= 0 && update {
			result = util.FailureResult(p.NotFoundMessage("权限资源信息"))
		} else {
			result = util.SuccessResult(p.Message("resource.success.saved", res.Name))
		}
	}
	return p.RenderJson(result)
}

func (p Privileges) DeleteResource(id uint32) revel.Result {
	_, err := p.Txn.Exec("delete rr.* from `m_role_resource` rr where rr.`res_id` = ?", id)
	//_, err = result.RowsAffected()
	if err != nil {
		panic(err)
		return p.RenderJson(util.ErrorResult(err.Error()))
	}
	_, err = p.Txn.Exec("delete r.* from `m_resource` r where r.`res_id` = ?", id)
	//_, err = result.RowsAffected()
	if err != nil {
		panic(err)
		return p.RenderJson(util.ErrorResult(err.Error()))
	}
	return p.RenderJson(util.SuccessResult(p.OperOkMessage()))
}

// query role pagination
func (c Privileges) RoleList(p, ps int) revel.Result {
	if ps <= 0 {
		ps = DEFAULT_PAGE_SIZE
	}
	sort := util.AscendingSort([]string{m.F_ROLE_CODE})
	pageable, err := util.NewPageable0(p, ps, sort)
	if err != nil {
		panic(err)
	}
	title := c.Message("role.title.list")
	pageRole := c.findAllRoles(pageable)
	return c.Render(title, pageRole)
}

func (p Privileges) RoleFormView(id uint32) revel.Result {
	var (
		title string
		role  *m.Role
	)

	if id == 0 {
		title = p.Message("role.title.addition")
	} else {
		title = p.Message("role.title.modification")
		role = m.ToRole(p.Txn.Get(m.Role{}, id))
	}
	return p.Render(title, role)
}

// add or modify role
func (p Privileges) SaveRole(role m.Role) revel.Result {
	var (
		err    error
		result *util.ResponseResult
	)
	if role.Id == 0 {
		err = p.Txn.Insert((&role).CreatedBy(p.connected()))
	} else {
		if exists := m.ToRole(p.Txn.Get(m.Role{}, role.Id)); exists != nil {
			_, err = p.Txn.Update(exists.UpdateBy(&role))
		}
	}
	if err != nil {
		result = util.ErrorResult(err.Error())
	} else {
		result = util.SuccessResult("保存成功")
	}
	return p.RenderJson(result)
}

func (p Privileges) DeleteRole(id uint32) revel.Result {
	_, err := p.Txn.Exec("delete ar.* from `m_admin_role` ar where ar.`role_id` = ?", id)
	//_, err = result.RowsAffected()
	if err != nil {
		panic(err)
		return p.RenderJson(util.ErrorResult(err.Error()))
	}
	_, err = p.Txn.Exec("delete rr.* from `m_role_resource` rr where rr.`role_id` = ?", id)
	//_, err = result.RowsAffected()
	if err != nil {
		panic(err)
		return p.RenderJson(util.ErrorResult(err.Error()))
	}
	_, err = p.Txn.Exec("delete r.* from `m_role` r where r.`role_id` = ?", id)
	if err != nil {
		panic(err)
		return p.RenderJson(util.ErrorResult(err.Error()))
	}
	return p.RenderJson(util.SuccessResult(p.OperOkMessage()))
}

func (p Privileges) ResourcePrivileges(rid uint32) revel.Result {
	role := m.ToRole(p.Txn.Get(m.Role{}, rid))
	mainResources := m.ToResources(p.Txn.Select(m.Resource{},
		m.QUERY_RESOURCE_BY_PID, 0))
	allSubResources := m.ToResources(p.Txn.Select(m.Resource{},
		m.BASE_QUERY_RESOURCE+"WHERE parent_id > 0"))
	mainResMap := make(map[uint32]int)
	for i, mainRes := range mainResources {
		mainResMap[mainRes.Id] = i
	}
	for _, subRes := range allSubResources {
		idx, ok := mainResMap[subRes.ParentId]
		if !ok {
			continue
		}
		mainRes := mainResources[idx]
		mainRes.Children = append(mainRes.Children, subRes)
	}
	var resIds []uint32
	rows, err := Dbm.Db.Query("select res_id from `m_role_resource` where role_id = ?", rid)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var resId uint32
		if err := rows.Scan(&resId); err != nil {
			panic(err)
		}
		resIds = append(resIds, resId)
	}
	if err := rows.Err(); err != nil {
		panic(err)
	}
	return p.Render(role, resIds, mainResources)
}

func (p Privileges) SavePrivileges() revel.Result {
	var allResIds []uint32
	roleId, err := strconv.ParseUint(p.Request.Form.Get("roleId"), 10, 32)
	if err != nil {
		panic(err)
	}
	mainIds := p.Request.Form["mainResource"]
	if len(mainIds) > 0 {
		for _, mainId := range mainIds {
			mid, err := strconv.ParseUint(mainId, 10, 32)
			if err != nil {
				panic(err)
			}
			allResIds = append(allResIds, uint32(mid))
			subIds := p.Request.Form[fmt.Sprintf("subResource%d", mid)]
			for _, subId := range subIds {
				sid, err := strconv.ParseUint(subId, 10, 32)
				if err != nil {
					panic(err)
				}
				allResIds = append(allResIds, uint32(sid))
			}
		}
	}
	_, err = p.Txn.Exec("delete rr.* from `m_role_resource` rr where rr.`role_id` = ?",
		uint32(roleId))
	if err != nil {
		panic(err)
	}
	for _, resId := range allResIds {
		_, err = p.Txn.Exec("insert into `m_role_resource` values(?, ?)",
			roleId, resId)
		if err != nil {
			panic(err)
		}
	}
	return p.RenderJson(util.SuccessResult("OK"))
}
