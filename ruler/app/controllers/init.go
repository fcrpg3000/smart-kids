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
package controllers

import (
	"github.com/robfig/revel"
	"reflect"
	"strings"
)

func init() {
	revel.OnAppStart(Init)
	revel.InterceptMethod((*GorpController).Begin, revel.BEFORE)
	revel.InterceptMethod(Application.AddAdmin, revel.BEFORE)
	revel.InterceptMethod(Application.AddMenus, revel.BEFORE)
	// revel.InterceptMethod(Hotels.checkUser, revel.BEFORE)
	revel.InterceptMethod((*GorpController).Commit, revel.AFTER)
	revel.InterceptMethod((*GorpController).Rollback, revel.FINALLY)

	revel.TemplateFuncs["gt"] = greaterThan
	revel.TemplateFuncs["replaceAll"] = replaceAll
}

func greaterThan(a, b interface{}) bool {
	switch a.(type) {
	case string:
		switch b.(type) {
		case string:
			return reflect.ValueOf(a).String() > reflect.ValueOf(b).String()
		}
	case int, int8, int16, int32, int64:
		switch b.(type) {
		case int, int8, int16, int32, int64:
			return reflect.ValueOf(a).Int() > reflect.ValueOf(b).Int()
		}
	case uint, uint8, uint16, uint32, uint64:
		switch b.(type) {
		case uint, uint8, uint16, uint32, uint64:
			return reflect.ValueOf(a).Uint() > reflect.ValueOf(b).Uint()
		}
	case float32, float64:
		switch b.(type) {
		case float32, float64:
			return reflect.ValueOf(a).Float() > reflect.ValueOf(b).Float()
		}
	}
	return false
}

func replaceAll(src, old, newVal interface{}) string {
	var newStr string
	s := reflect.ValueOf(src).String()
	o := reflect.ValueOf(old).String()
	switch newVal.(type) {
	case int, int8, int16, int32, int64:
		newStr = string(reflect.ValueOf(newVal).Int())
		break
	default:
		newStr = reflect.ValueOf(newVal).String()
	}
	return strings.Replace(s, o, newStr, 100)
}
