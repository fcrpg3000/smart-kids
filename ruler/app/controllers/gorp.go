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
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/revel"
	"smart-kids/ruler/app/models"
)

var (
	Dbm *gorp.DbMap
)

func Init() {
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	initAdmin()

	Dbm.TraceOn("[gorp]", revel.INFO)

}

func setColumnSizes(t *gorp.TableMap, colSizes map[string]int) {
	for col, size := range colSizes {
		t.ColMap(col).MaxSize = size
	}
}

func initAdmin() {
	// Register Admin
	t := Dbm.AddTableWithName(models.Admin{}, "m_admin").SetKeys(true, "Id")
	setColumnSizes(t, map[string]int{
		"AdminName":     50,
		"HashPassword":  100,
		"Salt":          100,
		"UserName":      50,
		"EmpName":       50,
		"EmpNo":         30,
		"CreatedByName": 50,
	})

	// Register Role
	t = Dbm.AddTableWithName(models.Role{}, "m_role").SetKeys(true, "Id")
	setColumnSizes(t, map[string]int{
		"Code":          50,
		"Name":          50,
		"Desc":          200,
		"CreatedByName": 50,
	})

	// Register Resource
	t = Dbm.AddTableWithName(models.Resource{}, "m_resource").SetKeys(true, "Id")
	setColumnSizes(t, map[string]int{
		"Code":          50,
		"Name":          50,
		"Desc":          200,
		"Url":           255,
		"CreatedByName": 50,
	})
}

type GorpController struct {
	*revel.Controller
	Txn *gorp.Transaction
}

func (c *GorpController) Begin() revel.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() revel.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() revel.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}
