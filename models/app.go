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
	"crypto/md5"
	_ "database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

// app module table names
const (
	DEVELOPER_TABLE   = "sk_developer"
	APP_TABLE         = "sk_app"
	APP_SESSION_TABLE = "sk_app_session"
)

// sk_developer fields constants
const (
	F_DEV_TYPE    = "dev_type"
	F_DEV_NAME    = "dev_name"
	F_PROVINCE_ID = "province_id"
	F_CITY_ID     = "city_id"
	F_PHONE       = "phone"
	F_DEV_IM_TYPE = "dev_im_type"
	F_DEV_IM      = "dev_im"
	F_DEV_SITE    = "dev_site"
	F_IS_TRUSTED  = "is_trusted"
)

// Developer's typeId and ImType enumeration
const (
	DT_PERSONAL  = int16(1) // 个人
	DT_COMPANY   = int16(2) // 公司
	DIT_GTALK    = int16(1) // GTalk
	DIT_QQ       = int16(2) // 腾讯QQ
	DIT_WANGWANG = int16(3) // 淘宝旺旺
)

// public var fields of Developer.
var (
	DeveloperFields = strings.Join([]string{
		F_USER_ID, F_USER_NAME, F_DEV_TYPE, F_DEV_NAME, F_PROVINCE_ID,
		F_CITY_ID, F_EMAIL, F_PHONE, F_DEV_IM_TYPE, F_DEV_IM, F_DEV_SITE,
		F_IS_TRUSTED, F_CREATED_TIME, F_LAST_MODIFIED_TIME,
	}, ", ")
)

// Smart Kids developer
type Developer struct {
	UserId           int64          `db:"user_id" json:"uid"`
	UserName         string         `db:"user_name" json:"userName"`
	DevType          int16          `db:"dev_type"`
	DevName          string         `db:"dev_name"`
	ProvinceId       int            `db:"province_id"`
	CityId           int            `db:"city_id"`
	Email            string         `db:"email"`
	Phone            string         `db:"phone"`
	DevImType        int16          `db:"dev_im_type"`
	DevIm            string         `db:"dev_im"`
	DevSite          string         `db:"dev_site"`
	IsTrusted        bool           `db:"is_trusted"`
	CreatedTime      mysql.NullTime `db:"created_time"`
	LastModifiedTime mysql.NullTime `db:"last_modified_time"`

	// Transient fields
	User     *User     `db:"-"`
	Location *Location `db:"-"`
}

func ToDeveloper(i interface{}, err error) *Developer {
	if err != nil {
		panic(err)
	}
	if i == nil || reflect.ValueOf(i).IsNil() {
		return nil
	}
	return i.(*Developer)
}

func ToDevelopers(results []interface{}, err error) []*Developer {
	if err != nil {
		panic(err)
	}
	size := len(results)
	developers := make([]*Developer, size)
	if size == 0 {
		return developers
	}
	for i, result := range results {
		developers[i] = result.(*Developer)
	}
	return developers
}

// sk_app fields constants
const (
	F_APP_ID         = "app_id"
	F_APP_NAME       = "app_name"
	F_APP_CATE_ID    = "cate_id"
	F_BASE_APP_OS    = "base_app_os"
	F_APP_OS         = "app_os"
	F_APP_URL        = "app_url"
	F_APP_SUMMARY    = "summary"
	F_APP_DESC       = "description"
	F_IS_BIND_DOMAIN = "is_bind_domain"
	F_TAG_ID1        = "tag_id1"
	F_TAG_ID2        = "tag_id2"
	F_TAG_ID3        = "tag_id3"
	F_APP_KEY        = "app_key"
	F_APP_SECRET     = "app_secret"
)

