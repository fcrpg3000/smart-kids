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
package util

import (
	_ "fmt"
	"reflect"
	"testing"
)

func assertTrue(t *testing.T, expr bool, arg1 string, arg2 string) {
	if expr {
		t.Log(arg1)
	} else {
		t.Error(arg2)
	}
}

func assertFalse(t *testing.T, expr bool, arg1 string, arg2 string) {
	if !expr {
		t.Log(arg1)
	} else {
		t.Error(arg2)
	}
}

func assertNull(t *testing.T, obj interface{}, arg1 string, arg2 string) {
	if obj == nil || reflect.ValueOf(obj).IsNil() {
		t.Log(arg1)
	} else {
		t.Error(arg2)
	}
}

type testBean struct {
	Id   int
	Name string
}

func TestEmptyPage(t *testing.T) {
	page := NewPage(nil, nil, 0)
	assertTrue(t, page.Current == 1, "PASS", "Empty page default current page is 1")
	assertTrue(t, page.PageSize == 0, "PASS", "Empty page size is 0")
	assertTrue(t, page.TotalPages == 1, "PASS", "Empty page defautl total pages is 1")
	assertTrue(t, page.NumberOfElements == 0, "PASS", "Empty page content length is 0")
	assertTrue(t, page.TotalElements == 0, "PASS", "Empty page default total elements is 0")
	assertFalse(t, page.HasPrevPage, "PASS", "Empty page no previous page.")
	assertTrue(t, page.IsFirstPage, "PASS", "Empty page current page is first.")
	assertFalse(t, page.HasNextPage, "PASS", "Empty page no next page.")
	assertNull(t, page.Sort, "PASS", "Empty page no sort.")
	assertNull(t, page.NextPageable(), "PASS", "Empty page no next pageable.")
	assertNull(t, page.PrevPageable(), "PASS", "Empty page no prev pageable.")
}

func TestBeanPage(t *testing.T) {
	testBeans := []testBean{
		testBean{2, "John"},
		testBean{4, "Shammy"},
		testBean{1, "Zone"},
		testBean{3, "Amy"},
	}
	sort := NewSort(DESC, []string{"Id"})
	pageable, err := NewPageable0(1, 2, sort)
	if err != nil {
		t.Error(err.Error())
	}
	var content = make([]interface{}, 0)
	content = append(content, testBeans[1], testBeans[3])
	page := NewPage(content, pageable, int64(len(testBeans)))
	assertTrue(t, page.TotalPages == 2, "PASS", "Page calculate error?")
}
