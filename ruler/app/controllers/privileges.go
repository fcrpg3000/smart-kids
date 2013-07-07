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
	"errors"
	"fmt"
	"github.com/robfig/revel"
	"log"
	_ "reflect"
	m "smart-kids/ruler/app/models"
	"smart-kids/util"
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
		" WHERE parent_id <= 0"))
}

func (p Privileges) addResource(res *m.Resource) error {
	if res.Id > 0 {
		return errors.New("The resource already exists.")
	}
	admin := p.connected()
	res.CreatedById = admin.Id
	res.CreatedByName = sql.NullString{admin.AdminName, true}
	fmt.Println("New resource: ", res)
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
	fmt.Println("Updated resource: ", res)
	return p.Txn.Update(existsRes)
}

func (p Privileges) loadResource(id int) *m.Resource {
	return m.ToResource(p.Txn.Get(m.Resource{}, id))
}

// Pagination resources
func (c Privileges) ResourceList(p, ps int) revel.Result {
	if ps <= 0 {
		ps = DEFAULT_PAGE_SIZE
	}
	sort := util.AscendingSort([]string{m.F_RES_URL})
	pageable, err := util.NewPageable0(p, ps, sort)
	if err != nil {
		panic(err)
	}
	title := c.Message("resource.title.list")
	pageResource := c.findAllResource(pageable)
	return c.Render(title, pageResource)
}

// Resource edit page
func (p Privileges) ResourceEdit(id int) revel.Result {
	topResources := p.findTopResources()
	title := p.Message("ResourceEdit.title.creation")
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
	title = p.Message("ResourceEdit.title.edit", res.Name)
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
		result = util.ErrorResult(p.Message("resource.errorEdit", err.Error()))
	} else {
		if row <= 0 && update {
			result = util.FailureResult(p.NotFoundMessage("权限资源信息"))
		} else {
			result = util.SuccessResult(p.Message("resource.successEdit", res.Name))
		}
	}
	return p.RenderJson(result)
}