var (
	AppFields = strings.Join([]string{
		F_APP_ID, F_APP_NAME, F_APP_CATE_ID, F_BASE_APP_OS, F_APP_OS,
		F_APP_URL, F_APP_SUMMARY, F_APP_DESC, F_IS_BIND_DOMAIN, F_TAG_ID1,
		F_TAG_ID2, F_TAG_ID3, F_USER_ID, F_USER_NAME, F_APP_KEY,
		F_APP_SECRET, F_CREATED_TIME, F_LAST_MODIFIED_TIME,
	}, ", ")
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type App struct {
	Id               int64          `db:"app_id"`
	Name             string         `db:"app_name"` // Unique Index
	CateId           int            `db:"cate_id"`
	BaseOsId         int            `db:"base_app_os"`
	OsId             int            `db:"app_os"`
	Url              string         `db:"app_url"` // Unique Index
	Summary          string         `db:"summary"`
	Description      string         `db:"description"`
	IsBindDomain     bool           `db:"is_bind_domain"`
	TagId1           int            `db:"tag_id1"`
	TagId2           int            `db:"tag_id2"`
	TagId3           int            `db:"tag_id3"`
	UserId           int64          `db:"user_id"`
	UserName         string         `db:"user_name"`
	AppKey           string         `db:"app_key"`    // Unique Index
	AppSecret        string         `db:"app_secret"` // Unique Index
	CreatedTime      mysql.NullTime `db:"created_time"`
	LastModifiedTime mysql.NullTime `db:"last_modified_time"`
}

func ToApp(i []interface{}, err error) *App {
	if len(i) == 0 || i[0] == nil || reflect.ValueOf(i[0]).IsNil() {
		return nil
	}
	return i[0].(*App)
}

func ToApps(results []interface{}, err error) []*App {
	if err != nil {
		panic(err)
	}
	size := len(results)
	apps := make([]*App, size)
	if size == 0 {
		return apps
	}
	for i, result := range results {
		apps[i] = result.(*App)
	}
	return apps
}

// sk_app fields constants
const (
	F_APP_ACCESS_TOKEN = "access_token"
	F_APP_AUTH_CODE    = "app_auth_code"
	F_LAST_ACCESS_TIME = "last_access_time"
)

var (
	AppSessionFields = strings.Join([]string{
		F_APP_ID, F_APP_NAME, F_APP_AUTH_CODE, F_APP_ACCESS_TOKEN, F_APP_KEY,
		F_APP_SECRET, F_LAST_ACCESS_TIME, F_CREATED_TIME, F_LAST_MODIFIED_TIME,
	}, ", ")
)

type AppSession struct {
	AppId            int64          `db:"app_id"`
	AppName          string         `db:"app_name"`
	AppAuthCode      string         `db:"app_auth_code"` // Dynamic change
	AccessToken      string         `db:"access_token"`  // Unique Index
	AppKey           string         `db:"app_key"`
	AppSecret        string         `db:"app_secret"`
	LastAccessTime   int64          `db:"last_access_time"`
	CreatedTime      mysql.NullTime `db:"created_time"`
	LastModifiedTime mysql.NullTime `db:"last_modified_time"`
}

// Flushes AppSession's AppAuthCode value, at the same time,
// flushes AccessToken value. Returns new AppAuthCode value.
func (a *AppSession) FlushAuthCode() string {
	authCode := string(random.Int63n(8))
	h := md5.New()
	h.Write([]byte(fmt.Sprintf("AuthCode{%s,%s,%s}", a.AppKey,
		a.AppSecret, authCode)))
	accessToken := fmt.Sprintf("%x", h.Sum(nil))
	a.AppAuthCode = authCode
	a.AccessToken = accessToken
	return a.AppAuthCode
}

func NewAppSession(app *App) *AppSession {
	timeNow := time.Now()
	session := &AppSession{
		AppId: app.Id, AppName: app.Name, AppKey: app.AppKey,
		AppSecret: app.AppSecret, LastAccessTime: timeNow.Unix(),
		CreatedTime:      mysql.NullTime{timeNow, true},
		LastModifiedTime: mysql.NullTime{timeNow, true},
	}
	session.FlushAuthCode()
	return session
}

func ToAppSession(i interface{}, err error) *AppSession {
	if err != nil {
		panic(err)
	}
	if i == nil || reflect.ValueOf(i).IsNil() {
		return nil
	}
	return i.(*AppSession)
}

func ToAppSessions(results []interface{}, err error) []*AppSession {
	if err != nil {
		panic(err)
	}
	size := len(results)
	appSessions := make([]*AppSession, size)
	if size == 0 {
		return appSessions
	}

	for i, result := range results {
		appSessions[i] = result.(*AppSession)
	}
	return appSessions
}

type AuthError struct {
	Error   string `json:"error"`
	Code    int    `json:"error_code"`
	Message string `json:"error_description"`
}

func (a *AuthError) Equals(other *AuthError) bool {
	if other == nil {
		return false
	}
	return a.Error == other.Error && a.Code == other.Code
}

// All AuthError enumeration instances.
var (
	Err_Redirect_URI_Mismatch     = &AuthError{"redirect_uri_mismatch", 21322, "重定向地址不匹配"}
	Err_Invalid_Request           = &AuthError{"invalid_request", 21323, "请求不合法"}
	Err_Invalid_Client            = &AuthError{"invalid_client", 21324, "client_id或client_secret参数无效"}
	Err_Invalid_Grant             = &AuthError{"invalid_grant", 21325, "提供的Access Grant是无效的、过期的或已撤销的"}
	Err_Unauthorized_Client       = &AuthError{"unauthorized_client", 21326, "客户端没有权限"}
	Err_Expired_Token             = &AuthError{"expired_token", 21327, "token过期"}
	Err_unsupported_grant_type    = &AuthError{"unsupported_grant_type", 21328, "不支持的 GrantType"}
	Err_unsupported_response_type = &AuthError{"unsupported_response_type", 21329, "不支持的 ResponseType"}
	Err_access_denied             = &AuthError{"access_denied", 21330, "用户或授权服务器拒绝授予数据访问权限"}
	Err_temporarily_unavailable   = &AuthError{"temporarily_unavailable", 21331, "服务暂时无法访问"}
)
