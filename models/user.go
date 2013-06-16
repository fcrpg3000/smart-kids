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
	"github.com/coopernurse/gorp"
	"github.com/robfig/revel"
	"regexp"
	"time"
)

var (
	Male         = &genderImpl{1, "男", "M"}
	Female       = &genderImpl{2, "女", "F"}
	SecretGender = &genderImpl{3, "保密", "N"}
	genderMap    = map[int]Gender{
		1: Male,
		2: Female,
		3: SecretGender}
)

func GenderOf(code int) Gender {
	gender, exists := genderMap[code]
	if exists {
		return gender
	}
	return SecretGender
}

type Gender interface {
	// gender code value
	Code() int
	// gender name
	Name() string
	// gender alias
	Alias() string
}

type genderImpl struct {
	code  int
	name  string
	alias string
}

func (g genderImpl) Code() int {
	return g.code
}

func (g genderImpl) Name() string {
	return g.name
}

func (g genderImpl) Alias() string {
	return g.alias
}

func (g genderImpl) String() string {
	return fmt.Sprintf("Gender(%d,%s,%s)", g.code, g.name, g.alias)
}

type User struct {
	UserId           int64     `db:"user_id"`
	UserName         string    `db:"user_name"`
	HashPassword     string    `db:"hash_password"`
	GenderCode       int       `db:"gender"`
	Email            string    `db:"email"`
	CommonlyEmail    string    `db:"commonly_email"`
	CreatedTime      time.Time `db:"created_time"`
	LastModifiedTime time.Time `db:"last_modified_time"`

	// Transient property
	Password string `db:"-"`
	Gender   Gender `db:"-"`
}

func (u *User) String() string {
	return fmt.Sprintf("User(%d, %s, %s, %s)", u.UserId, u.UserName, u.HashPassword, u.Email)
}

func NewUser(userName string, password string, gender Gender, others map[string]interface{}) (*User, error) {
	user := &User{}
	timeNow := time.Now()
	user.UserName = userName
	user.Password = password
	user.Gender = gender
	user.CreatedTime = timeNow
	user.LastModifiedTime = timeNow

	if len(others) == 0 {
		return user, nil
	}
	for key, val := range others {
		switch key {
		case "Email":
			user.Email = val.(string)
			break
		case "CommonlyEmail":
			user.CommonlyEmail = val.(string)
			break
		}
	}
	return user, nil
}

var userNameRegexp = regexp.MustCompile("^\\w+$")

func (user *User) Validate(v *revel.Validation) {
	v.Required(user.Gender)

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
	u.GenderCode = u.Gender.Code()
	return nil
}

func (u *User) PostGet() error {
	u.Gender = GenderOf(u.GenderCode)
	return nil
}

// UserAvatar struct
// ----------------------------------------------------------------------------

type UserAvatar struct {
	UserId           string    `db:"user_id"`
	UserName         string    `db:"user_name"`
	ImageDomain      string    `db:"image_domain"`
	AvatarPath       string    `db:"avatar_path"`       // 150x150 maybe
	SrcAvatarPath    string    `db:"src_avatar_path"`   // source size
	SmallAvatarPath  string    `db:"small_avatar_path"` // 80x80 maybe
	ThumbAvatarPath  string    `db:"thumb_avatar_path"` // 40x40 maybe
	AvatarName       string    `db:"avatar_name"`
	CreatedTime      time.Time `db:"created_time"`
	LastModifiedTime time.Time `db:"last_modified_time"`
}

// Returns user's normal avatar image url.
func (u UserAvatar) AvatarUrl() string {
	return u.avatarUrlInternal(u.AvatarPath)
}

// Returns user's small avatar image url.
func (u UserAvatar) SmallAvatarUrl() string {
	return u.avatarUrlInternal(u.SmallAvatarPath)
}

// Returns user's thumb avatar image url.
func (u UserAvatar) ThumbAvatarUrl() string {
	return u.avatarUrlInternal(u.ThumbAvatarPath)
}

