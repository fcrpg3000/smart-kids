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
)

var (
	Male         = &Gender{uint16(1), "男", "Male", "M"}
	Female       = &Gender{uint16(2), "女", "Female", "F"}
	SecretGender = &Gender{uint16(3), "保密", "Unknown", "N"}
	genderMap    = map[uint16]*Gender{
		Male.Code:         Male,
		Female.Code:       Female,
		SecretGender.Code: SecretGender,
	}
)

func GenderOf(code uint16) *Gender {
	gender, exists := genderMap[code]
	if exists {
		return gender
	}
	return SecretGender
}

type Gender struct {
	Code  uint16 `json:"code"`
	Name  string `json:"name"`
	EName string `json:"ename"`
	Alias string `json:"alias"`
}

func (g Gender) String() string {
	return fmt.Sprintf("Gender(%d,%s,%s)", g.Code, g.Name, g.Alias)
}

type Client struct {
	Code     uint16 `json:"code"`
	Name     string `json:"name"`
	IsMobile bool   `json:"isMobile"`
}

// 内置的客户端信息
var (
	UNKNOWN_CLIENT = &Client{uint16(0), "未知的客户端", false}
	WEB            = &Client{uint16(1), "普通版", false}
	WEB_PE         = &Client{uint16(2), "专业版", false}
	ANDROID        = &Client{uint16(3), "Android客户端", true}
	IPHONE         = &Client{uint16(4), "iPhone客户端", true}
	WINDOWS_PHONE  = &Client{uint16(5), "Windows.Phone客户端", true}
	IPAD           = &Client{uint16(6), "iPad客户端", true}
	FIREFOX        = &Client{uint16(7), "FirefoxOS客户端", true}
	clients        = map[uint16]*Client{
		WEB.Code:           WEB,
		WEB_PE.Code:        WEB_PE,
		ANDROID.Code:       ANDROID,
		IPHONE.Code:        IPHONE,
		WINDOWS_PHONE.Code: WINDOWS_PHONE,
		IPAD.Code:          IPAD,
		FIREFOX.Code:       FIREFOX,
	}
)

func ClientOf(code uint16, def *Client) *Client {
	if client, exists := clients[code]; exists {
		return client
	}
	return def
}

type Visibility struct {
	Code uint16 `json:"code"`
	Name string `json:"name"`
}

var (
	V_ALL        = &Visibility{uint16(1), "所有人可见"}
	V_FRIENDS    = &Visibility{uint16(2), "仅好友可见"}
	V_PASSWORD   = &Visibility{uint16(3), "凭密码访问"}
	V_SELF       = &Visibility{uint16(4), "仅自己可见"}
	Visibilities = map[uint16]*Visibility{
		V_ALL.Code:      V_ALL,
		V_FRIENDS.Code:  V_FRIENDS,
		V_PASSWORD.Code: V_PASSWORD,
		V_SELF.Code:     V_SELF,
	}
)

// Returns code mapping Visibility, def if code not found.
func VisibilityOf(code uint16, def *Visibility) *Visibility {
	if v, ok := Visibilities[code]; ok {
		return v
	}
	return def
}

// Implement Cacheable this interface can be cached.
type Cacheable interface {

	// Returns implement this interface struct object key.
	CacheKey() string
	// Returns implement this interface struct object id list key.
	IdListKey() string
}
