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
	"log"
	_ "reflect"
	"smart-kids/ruler/app/models"
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
	total, err = c.Txn.SelectInt(models.CountSql("res_url", "m_resource"))
	if err != nil {
		log.Printf("Count for resource error: %s", err.Error())
		panic(err)
	}
	if total == 0 {
		return util.NewPage(nil, pageable, total)
	}
	if pageable == nil {
		content, err = c.Txn.Select(models.Resource{}, models.BASE_QUERY_RESOURCE)
	} else {
		sqlBuilder := models.NewSqlBuilder(models.BASE_QUERY_RESOURCE)
		defaultSort, _ := util.AscendingSort([]string{"res_url"})
		models.PageOrderBy(sqlBuilder, pageable, defaultSort)
		content, err = c.Txn.Select(models.Resource{}, sqlBuilder.ToSqlString())
	}
	if err != nil {
		log.Printf("Select resources error: %s", err.Error())
		panic(err)
	}
	return util.NewPage(content, pageable, total)
}

// Pagination resources
func (c Privileges) ResourceList(p, ps int) revel.Result {
	if ps <= 0 {
		ps = DEFAULT_PAGE_SIZE
	}
	sort, _ := util.AscendingSort([]string{"res_url"})
	pageable, err := util.NewPageable0(p, ps, sort)
	if err != nil {
		panic(err)
	}
	pageResource := c.findAllResource(pageable)
	return c.Render(pageResource)
}
