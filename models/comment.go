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
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type BaseComment struct {
	Id          uint64         `db:"id"`
	UserId      uint64         `db:"user_id"`
	UserName    string         `db:"user_name"`
	UserEmail   string         `db:"user_email"`
	UserUrl     string         `db:"user_url"`
	ClientCode  uint16         `db:"client_code"`
	ClientIp    string         `db:"client_ip"`
	Title       sql.NullString `db:"title"`
	Content     string         `db:"content"`
	IsTop       bool           `db:"is_top"`
	CreatedTime time.Time      `db:"created_time"`

	Client *Client `db:"-"`
}

func (b *BaseComment) AsTop() *BaseComment {
	b.IsTop = true
	return b
}

func (b *BaseComment) PostGet(_ gorp.SqlExecutor) error {
	b.Client = ClientOf(b.ClientCode, UNKNOWN_CLIENT)
	return nil
}
