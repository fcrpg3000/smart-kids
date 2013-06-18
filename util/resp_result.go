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

const (
	SUCESS_CODE  = 1
	FAILURE_CODE = 0
	ERROR_CODE   = -1
)

type ResponseResult struct {
	Code    int
	Message string
	Values  map[string]interface{}
}

func (r *ResponseResult) IsSuccessful() bool {
	return r.Code == SUCESS_CODE
}

func (r *ResponseResult) IsFailed() bool {
	return r.Code == FAILURE_CODE
}

func (r *ResponseResult) IsError() bool {
	return r.Code == ERROR_CODE
}

func (r *ResponseResult) AddValue(k string, v interface{}) *ResponseResult {
	if len(k) == 0 || v == nil {
		return r
	}
	r.Values[k] = v
	return r
}

func (r *ResponseResult) AddValues(added map[string]interface{}) *ResponseResult {
	if len(added) == 0 {
		return r
	}
	for k, v := range added {
		if v != nil {
			r.Values[k] = v
		}
	}
	return r
}

func (r *ResponseResult) SetValues(newValues map[string]interface{}) *ResponseResult {
	r.Values = newValues
	return r
}

func SuccessResult(message string) *ResponseResult {
	return &ResponseResult{Code: SUCESS_CODE, Message: message}
}

func FailureResult(message string) *ResponseResult {
	return &ResponseResult{Code: FAILURE_CODE, Message: message}
}

func ErrorResult(message string) *ResponseResult {
	return &ResponseResult{Code: ERROR_CODE, Message: message}
}
