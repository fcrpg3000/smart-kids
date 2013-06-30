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
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	notSupportComparator = errors.New("The given values not support comparator.")
)

// Add params to the url.
// Examples:
// url1, url2 := "http://www.domain.com", "http://domain.com?top=1"
// newUrl1 := AddParamsToUrl(url1, map[string]string{
//     "param1": "value1",
//     "param2": "value2",
// })
// newUrl2 := AddParamsToUrl(url2, map[string]string{
//     "param1": "value1",
// })
// newUrl1 =&gt; "http://www.domain.com?param1=value1&param2=value2"
// newUrl2 =&gt; "http://domain.com?top=1&param1=value1"
func AddParamsToUrl(url string, params map[string]string) string {
	var (
		queryString []string
		sep         string
	)

	if strings.Contains(url, "?") {
		sep = "&"
	} else {
		sep = "?"
	}
	for k, v := range params {
		queryString = append(queryString, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join([]string{
		url, sep, strings.Join(queryString, "&"),
	}, "")
}

// >
func GreaterThan(a, b interface{}) bool {
	val, err := comparator(a, b)
	if err != nil {
		panic(err)
	}
	return val == 1
}

// >=
func GreaterThanOrEqual(a, b interface{}) bool {
	val, err := comparator(a, b)
	if err != nil {
		panic(err)
	}
	return val >= 0
}

// <
func LessThan(a, b interface{}) bool {
	val, err := comparator(a, b)
	if err != nil {
		panic(err)
	}
	return val == -1
}

// <=
func LessThanOrEqual(a, b interface{}) bool {
	val, err := comparator(a, b)
	if err != nil {
		panic(err)
	}
	return val <= 0
}

// Comparator two given values. Returns -1, 0 or 1 when a < b, a == b or a > b
// a and b are not the same type of error is returned
// Support type: int, int8, int16, int32, int64,
// uint, uint8, uint16, uint32, uint64,
// float32, float64, string
func comparator(a, b interface{}) (int, error) {
	if a == nil && b == nil {
		return 0, nil
	} else if a == nil && b != nil {
		return -1, nil
	} else if a != nil && b == nil {
		return 1, nil
	}
	a_type, b_type := reflect.TypeOf(a).Kind(), reflect.TypeOf(b).Kind()
	switch a_type {
	case reflect.String:
		if b_type != reflect.String {
			return -1, errors.New(fmt.Sprintf(
				"The given second's type is not string. %v", b_type))
		}
		a_str, b_str := reflect.ValueOf(a).String(), reflect.ValueOf(b).String()
		if a_str == b_str {
			return 0, nil
		} else if a_str > b_str {
			return 1, nil
		} else {
			return -1, nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var a_int, b_int int64
		a_int = reflect.ValueOf(a).Int()
		switch b_type {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			b_int = reflect.ValueOf(b).Int()
			break
		default:
			return -1, errors.New(fmt.Sprintf(
				"The given second's type is not int(8-64). %v", b_type))
		}
		if a_int == b_int {
			return 0, nil
		} else if a_int > b_int {
			return 1, nil
		} else {
			return -1, nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var a_uint, b_uint uint64
		a_uint = reflect.ValueOf(a).Uint()
		switch b_type {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			b_uint = reflect.ValueOf(b).Uint()
			break
		default:
			return -1, errors.New(fmt.Sprintf(
				"The given second's type is not uint(8-64). %v", b_type))
		}
		if a_uint == b_uint {
			return 0, nil
		} else if a_uint > b_uint {
			return 1, nil
		} else {
			return -1, nil
		}
	case reflect.Float32, reflect.Float64:
		var a_f, b_f float64
		a_f = reflect.ValueOf(a).Float()
		switch b_type {
		case reflect.Float32, reflect.Float64:
			b_f = reflect.ValueOf(b).Float()
			break
		default:
			return -1, errors.New(fmt.Sprintf(
				"The given second's type is not float(32|64). %v", b_type))
		}
		if a_f == b_f {
			return 0, nil
		} else if a_f > b_f {
			return 1, nil
		} else {
			return -1, nil
		}
	}
	return -1, notSupportComparator
}
