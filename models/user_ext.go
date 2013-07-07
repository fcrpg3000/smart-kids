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
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/go-sql-driver/mysql"
	_ "github.com/robfig/revel"
	"reflect"
	"strings"
	"time"
)

const (
	F_OPERATOR_ID   = "operator_id"
	F_OPERATOR_NAME = "operator_name"
	F_CAUSE         = "banned_cause"
	F_IS_PERMANENT  = "is_permanent"
	F_BANNED_TIME   = "banned_time"
	F_UNBAN_TIME    = "unban_time"
)

var (
	BannedUserFields = strings.Join([]string{
		F_ID, F_USER_ID, F_USER_NAME, F_OPERATOR_ID, F_OPERATOR_NAME,
		F_CAUSE, F_IS_PERMANENT, F_BANNED_TIME, F_UNBAN_TIME,
		F_CREATED_TIME, F_LAST_MODIFIED_TIME,
	}, ", ")
)

// BannedUser struct
// ----------------------------------------------------------------------------

type BannedUser struct {
	Id               uint           `db:"id"`
	UserId           uint64         `db:"user_id"`
	UserName         string         `db:"user_name"`
	OperatorId       int            `db:"operator_id"`
	OperatorName     string         `db:"operator_name"`
	Cause            string         `db:"banned_cause"`
	IsPermanent      bool           `db:"is_permanent"`
	BannedTime       mysql.NullTime `db:"banned_time"`
	UnbanTime        mysql.NullTime `db:"unban_time"`
	CreatedTime      mysql.NullTime `db:"created_time"`
	LastModifiedTime mysql.NullTime `db:"last_modified_time"`
}

// BannedUser instance default string
func (b BannedUser) String() string {
	return fmt.Sprintf("BannedUser{Id=%d,Target=(%d, %s),Operator(%d, %s),"+
		"Cause=\"%s\",Permanent=%v,Period=(%v - %v),"+
		"LastModified=%v}",
		b.Id, b.UserId, b.UserName, b.OperatorId, b.OperatorName,
		b.Cause, b.IsPermanent, b.BannedTime.Time, b.UnbanTime.Time,
		b.LastModifiedTime.Time)
}

func (b *BannedUser) PreInsert(_ gorp.SqlExecutor) error {
	timeNow := time.Now()
	b.CreatedTime = mysql.NullTime{timeNow, true}
	b.LastModifiedTime = mysql.NullTime{timeNow, true}
	return nil
}

func ToBannedUser(i interface{}, err error) *BannedUser {
	if err != nil {
		panic(err)
	}
	if i == nil || reflect.ValueOf(i).IsNil() {
		return nil
	}
	return i.(*BannedUser)
}

func ToBannedUsers(results []interface{}, err error) []*BannedUser {
	if err != nil {
		panic(err)
	}
	size := len(results)
	bannedUsers := make([]*BannedUser, size)
	if size == 0 {
		return bannedUsers
	}
	for i, r := range results {
		bannedUsers[i] = r.(*BannedUser)
	}
	return bannedUsers
}
