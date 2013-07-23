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
	_ "fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type PhotoAlbum struct {
	Id          uint64         `db:"id"`
	Name        string         `db:"album_name"`
	UserId      uint64         `db:"user_id"`
	Tags        sql.NullString `db:"tags"`
	PhotoCount  uint           `db:"photo_count"`
	ViewCount   uint           `db:"view_count"`
	FrontCover  uint64         `db:"front_cover"`
	VCode       uint16         `db:"v_code"`
	CreatedTime time.Time      `db:"created_time"`
}

func (p *PhotoAlbum) PreInsert(_ gorp.SqlExecutor) error {
	p.CreatedTime = time.Now()
	return nil
}

type Photo struct {
	Id               uint64         `db:"id"`
	UserId           uint64         `db:"user_id"`
	AlbumId          uint64         `db:"album_id"`
	Tags             sql.NullString `db:"tags"`
	Description      sql.NullString `db:"description"`
	SourceUrl        string         `db:"source_url"`
	MediumUrl        string         `db:"medium_url"`
	SmallUrl         string         `db:"small_url"`
	ThumbUrl         string         `db:"thumb_url"`
	ViewCount        uint           `db:"view_count"`
	CommentCount     uint           `db:"comment_count"`
	CreatedTime      time.Time      `db:"created_time"`
	LastModifiedTime time.Time      `db:"last_modified_time"`
}

type PhotoComment struct {
	BaseComment
	PhotoId uint64 `db:"photo_id"`
}
