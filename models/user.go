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
	USER_TABLE          = "sk_user"
	USER_IDENTITY_TABLE = "sk_user_identity"
	USER_DIGITAL_TABLE  = "sk_user_digital"
	USER_AVATAR_TABLE   = "sk_user_avatar"
	USER_INFO_TABLE     = "sk_user_info"
	BANNED_USER_TABLE   = "sk_banned_user"
)

var (
	Male         = Gender{int16(1), "男", "M"}
	Female       = Gender{int16(2), "女", "F"}
	SecretGender = Gender{int16(3), "保密", "N"}
	genderMap    = map[int16]Gender{
		int16(1): Male,
		int16(2): Female,
		int16(3): SecretGender}

	emptyNameAndPwd = errors.New("The userName and password must not be empty.")

	InsufficientBalanceError = errors.New("Insufficient.balance")
	InsufficientScoreError   = errors.New("Insufficient.score")
)

func GenderOf(code int16) Gender {
	gender, exists := genderMap[code]
	if exists {
		return gender
	}
	return SecretGender
}

type Gender struct {
	Code  int16  `json:"code"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

func (g Gender) String() string {
	return fmt.Sprintf("Gender(%d,%s,%s)", g.Code, g.Name, g.Alias)
}

// User table fields
const (
	F_HASH_PASSWORD = "hash_password"
	F_PASSWORD_SALT = "password_salt"
	F_EMAIL         = "email"
	F_SPARE_EMAIL   = "spare_email"
	F_GENDER_CODE   = "gender_code"
)

var (
	UserFields = strings.Join([]string{
		F_USER_ID, F_USER_NAME, F_HASH_PASSWORD, F_PASSWORD_SALT,
		F_EMAIL, F_SPARE_EMAIL, F_CREATED_TIME, F_LAST_MODIFIED_TIME,
	}, ", ")
)

// User model for mapping table `sk_user`
type User struct {
	UserId           uint64         `db:"user_id" json:"uid"`
	Email            string         `db:"email" json:"-"`        // login account
	UserName         string         `db:"user_name" json:"name"` // user name
	HashPassword     string         `db:"hash_password" json:"-"`
	PasswordSalt     string         `db:"password_salt" json:"-"`
	SpareEmail       sql.NullString `db:"spare_email" json:"-"`
	CreatedTime      mysql.NullTime `db:"created_time" json:"created"`
	LastModifiedTime mysql.NullTime `db:"last_modified_time" json:"-"`

	// Transient property
	Password string `db:"-" json:"-"`
}

func (u User) String() string {
	return fmt.Sprintf("User(%d, %s, %s, %s)", u.UserId, u.UserName, u.HashPassword, u.Email)
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

// UserDigital struct
// ----------------------------------------------------------------------------

const (
	F_TOTAL_SCORE           = "total_score"
	F_SCORE                 = "user_score"
	F_GRADE                 = "user_grade"
	F_TOTAL_AMOUNT          = "total_amount"
	F_BALANCE               = "balance"
	F_THREADS               = "threads"
	F_POSTS                 = "posts"
	F_POSTS_REPLIES         = "posts_replies"
	F_NEWS_COMMENTS         = "news_comments"
	F_IMAGE_COMMENTS        = "image_comments"
	F_IMAGE_COMMENT_REPLIES = "image_comment_replies"
)

var (
	UserDigitalFields = strings.Join([]string{
		F_USER_ID, F_USER_NAME, F_TOTAL_SCORE, F_SCORE, F_GRADE,
		F_TOTAL_AMOUNT, F_BALANCE, F_THREADS, F_POSTS, F_POSTS_REPLIES,
		F_NEWS_COMMENTS, F_IMAGE_COMMENTS, F_IMAGE_COMMENT_REPLIES,
		F_LAST_MODIFIED_TIME,
	}, ", ")
)

type UserDigital struct {
	UserId              uint64         `db:"user_id"`
	UserName            string         `db:"user_name"`
	TotalScore          uint64         `db:"total_score"`
	Score               uint64         `db:"user_score"`
	Grade               uint           `db:"user_grade"`
	TotalAmount         uint64         `db:"total_amount"`
	Balance             uint64         `db:"balance"`
	Threads             uint           `db:"threads"`
	Posts               uint           `db:"posts"`
	PostsReplies        uint           `db:"posts_replies"`
	NewsComments        uint           `db:"news_comments"`
	ImageComments       uint           `db:"image_comments"`
	ImageCommentReplies uint           `db:"image_comment_replies"`
	LastModifiedTime    mysql.NullTime `db:"last_modified_time"`
}

func NewDigital(user *User) *UserDigital {
	digital := &UserDigital{UserId: user.UserId, UserName: user.UserName}
	return digital
}

// Add the score, and the total score.
func (u *UserDigital) AddScore(score uint64) *UserDigital {
	u.TotalScore = u.TotalScore + score
	u.Score = u.Score + score
	return u
}

// Add current and total score, and Returns old total score and old current score.
func (u *UserDigital) AddAndGetOldScore(score uint64) (oldTotalScore uint64, oldScore uint64) {
	oldTotalScore = u.TotalScore
	oldScore = u.Score
	u.AddScore(score)
	return oldTotalScore, oldScore
}

// Subtract the score, but not change total score.
func (u *UserDigital) SubtractScore(score uint64) (uint64, error) {
	oldScore := u.Score
	if u.Score < score {
		return oldScore, InsufficientScoreError
	}
	u.Score = u.Score - score
	return oldScore, nil
}

func (u *UserDigital) IncrementGrade() *UserDigital {
	u.Grade = u.Grade + 1
	return u
}

func (u *UserDigital) AddBalance(delta uint64) *UserDigital {
	u.Balance = u.Balance + delta
	u.TotalAmount = u.TotalAmount + delta
	return u
}

func (u *UserDigital) SubtractBalance(subtrahend uint64) (uint64, error) {
	oldBalance := u.Balance
	if u.Balance < subtrahend {
		return oldBalance, InsufficientBalanceError
	}
	u.Balance = u.Balance - subtrahend
	return oldBalance, nil
}

func (u UserDigital) String() string {
	return fmt.Sprintf("UserDigital{UserId=%d, UserName=%s, TotalScore=%d, "+
		"Score=%d, Grade=%d, TotalAmount=%d, Balance=%d, Threads=%d, Posts=%d, "+
		"PostsReplies=%d, NewsComments=%d, ImageComments=%d, ImageCommentReplies=%d, "+
		"LastModifiedTime=%v}", u.UserId, u.UserName, u.TotalScore, u.Score,
		u.Grade, u.TotalAmount, u.Balance, u.Threads, u.Posts, u.PostsReplies,
		u.NewsComments, u.ImageComments, u.ImageCommentReplies, u.LastModifiedTime)
}

// UserIdentity struct
// ----------------------------------------------------------------------------

// UserIdentity table field name constants
const (
	F_READ_NAME     = "real_name"
	F_IDCARD        = "idcard"
	F_IS_VARIFIED   = "is_varified"
	F_ID_IMG_DOMAIN = "id_img_domain"
	F_ID_IMG_PATH   = "id_img_path"
	F_VARIFIED_DATE = "varified_date"
)

// UserIdentity variables
var (
	UserIdentityFields = strings.Join([]string{
		F_USER_ID, F_READ_NAME, F_GENDER_CODE, F_IDCARD, F_IS_VARIFIED,
		F_ID_IMG_DOMAIN, F_ID_IMG_PATH, F_VARIFIED_DATE, F_LAST_MODIFIED_TIME,
	}, ", ")
)

type UserIdentity struct {
	UserId           uint64         `db:"user_id"`
	RealName         string         `db:"real_name"` // user real name
	GenderCode       int16          `db:"gender_code"`
	Idcard           string         `db:"idcard"`      // identity card number
	IsVarified       bool           `db:"is_varified"` // 0 or 1 in db
	IdImgDomain      sql.NullString `db:"id_img_domain"`
	IdImgPath        sql.NullString `db:"id_img_path"`
	VarifiedDate     mysql.NullTime `db:"varified_date"` // 0000-00-00 in db
	LastModifiedTime mysql.NullTime `db:"last_modified_time"`

	Gender Gender `db:"-"`
}

func (u *UserIdentity) String() string {
	return fmt.Sprintf("UserIdentity{UserId=%d, RealName=%s, Idcard=%s, "+
		"IsVarified=%v, IdImgDomain=%v, IdImgPath=%v, VarifiedDate=%v, LastModifiedTime=%v}",
		u.UserId, u.RealName, u.Idcard, u.IsVarified, u.IdImgDomain, u.IdImgPath,
		u.VarifiedDate, u.LastModifiedTime)
}

func (u *UserIdentity) PreInsert(_ gorp.SqlExecutor) error {
	timeNow := time.Now()
	if u.IdImgDomain.Valid && u.IdImgPath.Valid && len(u.Idcard) > 0 &&
		len(u.RealName) > 0 {
		u.IsVarified = true
		u.VarifiedDate = mysql.NullTime{timeNow, true}
	}
	if u.Gender.Code > int16(0) {
		u.GenderCode = u.Gender.Code
	}
	u.LastModifiedTime = mysql.NullTime{timeNow, true}
	return nil
}

func (u *UserIdentity) PostGet(_ gorp.SqlExecutor) error {
	if u.GenderCode > 0 {
		u.Gender = GenderOf(u.GenderCode)
	}
	return nil
}

func (u UserIdentity) IdImageUrl() string {
	if !u.IsVarified {
		return ""
	}
	return fmt.Sprintf("%s%s", u.IdImgDomain, u.IdImgPath)
}

// UserAvatar table field name constants.
const (
	F_IMAGE_DOMAIN      = "image_domain"
	F_AVATAR_PATH       = "avatar_path"
	F_SRC_AVATAR_PATH   = "src_avatar_path"
	F_SMALL_AVATAR_PATH = "small_avatar_path"
	F_THUMB_AVATAR_PATH = "thumb_avatar_path"
	F_AVATAR_NAME       = "avatar_name"
)

// UserAvatar variables
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
	Id               uint64         `db:"id"`
	UserId           uint64         `db:"user_id"`
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

// sk_user_info fields constant
const (
	F_NICKNAME         = "nickname"
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
	UserId           uint64         `db:"user_id" json:"uid"`        // not autoincrement
	UserName         string         `db:"user_name" json:"userName"` // just redundancy field
	Nickname         sql.NullString `db:"nickname" json:"nickname,omitempty"`
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