// Returns user's avatar image url of the specified size path.
func (u UserAvatar) avatarUrlInternal(path string) string {
	return fmt.Sprintf("http://%s%s%s", u.ImageDomain, path, u.AvatarName)
}

// BannedUser struct
// ----------------------------------------------------------------------------

type BannedUser struct {
	Id                 int64     `db:"id"`
	UserId             int64     `db:"user_id"`
	UserName           string    `db:"user_name"`
	OperatorId         int64     `db:"operator_id"`
	OperatorName       string    `db:"operator_name"`
	Cause              string    `db:"banned_cause"`
	IsPermanent        bool      `db:"is_permanent"`
	BannedTime         time.Time `db:"banned_time"`
	UnbanTime          time.Time `db:"unban_time"`
	CreatedTime        time.Time `db:"created_time"`
	LastModifiedTime   time.Time `db:"last_modified_time"`
	LastModifiedById   int64     `db:"last_modified_by_id"`
	LastModifiedByName string    `db:"last_modified_by_name"`
}

// BannedUser instance default string
func (b *BannedUser) String() string {
	return fmt.Sprintf(`BannedUser{
	Id=%d,
	Target=(%d, %s),
	Operator(%d, %s),
	Cause="%s",
	Permanent=%v,
	Period=(%v - %v),
	LastModified=(time=%v, id=%d, name=%s)
}`, b.Id, b.UserId, b.UserName, b.OperatorId, b.OperatorName,
		b.Cause, b.IsPermanent, b.BannedTime, b.UnbanTime,
		b.LastModifiedTime, b.LastModifiedById, b.LastModifiedByName)
}

// UserInfo struct
// ----------------------------------------------------------------------------

type UserInfo struct {
	UserId           int64     `db:"user_id"`   // not autoincrement
	UserName         string    `db:"user_name"` // just redundancy field
	Nickname         string    `db:"nickname"`
	GenderCode       int       `db:"gender_code"`
	CalendarMode     int16     `db:"calendar_mode"`
	DateOfBirthStr   string    `db:"date_of_birth"`
	HtCountryId      int       `db:"ht_country_id"`
	HtStateId        int       `db:"ht_state_id"`
	HtCityId         int       `db:"ht_city_id"`
	HtDistId         int       `db:"ht_dist_id"`
	PorCountryId     int       `db:"por_country_id`
	PorStateId       int       `db:"por_state_id"`
	PorCityId        int       `db:"por_city_id"`
	PorDistId        int       `db:"por_dist_id"`
	OtherState       string    `db:"other_state"`
	EduId            int       `db:"edu_id"`
	FeelingId        int       `db:"feeling_id"`
	BloodTypeId      int       `db:"blood_type_id"`
	ConstellationId  int       `db:"constellation_id"`
	CreatedTime      time.Time `db:"created_time"`
	LastModifiedTime time.Time `db:"last_modified_time"`

	// Transient
	Gender           Gender         `db:"-"`
	DateOfBirth      time.Time      `db:"-"`
	User             *User          `db:"-"`
	Hometown         *Location      `db:"-"`
	PlaceOfResidence *Location      `db:"-"`
	Education        *Education     `db:"-"`
	Feeling          *Feeling       `db:"-"`
	BloodType        *BloodType     `db:"-"`
	Constellation    *Constellation `db:"-"`
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
	u.userInfo.Nickname = nickname
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
	if u.Gender != nil {
		u.GenderCode = u.Gender.Code()
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
		if u.Gender = GenderOf(u.GenderCode); u.Gender == nil {
			return fmt.Errorf("Error GenderCode => %d", u.GenderCode)
		}
	}
	if len(u.DateOfBirthStr) > 0 {
		if u.DateOfBirth, err = time.Parse("2006-01-02", u.DateOfBirthStr); err != nil {
			return fmt.Errorf("Error parsing date of birth '%s'", u.DateOfBirthStr)
		}
	}
	return nil
}
