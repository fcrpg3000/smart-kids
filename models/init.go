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
	_ "reflect"
	"regexp"
)

// forum module table name constants
const (
	FORUM_TABLE             = "sk_forum"
	FORUM_FIELD_TABLE       = "sk_forum_field"
	FORUM_FIELD_VALUE_TABLE = "sk_forum_field_value"
	FORUM_THREAD_TABLE      = "sk_forum_thread"
)

// model's shared field name constants
const (
	F_ID                 = "id"
	F_ID_ALIAS           = "id_alias"
	F_USER_ID            = "user_id"
	F_USER_NAME          = "user_name"
	F_CREATED_TIME       = "created_time"
	F_LAST_MODIFIED_TIME = "last_modified_time"
	F_STATUS             = "status"
	F_SORT_ORDER         = "sort_order"
	F_TITLE              = "title"
	F_SUMMARY            = "summary"
)

// shared field value constants
const (
	STATUS_DELETED = int16(-1)
)

// request params
const (
	PARAM_CLIENT_ID     = "client_id"
	PARAM_CLIENT_SECRET = "client_secret"
)

var (
	IdAliasRule = regexp.MustCompile("[0-9a-zA-Z][\\w\\-]+")
)
