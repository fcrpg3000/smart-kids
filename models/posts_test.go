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
	"testing"
	"time"
)

func TestDescription(t *testing.T) {
	createdTime, _ := time.Parse("2006-01-02", "2013-06-15")
	lastModifiedTime := createdTime
	p := Posts{int64(1), int64(-1), "这里是内容！", createdTime, lastModifiedTime,
		int64(1), "king4go", WEB.Id, nil, nil}
	p.PostGet()
	description := fmt.Sprintf("%s: %s %s From %s",
		p.UserName, p.Content, p.CreatedTime.Format(time.RFC3339),
		p.Client.Name)
	if p.Description() != description {
		t.Error("#Description() format error. actual: ", description)
	} else {
		t.Log("#Description() format is ok. description: ", description)
	}
}

func TestNewForwardError(t *testing.T) {
	// Returns error if src posts id is zero or negative
	posts, err := NewForward(int64(0), &User{}, WEB, "New content string")
	if err != nil && posts == nil {
		t.Log("The source Posts id is 0, Returns a error.", err.Error())
	} else {
		t.Error("The source Posts id is 0, But not error. ")
	}
	// Returns error if user is nil
	posts, err = NewForward(int64(1), nil, WEB_PE, "New content string")
	if err != nil && posts == nil {
		t.Log("The user Ptr is nil, Returns a error.", err.Error())
	} else {
		t.Error("The user Ptr is nil, But not error.")
	}
	// Returns error if client is nil
	posts, err = NewForward(int64(1), &User{}, nil, "New content string")
	if err != nil && posts == nil {
		t.Log("The client Ptr is nil, Returns a error.", err.Error())
	} else {
		t.Error("The client Ptr is nil, But not error.")
	}

	// Returns error if content length is 0
	posts, err = NewForward(int64(1), &User{}, IPAD, "")
	if err != nil && posts == nil {
		t.Log("The content length is 0, Returns a error.", err.Error())
	} else {
		t.Error("The content length is 0, But not error.")
	}
}

func TestNewForward(t *testing.T) {
	posts, err := NewForward(int64(1), &User{}, ANDROID, "I'm here anytime.")
	if err != nil && posts == nil {
		t.Error("NewForward is error ", err.Error())
	} else {
		t.Log("NewForward returns a new Posts Ptr.")
	}
}
