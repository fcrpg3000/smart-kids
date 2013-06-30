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
	"errors"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/go-sql-driver/mysql"
	"time"
)

type Posts struct {
	Id               int64          `db:"posts_id"`
	SrcId            int64          `db:"posts_src_id"`
	UserId           int64          `db:"user_id"`
	UserName         string         `db:"user_name"`
	Content          string         `db:"posts_content"`
	ClientType       int            `db:"client_type"`
	CreatedTime      mysql.NullTime `db:"created_time"`
	LastModifiedTime mysql.NullTime `db:"last_modified_time"`

	// Transient
	Client *Client `db:"-" json:"client,omitempty"`
	User   *User   `db:"-" json:"user,omitempty"`
}

// Posts description string
func (p *Posts) Description() string {
	return fmt.Sprintf("%s: %s %s From %s",
		p.UserName, p.Content, p.CreatedTime.Time.Format(time.RFC3339),
		p.Client.Name)
}

var (
	negativePostsId = errors.New("The given source Posts Id must not be 0 or positive.")
	emptyContent    = errors.New("The posts content must not be nil or empty string.")
	nullUser        = errors.New("The given user Ptr must not be nil.")
	nullClient      = errors.New("The given client Ptr must not be nil.")
)

// Create a new Posts base on source Posts id, current user, client and posts content
// Returns error if srcId <= 0 || user == nil || client == nil
func NewForward(srcId int64, user *User, client *Client, content string) (*Posts, error) {
	if srcId <= 0 {
		return nil, negativePostsId
	}
	return newPostsInternal(srcId, user, client, content)
}

// Creates and Returns a new Posts Ptr base on current User, client and posts content.
// Returns error if user == nil || client == nil || len(content) == 0
func NewPosts(user *User, client *Client, content string) (*Posts, error) {
	return newPostsInternal(int64(-1), user, client, content)
}

func newPostsInternal(srcId int64, user *User, client *Client, content string) (*Posts, error) {
	if user == nil {
		return nil, nullUser
	}
	if client == nil {
		return nil, nullClient
	}
	if len(content) == 0 {
		return nil, emptyContent
	}
	posts := &Posts{}
	timeNow := time.Now()
	posts.SrcId = srcId
	posts.Client = client
	posts.User = user
	posts.Content = content
	posts.CreatedTime = mysql.NullTime{timeNow, true}
	posts.LastModifiedTime = mysql.NullTime{timeNow, true}
	return posts, nil
}

// These hooks work around one thing:
// - Gorp's lack of support for loading relations automatically.
func (p *Posts) PreInsert(_ gorp.SqlExecutor) error {
	if p.Client != nil {
		p.ClientType = p.Client.Id
	}

	if p.User != nil {
		p.UserId = p.User.UserId
		p.UserName = p.User.UserName
	}
	return nil
}

func (p *Posts) PostGet() error {
	p.Client = ClientOf(p.ClientType)

	return nil
}
