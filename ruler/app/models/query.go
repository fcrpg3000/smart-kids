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

const (
	baseQueryForAdmin = "SELECT admin_id, admin_name, hash_password, pwd_salt, " +
		"user_id, user_name, emp_name, emp_no, created_by_id, " +
		"created_by_name, created_time, last_modified_time, " +
		"is_enabled, last_ip FROM m_admin "
	BASE_QUERY_ROLE = "SELECT role_id, role_name, role_code, role_desc, created_by_id, " +
		"created_by_name, created_time, last_modified_time FROM m_role "
	QUERY_ADMIN_BY_NAME = baseQueryForAdmin + "WHERE admin_name = ?"
	QUERY_ROLE_BY_NAME  = BASE_QUERY_ROLE + "WHERE role_name = ?"
	QUERY_ROLE_BY_CODE  = BASE_QUERY_ROLE + "WHERE role_code = ?"
)
