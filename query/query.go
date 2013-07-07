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
package query

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"smart-kids/util"
	"strings"
)

// Private const
const (
	countQueryTpl       = "SELECT count(%s) FROM %s %s"
	simpleQueryTpl      = "SELECT %s FROM %s %s"
	deleteAllQueryTpl   = "DELETE FROM %s x "
	defaultAlias        = "x"
	countReplacementTpl = "SELECT count(%s) $5$6$7"
	simpleCountValue    = "$2"
	complexCountValue   = "$3$6"
	identifier          = "[\\w._$]+"
	equalsConditionTpl  = "%s.%s = ?"
)

// query variables
var (
	identifierGroup  = fmt.Sprintf("(%s)", identifier)
	leftJoin         = fmt.Sprintf("left (outer )?join %s (as )?%s", identifier, identifierGroup)
	leftJoinRegexp   = regexp.MustCompile(leftJoin)
	aliasMatchRegexp = regexp.MustCompile(fmt.Sprintf("(?: from)(?: )+%s(?: as)*(?: )+(\\w*)", identifierGroup))
	countMatchRegexp = regexp.MustCompile(fmt.Sprintf("((?i:select)\\s+((distinct )?(.+?)?)\\s+)?((?i:from)\\s+"+
		"%s(?:\\s+as)?\\s+)%s(.*)", identifier, identifierGroup))
	fromMatchRegexp = regexp.MustCompile("(?i:from)")
)

// error variables
var (
	defaultSortError = errors.New("The pageable's sort is nil, the default sort must not be nil.")
)

// inner set implementation
type stringSet struct {
	a []string
}

func (s stringSet) contains(value string) bool {
	return s.indexOf(value) != -1
}

func (s stringSet) indexOf(value string) int {
	for i, entry := range s.a {
		if value == entry {
			return i
		}
	}
	return -1
}

func (s *stringSet) add(value string) *stringSet {
	if s.contains(value) {
		return s
	}
	s.a = append(s.a, value)
	return s
}

func (s stringSet) toArray() []string {
	return s.a
}

// Query string utilities functions
// --------------------------------------------------------------------------------------

func SimpleQuerySql(fields, table, alias string) string {
	return fmt.Sprintf(simpleQueryTpl, fields, table, alias)
}

// Returns the query string to execute an exists query for the given id attributes.
// Examples:
//    tableName, placeHolder := "user", "username"
//    ExistsQueryString(tableName, placeHolder, []string{"is_vip"})
// Returns:
//    "SELECT count(username) FROM user x WHERE x.is_vip = ? AND 1 = 1"
//
// Param tableName the name of the table to create the query for, must not be empty.
// Param placeHolder the placeholder for the count clause, must not be empty.
// Param idAttributes the id attributes for the table.
func ExistsQueryString(tableName, placeHolder string, idAttributes []string) string {
	querySlice := []string{CountSql(placeHolder, tableName)}
	querySlice = append(querySlice, " WHERE ")

	for _, attribute := range idAttributes {
		querySlice = append(querySlice, fmt.Sprintf(equalsConditionTpl, "x", attribute))
		querySlice = append(querySlice, " AND ")
	}
	querySlice = append(querySlice, "1 = 1")
	return strings.Join(querySlice, "")
}

// Adds "order by" clause to the SQL query.
// Uses the default alias to bind the sorting property to.
func ApplySorting(query string, sort *util.Sort) string {
	return ApplySortingWithAlias(query, sort, defaultAlias)
}

// Adds "order by" clause to the SQL query.
func ApplySortingWithAlias(query string, sort *util.Sort, alias string) string {
	if sort == nil || reflect.ValueOf(sort).IsNil() || len(sort.Orders) == 0 {
		return query
	}
	q := []string{query}
	if strings.Contains(query, "order by") {
		q = append(q, ", ")
	} else {
		q = append(q, " order by ")
	}
	aliases := getOuterJoinAliases(query)

	for i, order := range sort.Orders {
		if i > 0 {
			q = append(q, ", ")
		}
		q = append(q, getOrderClause(aliases, alias, order))
	}
	return strings.Join(q, "")
}

// Resolves the alias for the entity to be retrieved from the given SQL query.
func DetectAlias(query string) (string, bool) {
	submatches := aliasMatchRegexp.FindAllStringSubmatch(query, -1)
	if len(submatches) == 0 {
		return "", false
	}
	submatch := submatches[0]
	if len(submatch) < 3 {
		return "", false
	}
	target := submatch[2]
	if len(target) == 0 {
		return "", false
	} else {
		return target, true
	}
}

