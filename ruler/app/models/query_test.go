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
	"smart-kids/util"
	"strings"
	"testing"
)

func assertTrue(t *testing.T, expr bool, arg1, arg2 string) {
	if expr {
		t.Log(arg1)
	} else {
		t.Error(arg2)
	}
}

func TestExistsQueryStringCorrectly(t *testing.T) {
	tableName := "user_profile"
	placeHolder := "user_name"
	idAttributes := []string{"is_enabled", "user_level"}
	query := ExistsQueryString(tableName, placeHolder, idAttributes)
	actualQuery := "SELECT count(user_name) FROM user_profile x " +
		"WHERE x.is_enabled = ? AND x.user_level = ? AND 1 = 1"
	if query != actualQuery {
		t.Errorf("ExistsQueryString error:\nreturn: %s\nacutal: %s", query, actualQuery)
	} else {
		t.Log("PASS")
	}
}

// ApplySortingWithAlias function test case when existing order by clauses.
func TestExtendsExistingOrderByClausesCorrectly(t *testing.T) {
	order := util.DescendingOrder("created_time")
	sort := &util.Sort{[]*util.Order{order}}
	query := "SELECT * FROM m_resource r order by r.res_name asc"
	extendOrderQuery := ApplySortingWithAlias(query, sort, "r")
	if strings.Contains(extendOrderQuery, "order by r.res_name asc, r.created_time desc") {
		t.Log("PASS")
	} else {
		t.Error("Extends order by clause error, Actual: ", extendOrderQuery)
	}
}

// ApplySortingWithAlias function test case to use ignoreCase true.
func TestIgnoreCaseOrderingCorrectly(t *testing.T) {
	order := util.NewOrder(util.ASC, "res_name", true)
	sort := &util.Sort{[]*util.Order{order}}
	query := "SELECT * FROM m_resource r"
	orderQuery := ApplySortingWithAlias(query, sort, "r")
	if strings.Contains(orderQuery, "order by lower(r.res_name) asc") {
		t.Log("PASS")
	} else {
		t.Error("OrderBy clause error, Actual: ", orderQuery)
	}
}

// getOuterJoinAliases function test case
func TestGetOuterJoinAliases(t *testing.T) {
	aliases := getOuterJoinAliases("select * from users u left outer join foo b2_$ar where ...")
	assertTrue(t, len(aliases) == 1, "PASS",
		fmt.Sprintf("Aliases length actual: %d", len(aliases)))
	assertTrue(t, aliases[0] == "b2_$ar", "PASS", "First alias actual: "+aliases[0])

	aliases = getOuterJoinAliases("select * from users u left join foo b2_$ar where ...")
	assertTrue(t, len(aliases) == 1, "PASS",
		fmt.Sprintf("Aliases length actual: %d", len(aliases)))
	assertTrue(t, aliases[0] == "b2_$ar", "PASS", "First alias actual: "+aliases[0])

	aliases = getOuterJoinAliases("select * from users u left join foo as b2_$ar, " +
		"left join bar as foo where ...")
	assertTrue(t, len(aliases) == 2, "PASS",
		fmt.Sprintf("Aliases length actual: %d", len(aliases)))
	assertTrue(t, aliases[0] == "b2_$ar", "PASS", "First alias is "+aliases[0])
	assertTrue(t, aliases[1] == "foo", "PASS", "Second alias is "+aliases[1])
}

func TestDetectAliasCorrectly(t *testing.T) {
	if alias, ok := DetectAlias("SELECT * FROM USER_PROFILE U"); ok {
		t.Log("Find alias: PASS")
		if alias != "U" {
			t.Errorf("Matcher Alias fail, target: U Actual: %s", alias)
		} else {
			t.Log("Matches alias: %s, PASS", alias)
		}
	}
	if alias, ok := DetectAlias("select * from  user_profile_status ups"); ok {
		t.Log("Find alias: PASS")
		if alias != "ups" {
			t.Errorf("Matches alias fail, target: ups Actual: %s", alias)
		} else {
			t.Logf("Matches alias: %s, PASS", alias)
		}
	}
}

func TestCreateCountQueryCorrectly(t *testing.T) {
	originalQuery := "select * from users u"
	actual := CreateCountQuery(originalQuery)
	target := "select count(*) from users u"
	if target != actual {
		t.Errorf("Creates count query error, actual: %s, target: %s", actual, target)
	} else {
		t.Log("PASS")
	}
}

func TestCreateCountQueryForCapitalLetterSQL(t *testing.T) {
	query := "SELECT u.user_name FROM users u WHERE u.is_enabled = ?"
	actual := CreateCountQuery(query)
	target := "SELECT count(u.user_name) FROM users u WHERE u.is_enabled = ?"
	if target != actual {
		t.Errorf("Creates count query error, actual: %s, target: %s", actual, target)
	} else {
		t.Log("PASS")
	}
}
