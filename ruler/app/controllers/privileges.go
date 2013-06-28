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
	_ "fmt"
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
	return loadResources(p.Txn.Select(m.Resource{}, m.BASE_QUERY_RESOURCE+
		" WHERE parent_id <= 0"))
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

func (p Privileges) loadResource(id int) *m.Resource {
	obj, err := p.Txn.Get(m.Resource{}, id)
	if err != nil {
		panic(err)
	}
	if obj == nil {
		return nil
	}
	return obj.(*m.Resource)
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
	pageResource := c.findAllResource(pageable)
	return c.Render(pageResource)
}

// Resource edit page
func (p Privileges) ResourceEdit(id int) revel.Result {
	var topResources []*m.Resource
	title := p.Message("ResourceEdit.title.creation")
	if id <= 0 {
		topResources = p.findTopResources()
		return p.Render(title, topResources)
	}
	res := p.loadResource(id)
	if res == nil {
		topResources = p.findTopResources()
		return p.Render(title, topResources)
	}
	if res.ParentId > 0 {
		res.Parent = p.loadResource(res.ParentId)
		topResources = p.findTopResources()
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
	)
	if res.Id > 0 { // Update
		row, err = p.updateResource(&res)
	} else {
		err = p.Txn.Insert(&res)
	}
	if err != nil {
		result = util.ErrorResult(p.Message("resource.errorEdit", err.Error()))
	} else {
		if row <= 0 {
			result = util.FailureResult(p.NotFoundMessage("权限资源信息"))
		} else {
			result = util.SuccessResult(p.Message("resource.successEdit", res.Name))
		}
	}
	return p.RenderJson(result)
}