// Returns count query SQL string of the specified count field and table name.
func CountSql(f string, tbl string) string {
	return simpleCountQuery(f, tbl, defaultAlias)
}

func simpleCountQuery(f, tbl, alias string) string {
	countQuery := fmt.Sprintf(countQueryTpl, f, tbl, alias)
	if len(alias) == 0 {
		return countQuery[0 : len(countQuery)-1]
	}
	return countQuery
}

// Creates a count projected query from the given orginal query.
// Returns empty string if originalQuery is empty.
//
// Param originalQuery must not be empty
func CreateCountQuery(originalQuery string) string {
	if len(originalQuery) == 0 {
		return ""
	}
	matches := countMatchRegexp.FindStringSubmatch(originalQuery)
	var (
		variable         string
		countReplacement string
		useVariable      bool
		multiFields      bool
		notMatch         bool
	)
	notMatch = len(matches) == 0
	if !notMatch {
		variable = matches[4]
		// for i, match := range matches {
		// 	fmt.Printf("$_%d = > %s\n", i, match)
		// }
		if !strings.Contains(variable, "distinct") && strings.Contains(variable, ",") {
			multiFields = true
		}
	}
	if multiFields || notMatch {
		fromIdx := fromMatchRegexp.FindStringIndex(originalQuery)
		if len(fromIdx) == 0 {
			return originalQuery
		}
		return simpleCountQuery("*", originalQuery[fromIdx[0]+5:], "")
	}

	useVariable = len(variable) > 0 && strings.Index(variable, "new") != 0 &&
		strings.Index(variable, "count(") != 0
	if useVariable {
		countReplacement = fmt.Sprintf(countReplacementTpl, simpleCountValue)
	} else {
		countReplacement = fmt.Sprintf(countReplacementTpl, complexCountValue)
	}
	return countMatchRegexp.ReplaceAllString(originalQuery, countReplacement)
}

func getOuterJoinAliases(query string) []string {
	result := stringSet{}
	matchers := leftJoinRegexp.FindAllStringSubmatch(query, -1)

	for _, matcher := range matchers {
		if len(matcher) == 0 {
			continue
		}
		if len(matcher) > 3 {
			group := matcher[3]
			if len(group) > 0 {
				result.add(group)
			}
		}
	}
	return result.toArray()
}

// Returns the order clause for the given Order.
// Will prefix the clause with the given alias if the referenced
// property refers to a join alias.
//
// Param joinAliases the join aliases of the original query.
// Param alias the alias for the root entity.
// Param order the order object to build the clause for.
func getOrderClause(joinAliases []string, alias string, order *util.Order) string {
	qualifyReference := true

	for _, joinAlias := range joinAliases {
		if strings.Index(joinAlias, order.Property) == 0 {
			qualifyReference = false
			break
		}
	}
	var (
		reference string
		wrapped   string
	)
	if qualifyReference {
		reference = fmt.Sprintf("%s.%s", alias, order.Property)
	} else {
		reference = order.Property
	}
	if order.IgnoreCase {
		wrapped = fmt.Sprintf("lower(%s)", reference)
	} else {
		wrapped = reference
	}
	return fmt.Sprintf("%s %s", wrapped, strings.ToLower(order.Direction))
}

type SqlBuilder struct {
	sqlParts []string
}

func (s *SqlBuilder) Append(str string) *SqlBuilder {
	s.sqlParts = append(s.sqlParts, str)
	return s
}

func (s SqlBuilder) ToSqlString() string {
	return strings.Join(s.sqlParts, "")
}

func NewSqlBuilder(base string) *SqlBuilder {
	sqlBuilder := &SqlBuilder{}
	if len(base) > 0 {
		return sqlBuilder.Append(base)
	}
	return sqlBuilder
}

func (s *SqlBuilder) PageOrderBy(pageable *util.Pageable, defaultSort *util.Sort) *SqlBuilder {
	sort := pageable.Sort
	if sort == nil || reflect.ValueOf(sort).IsNil() {
		if defaultSort == nil || reflect.ValueOf(defaultSort).IsNil() {
			panic(defaultSortError)
		}
		sort = defaultSort
	}
	s.Append(sort.SqlString())
	s.Append(fmt.Sprintf(" LIMIT %d, %d", pageable.Offset, pageable.PageSize))
	return s
}
