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
	m "smart-kids/models"
	"smart-kids/query"
	"smart-kids/util"
)

var (
	appListSql = query.SimpleQuerySql(m.AppFields, m.APP_TABLE, "x")
)

type AppController struct {
	Application
}

func (a AppController) findPageApp(pageable *util.Pageable) *util.Page {
	var (
		total   int64
		content []interface{}
		err     error
	)
	total, err = a.Txn.SelectInt(query.CountSql(m.F_APP_NAME, m.APP_TABLE))
	if total == 0 || err != nil {
		return util.NewPage(nil, pageable, total)
	}
	if pageable == nil {
		content, err = a.Txn.Select(m.App{}, appListSql)
	} else {
		sql := query.NewSqlBuilder(appListSql).
			PageOrderBy(pageable, util.AscendingSort([]string{m.F_APP_NAME})).
			ToSqlString()
		content, err = a.Txn.Select(m.App{}, sql)
	}
	if err != nil {
		panic(err)
	}
	return util.NewPage(content, pageable, total)
}

// App models pagination
func (a AppController) AppList(p, ps int) revel.Result {
	if ps <= 1 {
		ps = DEFAULT_PAGE_SIZE
	}
	if p <= 0 {
		p = 1
	}
	pageable, err := util.NewPageable(p, ps, util.ASC, []string{m.F_APP_NAME})
	if err != nil { // never heppen
		log.Fatalf("Error for %s", err.Error())
		panic(err)
	}
	pageApp := a.findPageApp(pageable)
	title := a.Message("AppList.title.list")
	return a.Render(title, pageApp)
}
