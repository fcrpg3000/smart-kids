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
	"github.com/robfig/revel"
	"time"
)

// Forum model
type Forum struct {
	Id               uint           `db:"id"`
	IdAlias          string         `db:"id_alias"`
	Title            string         `db:"title"`
	Summary          string         `db:"summary"`
	SortOrder        uint16         `db:"sort_order"`
	Status           int16          `db:"status"`
	CreatedTime      mysql.NullTime `db:"created_time"`
	LastModifiedTime mysql.NullTime `db:"last_modified_time"`
}

func (f Forum) String() string {
	return fmt.Sprintf("Forum{Id=%d, IdAlias=%s, Title=%s, SortOrder=%d, Status=%d, "+
		"CreatedTime=%v, LastModifiedTime=%v, Summary=%s}", f.Id, f.IdAlias, f.Title, f.SortOrder,
		f.Status, f.CreatedTime, f.LastModifiedTime, f.Summary)
}

// Returns true if this Forum is invalid, otherwise false.
func (f *Forum) IsInvalid() bool {
	return f.Status == STATUS_DELETED
}

// Returns true if this forum is valid, otherwise false.
func (f *Forum) IsValid() bool {
	return !f.IsInvalid()
}

// The forum field validate function
func (f *Forum) Validate(v *revel.Validation) {
	v.Check(f.IdAlias,
		revel.Required{},
		revel.MinSize{3},
		revel.MaxSize{50},
		revel.Match{IdAliasRule},
	)
	v.Check(f.Title,
		revel.Required{},
		revel.MinSize{5},
		revel.MaxSize{30},
	)
	v.Check(f.Summary,
		revel.Required{},
		revel.MinSize{1},
		revel.MaxSize{300},
	)
}

// pre-insert hook function
func (f *Forum) PreInsert(_ gorp.SqlExecutor) error {
	timeNow := time.Now()
	f.CreatedTime = mysql.NullTime{timeNow, true}
	f.LastModifiedTime = mysql.NullTime{timeNow, true}
	return nil
}

func NewForum(idAlias, title, summary string) *Forum {
	forum := &Forum{IdAlias: idAlias, Title: title, Summary: summary}
	forum.Status = int16(1)
	forum.SortOrder = uint16(0)
	return forum
}

// Forum's Field model
type ForumField struct {
	Id          uint           `db:"id"`
	ForumId     uint           `db:"forum_id"`
	Name        string         `db:"field_name"`
	Summary     sql.NullString `db:"summary"`
	Rule        sql.NullString `db:"field_rule"`
	FieldTypeId int16          `db:"field_type"`
	SortOrder   uint16         `db:"sort_order"`
	Required    bool           `db:"required"`
	Options     int            `db:"options"`
}

// Forum's Field value model
type ForumFieldValue struct {
	Id        uint           `db:"id"`
	FieldId   uint           `db:"field_id"`  // ForumField.Id
	ParentId  uint           `db:"parent_id"` // >> FieldId
	Name      string         `db:"field_name"`
	Value     sql.NullString `db:"field_value"`
	SortOrder uint16         `db:"sort_order"`
	IsDefault bool           `db:"is_default"`
}

type Thread struct {
	Id               uint64         `db:"id"`
	IdAlias          string         `db:"id_alias"`
	UserId           uint64         `db:"user_id"`
	ForumId          uint           `db:"forum_id"`
	TypeId           int16          `db:"type_id"`
	Title            string         `db:"title"`
	Content          string         `db:"content"`
	Tags             sql.NullString `db:"tags"`
	SourceUrl        sql.NullString `db:"source_url"`
	ViewCount        uint           `db:"view_count"`
	ReplyCount       uint           `db:"reply_count"`
	LastPostId       uint64         `db:"last_post_id"`
	LastPostUserId   uint64         `db:"last_post_user_id"`
	LastPostTime     mysql.NullTime `db:"last_post_time"`
	AsTop            bool           `db:"as_top"`
	AsGood           bool           `db:"as_good"`
	ClientIp         sql.NullString `db:"client_ip"`
	CreatedTime      mysql.NullTime `db:"created_time"`
	LastModifiedTime mysql.NullTime `db:"last_modified_time"`
	Options          int            `db:"options"`
	Status           int16          `db:"status"`
}

type Posts struct {
	Id          uint64         `db:"id"`
	ThreadId    uint64         `db:"thread_id"`
	UserId      uint64         `db:"user_id"`
	UserName    string         `db:"user_name"`
	UserEmail   sql.NullString `db:"user_email"` // redundant field
	UserUrl     sql.NullString `db:"user_url"`   // redundant field
	Title       string         `db:"title"`
	Content     string         `db:"content"`
	ClientIp    string         `db:"client_ip"`
	CreatedTime mysql.NullTime `db:"created_time"`
	Options     int            `db:"options"`
	Status      int16          `db:"status"`
}

type PostsReply struct {
	Id uint64 `db:"id"`
	PostsId uint64 `db:"posts_id"`
	UserId uint64 `db:"user_id"`
	UserName string `db:"user_name"`
	UserEmail   sql.NullString `db:"user_email"` // redundant field
	UserUrl     sql.NullString `db:"user_url"`   // redundant field
	
}