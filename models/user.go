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
	"database/sql"
	_ "encoding/hex"
	"errors"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/go-sql-driver/mysql"
	"github.com/robfig/revel"
	"reflect"
	"regexp"
	"strings"
	"time"
)

const (
	USER_TABLE        = "sk_user"
	USER_AVATAR_TABLE = "sk_user_avatar"
	USER_INFO_TABLE   = "sk_user_info"
	BANNED_USER_TABLE = "sk_banned_user"
	// commons table fields
	F_ID                 = "id"
	F_USER_ID            = "user_id"
	F_USER_NAME          = "user_name"
	F_CREATED_TIME       = "created_time"
	F_LAST_MODIFIED_TIME = "last_modified_time"
)

var (
	Male         = Gender{1, "男", "M"}
	Female       = Gender{2, "女", "F"}
	SecretGender = Gender{3, "保密", "N"}
	genderMap    = map[int]Gender{
		1: Male,
		2: Female,
		3: SecretGender}

	emptyNameAndPwd = errors.New("The userName and password must not be empty.")
)

func GenderOf(code int) Gender {
	gender, exists := genderMap[code]
	if exists {
		return gender
	}
	return SecretGender
}

type Gender struct {
	Code  int    `json:"code"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

func (g Gender) String() string {
	return fmt.Sprintf("Gender(%d,%s,%s)", g.Code, g.Name, g.Alias)
}

// User table fields
const (
	F_HASH_PASSWORD  = "hash_password"
	F_PASSWORD_SALT  = "password_salt"
	F_EMAIL          = "email"
	F_COMMONLY_EMAIL = "commonly_email"
)

var (
	UserFields = strings.Join([]string{
		F_USER_ID, F_USER_NAME, F_HASH_PASSWORD, F_PASSWORD_SALT,
		F_EMAIL, F_COMMONLY_EMAIL, F_CREATED_TIME, F_LAST_MODIFIED_TIME,
	}, ", ")
)

// User model for mapping table `sk_user`
type User struct {
	UserId           int64          `db:"user_id" json:"uid"`
	UserName         string         `db:"user_name" json:"userName"`
	HashPassword     string         `db:"hash_password" json:"-"`
	PasswordSalt     string         `db:"password_salt" json:"-"`
	Email            sql.NullString `db:"email" json:"-"`
	CommonlyEmail    sql.NullString `db:"commonly_email" json:"-"`
	CreatedTime      mysql.NullTime `db:"created_time" json:"register"`
	LastModifiedTime mysql.NullTime `db:"last_modified_time" json:"-"`

	// Transient property
	Password    string `db:"-" json:"-"`
	EmailValue  string `db:"-" json:"email,omitempty"`
	Email2Value string `db:"-" json:"email2,omitempty"`
}

func (u User) String() string {
	return fmt.Sprintf("User(%d, %s, %s, %s)", u.UserId, u.UserName, u.HashPassword, u.Email)
}

func NewUser(userName string, password string, others map[string]interface{}) (*User, error) {
	user := &User{}
	timeNow := time.Now()
	if len(userName) == 0 || len(password) == 0 {
		return nil, emptyNameAndPwd
	}
	user.UserName = userName
	user.Password = password
	user.CreatedTime = mysql.NullTime{timeNow, true}
	user.LastModifiedTime = mysql.NullTime{timeNow, true}

	if len(others) == 0 {
		return user, nil
	}
	for key, val := range others {
		switch key {
		case "Email":
			if len(val.(string)) != 0 {
				user.Email = sql.NullString{val.(string), true}
			} else {
				user.Email = sql.NullString{"", false}
			}
			break
		case "CommonlyEmail":
			if len(val.(string)) != 0 {
				user.CommonlyEmail = sql.NullString{val.(string), true}
			} else {
				user.CommonlyEmail = sql.NullString{"", false}
			}
			break
		}
	}
	return user, nil
}

var userNameRegexp = regexp.MustCompile("^\\w+$")

func (user *User) Validate(v *revel.Validation) {

	v.Check(user.UserName,
		revel.Required{},
		revel.MinSize{8},
		revel.MaxSize{50},
		revel.Match{userNameRegexp},
	)

	ValidatePassword(v, user.Password).Key("user.Password")
}

func ValidatePassword(v *revel.Validation, password string) *revel.ValidationResult {
	return v.Check(password,
		revel.Required{},
		revel.MinSize{8},
		revel.MaxSize{16},
	)
}

func (u *User) PreInsert(_ gorp.SqlExecutor) error {
	return nil
}

func (u *User) PostGet() error {
	if u.Email.Valid {
		u.EmailValue = u.Email.String
	}
	if u.CommonlyEmail.Valid {
		u.Email2Value = u.CommonlyEmail.String
	}
	return nil
}

func ToUser(i interface{}, err error) *User {
	if err != nil {
		panic(err)
	}
	if i == nil || reflect.ValueOf(i).IsNil() {
		return nil
	}
	return i.(*User)
}

func ToUsers(results []interface{}, err error) []*User {
	if err != nil {
		panic(err)
	}
	size := len(results)
	users := make([]*User, size)
	if size == 0 {
		return users
	}
	for i, result := range results {
		users[i] = result.(*User)
	}
	return users
}

const (
	F_IMAGE_DOMAIN      = "image_domain"
	F_AVATAR_PATH       = "avatar_path"
	F_SRC_AVATAR_PATH   = "src_avatar_path"
	F_SMALL_AVATAR_PATH = "small_avatar_path"
	F_THUMB_AVATAR_PATH = "thumb_avatar_path"
	F_AVATAR_NAME       = "avatar_name"
)

var (
	UserAvatarFields = strings.Join([]string{
		F_USER_ID, F_USER_NAME, F_IMAGE_DOMAIN, F_AVATAR_PATH,
		F_SRC_AVATAR_PATH, F_SMALL_AVATAR_PATH, F_THUMB_AVATAR_PATH,
		F_AVATAR_NAME, F_CREATED_TIME, F_LAST_MODIFIED_TIME,
	}, ", ")
)

// UserAvatar struct
// ----------------------------------------------------------------------------

type UserAvatar struct {
	UserId           string         `db:"user_id"`
	UserName         string         `db:"user_name"`
	ImageDomain      string         `db:"image_domain"`
	AvatarPath       sql.NullString `db:"avatar_path"`       // 150x150 maybe
	SrcAvatarPath    sql.NullString `db:"src_avatar_path"`   // source size
	SmallAvatarPath  sql.NullString `db:"small_avatar_path"` // 80x80 maybe
	ThumbAvatarPath  sql.NullString `db:"thumb_avatar_path"` // 40x40 maybe
	AvatarName       sql.NullString `db:"avatar_name"`
	CreatedTime      mysql.NullTime `db:"created_time"`
	LastModifiedTime mysql.NullTime `db:"last_modified_time"`
}

// Returns user's source avatar image url.
func (u UserAvatar) SrcAvatarUrl() string {
	if u.SrcAvatarPath.Valid {
		return u.avatarUrlInternal(u.SrcAvatarPath.String)
	}
	return ""
}

// Returns user's normal avatar image url.
func (u UserAvatar) AvatarUrl() string {
	if u.AvatarPath.Valid {
		return u.avatarUrlInternal(u.AvatarPath.String)
	}
	return ""
}

// Returns user's small avatar image url.
func (u UserAvatar) SmallAvatarUrl() string {
	if u.SmallAvatarPath.Valid {
		return u.avatarUrlInternal(u.SmallAvatarPath.String)
	}
	return ""
}

// Returns user's thumb avatar image url.
func (u UserAvatar) ThumbAvatarUrl() string {
	if u.ThumbAvatarPath.Valid {
		return u.avatarUrlInternal(u.ThumbAvatarPath.String)
	}
	return ""
}

// Returns user's avatar image url of the specified size path.
func (u UserAvatar) avatarUrlInternal(path string) string {
	return fmt.Sprintf("http://%s%s%s", u.ImageDomain, path, u.AvatarName)
}

func ToUserAvatar(i interface{}, err error) *UserAvatar {
	if err != nil {
		panic(err)
	}
	if i == nil || reflect.ValueOf(i).IsNil() {
		return nil
	}
	return i.(*UserAvatar)
}

func ToUserAvatars(results []interface{}, err error) []*UserAvatar {
	if err != nil {
		panic(err)
	}
	size := len(results)
	userAvatars := make([]*UserAvatar, size)
	if size == 0 {
		return userAvatars
	}
	for i, r := range results {
		userAvatars[i] = r.(*UserAvatar)
	}
	return userAvatars
}

const (
	F_OPERATOR_ID   = "operator_id"
	F_OPERATOR_NAME = "operator_name"
	F_CAUSE         = "banned_cause"
	F_IS_PERMANENT  = "is_permanent"
	F_BANNED_TIME   = "banned_time"
	F_UNBAN_TIME    = "unban_time"
)

var (
	BannedUserFields = strings.Join([]string{
		F_ID, F_USER_ID, F_USER_NAME, F_OPERATOR_ID, F_OPERATOR_NAME,
		F_CAUSE, F_IS_PERMANENT, F_BANNED_TIME, F_UNBAN_TIME,
		F_CREATED_TIME, F_LAST_MODIFIED_TIME,
	}, ", ")
)

// BannedUser struct
// ----------------------------------------------------------------------------

type BannedUser struct {
	Id               int64          `db:"id"`
	UserId           int64          `db:"user_id"`
	UserName         string         `db:"user_name"`
	OperatorId       int            `db:"operator_id"`
	OperatorName     string         `db:"operator_name"`
	Cause            string         `db:"banned_cause"`
	IsPermanent      bool           `db:"is_permanent"`
	BannedTime       mysql.NullTime `db:"banned_time"`
	UnbanTime        mysql.NullTime `db:"unban_time"`
	CreatedTime      mysql.NullTime `db:"created_time"`
	LastModifiedTime mysql.NullTime `db:"last_modified_time"`
}

// BannedUser instance default string
func (b BannedUser) String() string {
	return fmt.Sprintf("BannedUser{Id=%d,Target=(%d, %s),Operator(%d, %s),"+
		"Cause=\"%s\",Permanent=%v,Period=(%v - %v),"+
		"LastModified=%v}",
		b.Id, b.UserId, b.UserName, b.OperatorId, b.OperatorName,
		b.Cause, b.IsPermanent, b.BannedTime.Time, b.UnbanTime.Time,
		b.LastModifiedTime.Time)
}

func ToBannedUser(i interface{}, err error) *BannedUser {
	if err != nil {
		panic(err)
	}
	if i == nil || reflect.ValueOf(i).IsNil() {
		return nil
	}
	return i.(*BannedUser)
}

func ToBannedUsers(results []interface{}, err error) []*BannedUser {
	if err != nil {
		panic(err)
	}
	size := len(results)
	bannedUsers := make([]*BannedUser, size)
	if size == 0 {
		return bannedUsers
	}
	for i, r := range results {
		bannedUsers[i] = r.(*BannedUser)
	}
	return bannedUsers
}

// sk_user_info fields constant
const (
	F_NICKNAME         = "nickname"
	F_GENDER_CODE      = "gender_code"
	F_CALENDAR_MODE    = "calendar_mode"
	F_DATE_OF_BIRTH    = "date_of_birth"
	F_HT_COUNTRY_ID    = "ht_country_id"
	F_HT_STATE_ID      = "ht_state_id"
	F_HT_CITY_ID       = "ht_city_id"
	F_HT_DIST_ID       = "ht_dist_id"
	F_POR_COUNTRY_ID   = "por_country_id"
	F_POR_STATE_ID     = "por_state_id"
	F_POR_CITY_ID      = "por_city_id"
	F_POR_DIST_ID      = "por_dist_id"
	F_OTHER_STATE      = "other_state"
	F_EDU_ID           = "edu_id"
	F_FEELING_ID       = "feeling_id"
	F_BLOOD_TYPE_ID    = "blood_type_id"
	F_CONSTELLATION_ID = "constellation_id"
)

var (
	UserInfoFields = strings.Join([]string{
		F_USER_ID, F_USER_NAME, F_NICKNAME, F_GENDER_CODE, F_CALENDAR_MODE,
		F_DATE_OF_BIRTH, F_HT_COUNTRY_ID, F_HT_STATE_ID, F_HT_CITY_ID,
		F_HT_DIST_ID, F_POR_COUNTRY_ID, F_POR_STATE_ID, F_POR_CITY_ID, F_POR_DIST_ID,
		F_OTHER_STATE, F_EDU_ID, F_FEELING_ID, F_BLOOD_TYPE_ID, F_CONSTELLATION_ID,
	}, ", ")
)

// UserInfo struct
// ----------------------------------------------------------------------------

type UserInfo struct {
	UserId           int64          `db:"user_id" json:"uid"`        // not autoincrement
	UserName         string         `db:"user_name" json:"userName"` // just redundancy field
	Nickname         sql.NullString `db:"nickname" json:"nickname,omitempty"`
	GenderCode       int            `db:"gender_code"`
	CalendarMode     int16          `db:"calendar_mode"`
	DateOfBirthStr   sql.NullString `db:"date_of_birth" json:"-"`
	HtCountryId      int            `db:"ht_country_id" json:"-"`
	HtStateId        int            `db:"ht_province_id" json:"-"`
	HtCityId         int            `db:"ht_city_id" json:"-"`
	HtDistId         int            `db:"ht_dist_id" json:"-"`
	PorCountryId     int            `db:"por_country_id" json:"-"`
	PorStateId       int            `db:"por_province_id" json:"-"`
	PorCityId        int            `db:"por_city_id" json:"-"`
	PorDistId        int            `db:"por_dist_id" json:"-"`
	OtherState       sql.NullString `db:"other_state"`
	EduId            int            `db:"edu_id" json:"-"`
	FeelingId        int            `db:"feeling_id" json:"-"`
	BloodTypeId      int            `db:"blood_type_id" json:"-"`
	ConstellationId  int            `db:"constellation_id" json:"-"`
	CreatedTime      mysql.NullTime `db:"created_time"`
	LastModifiedTime mysql.NullTime `db:"last_modified_time"`

	// Transient
	Gender           Gender         `db:"-" json:"gender,omitempty"`
	DateOfBirth      time.Time      `db:"-" json:"dateOfBirth,omitempty"`
	User             *User          `db:"-" json:"user,omitempty"`
	Hometown         *Location      `db:"-" json:"hometown,omitempty"`
	PlaceOfResidence *Location      `db:"-" json:"placeOfResidence,omitempty"`
	Education        *Education     `db:"-" json:"education,omitempty"`
	Feeling          *Feeling       `db:"-" json:"feeling,omitempty"`
	BloodType        *BloodType     `db:"-" json:"bloodType,omitempty"`
	Constellation    *Constellation `db:"-" json:"constellation,omitempty"`
}

// UserInfo struct builder
// Examples:
// builder := &UserInfoBuilder{&UserInfo{}}
// userInfo := builder.User(user).Nickname("MyName")
// .Education(EducationOf(2))...Feeling(FeelingOf(2)).Builder()
type UserInfoBuilder struct {
	userInfo *UserInfo
}

// Set a pointer to user for this builder
func (u *UserInfoBuilder) User(user *User) *UserInfoBuilder {
	u.userInfo.User = user
	return u
}

// Set a nickname for this builder
func (u *UserInfoBuilder) Nickname(nickname string) *UserInfoBuilder {
	if len(nickname) == 0 {
		u.userInfo.Nickname = sql.NullString{"", false}
	} else {
		u.userInfo.Nickname = sql.NullString{nickname, true}
	}
	return u
}

// Set a pointer to Location for this builder
func (u *UserInfoBuilder) Hometown(hometown *Location) *UserInfoBuilder {
	u.userInfo.Hometown = hometown
	return u
}

// Set a pointer to Location for this builder
func (u *UserInfoBuilder) PlaceOfResidence(por *Location) *UserInfoBuilder {
	u.userInfo.PlaceOfResidence = por
	return u
}

// Set a pointer to Education for this builder
func (u *UserInfoBuilder) Education(edu *Education) *UserInfoBuilder {
	u.userInfo.Education = edu
	return u
}

// Set a pointer to Feeling for this builder
func (u *UserInfoBuilder) Feeling(feeling *Feeling) *UserInfoBuilder {
	u.userInfo.Feeling = feeling
	return u
}

// Set a pointer to BloodType for this builder
func (u *UserInfoBuilder) BloodType(bloodType *BloodType) *UserInfoBuilder {
	u.userInfo.BloodType = bloodType
	return u
}

// Set a pointer to Constellation for this builder
func (u *UserInfoBuilder) Constellation(constellation *Constellation) *UserInfoBuilder {
	u.userInfo.Constellation = constellation
	return u
}

// Set a Gender for this builder
func (u *UserInfoBuilder) Gender(gender Gender) *UserInfoBuilder {
	u.userInfo.Gender = gender
	return u
}

// Set date of birth and calendar mode for this builder
func (u *UserInfoBuilder) DateOfBirth(dateOfBirth time.Time, mode int16) *UserInfoBuilder {
	u.userInfo.DateOfBirth = dateOfBirth
	u.userInfo.CalendarMode = mode
	return u
}

// Returns a pointer pointing to UserInfo from this builder
func (u *UserInfoBuilder) Builder() *UserInfo {
	return u.userInfo
}

// Gorp's lack of support for loading relations automatically.
func (u *UserInfo) PreInsert(_ gorp.SqlExecutor) error {
	if u.User != nil {
		u.UserId = u.User.UserId
		u.UserName = u.User.UserName
	}
	if u.Gender.Code > 0 {
		u.GenderCode = u.Gender.Code
	}
	if u.Hometown != nil {
		if u.Hometown.Country != nil {
			u.HtCountryId = u.Hometown.Country.Id
		}
		if u.Hometown.Province != nil {
			u.HtStateId = u.Hometown.Province.Id
		}
		if u.Hometown.City != nil {
			u.HtCityId = u.Hometown.City.Id
		}
		if u.Hometown.District != nil {
			u.HtDistId = u.Hometown.District.Id
		}
	}
	if u.PlaceOfResidence != nil {
		if u.PlaceOfResidence.Country != nil {
			u.PorCountryId = u.PlaceOfResidence.Country.Id
		}
		if u.PlaceOfResidence.Province != nil {
			u.PorStateId = u.PlaceOfResidence.Province.Id
		}
		if u.PlaceOfResidence.City != nil {
			u.PorCityId = u.PlaceOfResidence.City.Id
		}
		if u.PlaceOfResidence.District != nil {
			u.PorDistId = u.PlaceOfResidence.District.Id
		}
	}
	if u.Education != nil {
		u.EduId = u.Education.Id
	}
	if u.Feeling != nil {
		u.FeelingId = u.Feeling.Id
	}
	if u.BloodType != nil {
		u.BloodTypeId = u.BloodTypeId
	}
	if u.Constellation != nil {
		u.ConstellationId = u.Constellation.Id
	}
	return nil
}

// Gorp's lack of support for loading relations automatically.
func (u *UserInfo) PostGet(exe gorp.SqlExecutor) error {

	obj, err := exe.Get(User{}, u.UserId)
	if err != nil {
		return fmt.Errorf("Error loading a UserInfo's User(%d) %s", u.UserId, err)
	}
	u.User = obj.(*User)
	if u.EduId > 0 {
		if u.Education = EducationOf(u.EduId); u.Education == nil {
			return fmt.Errorf("Error EduId => %d", u.EduId)
		}
	}
	if u.FeelingId > 0 {
		if u.Feeling = FeelingOf(u.FeelingId); u.Feeling == nil {
			return fmt.Errorf("Error FeelingId => %d", u.FeelingId)
		}
	}
	if u.BloodTypeId > 0 {
		if u.BloodType = BloodTypeOf(u.BloodTypeId); u.BloodType == nil {
			return fmt.Errorf("Error BloodTypeId => %d", u.BloodTypeId)
		}
	}
	if u.ConstellationId > 0 {
		if u.Constellation = ConstellationOf(u.ConstellationId); u.Constellation == nil {
			return fmt.Errorf("Error ConstellationId => %d", u.ConstellationId)
		}
	}
	if u.GenderCode > 0 {
		if u.Gender = GenderOf(u.GenderCode); u.Gender.Code <= 0 {
			return fmt.Errorf("Error GenderCode => %d", u.GenderCode)
		}
	}
	if u.DateOfBirthStr.Valid {
		if u.DateOfBirth, err = time.Parse("2006-01-02", u.DateOfBirthStr.String); err != nil {
			return fmt.Errorf("Error parsing date of birth '%v'", u.DateOfBirthStr)
		}
	}
	return nil
}

func ToUserInfo(i interface{}, err error) *UserInfo {
	if err != nil {
		panic(err)
	}
	if i == nil || reflect.ValueOf(i).IsNil() {
		return nil
	}
	return i.(*UserInfo)
}

func ToUserInfos(results []interface{}, err error) []*UserInfo {
	if err != nil {
		panic(err)
	}
	size := len(results)
	userInfos := make([]*UserInfo, size)
	if size == 0 {
		return userInfos
	}
	for i, r := range results {
		userInfos[i] = r.(*UserInfo)
	}
	return userInfos
}
