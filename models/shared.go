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

type Client struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsMobile bool   `json:"isMobile"`
}

// 内置的客户端信息
var (
	WEB           = &Client{1, "普通版", false}
	WEB_PE        = &Client{2, "专业版", false}
	ANDROID       = &Client{3, "Android客户端", true}
	IPHONE        = &Client{4, "iphone客户端", true}
	WINDOWS_PHONE = &Client{5, "Windows.Phone客户端", true}
	IPAD          = &Client{6, "ipad客户端", true}
	FIREFOX       = &Client{7, "FirefoxOS客户端", true}
	clients       = map[int]*Client{
		1: WEB,
		2: WEB_PE,
		3: ANDROID,
		4: IPHONE,
		5: WINDOWS_PHONE,
		6: IPAD,
		7: FIREFOX,
	}
)

func ClientOf(id int) *Client {
	client, exists := clients[id]
	if exists {
		return client
	}
	return nil
}

// Implement Cacheable this interface can be cached.
type Cacheable interface {

	// Returns implement this interface struct object key.
	CacheKey() string
	// Returns implement this interface struct object id list key.
	IdListKey() string
}

// mapped table dict_country
type Country struct {
	Id    int    `db:"id" json:"id"`
	Name  string `db:"c_name" json:"cname"`
	EName string `db:"e_name" json:"ename"`
	Code  string `db:"d_code" json:"code"`

	// Transient
	Provinces []*Province `db:"-" json:"provinces,omitempty"`
}

func (c Country) CacheKey() string {
	return fmt.Sprintf("country:%d", c.Id)
}

func (c Country) IdListKey() string {
	return "country:id:list"
}

// province or state
// mapped table `dict_province`
type Province struct {
	Id             int    `db:"id" json:"id"`
	CountryId      int    `db:"country_id" json:"countryId"`
	Name           string `db:"c_name" json:"cname"`
	ShortName      string `db:"s_name" json:"sname"`
	FullName       string `db:"f_name" json:"fname"`
	Code           string `db:"d_code" json:"code"`
	IsState        bool   `db:"is_state" json:"isState"`
	IsMunicipality bool   `db:"is_municipality" json:"isMunicipality"`

	// Transient
	Country *Country `db:"-" json:"country,omitempty"`
	Cities  []*City  `db:"-" json:"cities,omitempty"`
}

func (p Province) CacheKey() string {
	return fmt.Sprintf("province:%d", p.Id)
}

func (p Province) IdListKey() string {
	return "province:id:list"
}

// mapped table `dict_city`
type City struct {
	Id        int    `db:"id" json:"id"`
	CountryId int    `db:"country_id", json:"countryId"`
	ProId     int    `db:"pro_id" json:"proId"`
	Name      string `db:"c_name" json:"name"`
	Code      string `db:"d_code" json:"code"`

	// Transient
	Country   *Country    `db:"-" json:"country,omitempty"`
	Province  *Province   `db:"-" json:"province,omitempty"`
	Districts []*District `db:"-" json:"districts,omitempty"`
}

func (c City) CacheKey() string {
	return fmt.Sprintf("city:%d", c.Id)
}

func (p City) IdListKey() string {
	return "city:id:list"
}

// mapped table `dict_district`
type District struct {
	Id        int    `db:"id" json:"id"`
	CountryId int    `db:"country_id" json:"countryId"`
	CityId    int    `db:"city_id" json:"cityId"`
	ProId     int    `db:"pro_id" json:"proId"`
	Name      string `db:"c_name" json:"cname"`
	ShortName string `db:"s_name" json:"sname"`
	FullName  string `db:"f_name" json:"fname"`
	Code      string `db:"d_code" json:"code"`

	// Transient
	Country  *Country  `db:"-" json:"country,omitempty"`
	Province *Province `db:"-" json:"province,omitempty"`
	City     *City     `db:"-" json:"city,omitempty"`
}

func (d District) CacheKey() string {
	return fmt.Sprintf("district:%d", d.Id)
}

func (p District) IdListKey() string {
	return "district:id:list"
}

// 所在地、位置信息
type Location struct {
	Country  *Country
	Province *Province
	City     *City
	District *District
}

type Education struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	EName string `json:"ename"`
	Level int    `json:"lv"`
}

// Returns Education object string
func (e Education) String() string {
	return fmt.Sprintf("Education{%d,\"%s\",\"%s\",%d}",
		e.Id, e.Name, e.EName, e.Level)
}

// Education instances
var (
	EduPostdoctoral  = &Education{1, "博士后", "Postdoctoral", 9}
	EduDoctor        = &Education{2, "博士", "Doctor", 8}
	EduMaster        = &Education{3, "硕士", "Master", 7}
	EduUndergraduate = &Education{4, "大学本科", "Undergraduate", 6}
	EduCollege       = &Education{5, "专科", "College", 5}
	EduSenior        = &Education{6, "高中/职高", "Senior", 4}
	EduSecondary     = &Education{7, "中专", "Secondary", 3}
	EduJunior        = &Education{8, "初中", "Junior", 2}
	EduPrimary       = &Education{9, "小学", "Primary", 1}
	educationMap     = map[int]*Education{
		1: EduPostdoctoral,
		2: EduDoctor,
		3: EduMaster,
		4: EduUndergraduate,
		5: EduCollege,
		6: EduSenior,
		7: EduSecondary,
		8: EduJunior,
		9: EduPrimary,
	}
)

// Returns Education of the specified id.
func EducationOf(id int) *Education {
	if edu, ok := educationMap[id]; ok {
		return edu
	}
	return nil
}

type Feeling struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	EName string `json:"ename"`
}

// Returns Feeling object string
func (f Feeling) String() string {
	return fmt.Sprintf("Feeling{%d,\"%s\",\"%s\"}", f.Id, f.Name, f.EName)
}

// Feeling instances
var (
	FL_Single    = &Feeling{1, "单身", "Single"}
	FL_InLove    = &Feeling{2, "恋爱中", "In love"}
	FL_Engadged  = &Feeling{3, "已定婚", "Engadged"}
	FL_Marriaged = &Feeling{4, "已婚", "Marriaged"}
	FL_Separated = &Feeling{5, "分居", "Separated"}
	FL_Divorced  = &Feeling{6, "离异", "Divorced"}
	FL_Secret    = &Feeling{7, "保密", "Secret"}
	feelingMap   = map[int]*Feeling{
		1: FL_Single,
		2: FL_InLove,
		3: FL_Engadged,
		4: FL_Marriaged,
		5: FL_Separated,
		6: FL_Divorced,
		7: FL_Secret,
	}
)

// Returns Feeling of the specified id.
func FeelingOf(id int) *Feeling {
	if feeling, ok := feelingMap[id]; ok {
		return feeling
	}
	return nil
}

// Returns all feelings
func AllFeelings() []*Feeling {
	results := make([]*Feeling, 7)
	i := 0
	for _, feeling := range feelingMap {
		results[i] = feeling
		i++
	}
	return results
}

type BloodType struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	EName string `json:"ename"`
}

// Returns BloodType object string
func (b BloodType) String() string {
	return fmt.Sprintf("BloodType{%d,\"%s\",\"%s\"}", b.Id, b.Name, b.EName)
}

// BloodType instances
var (
	Blood_A     = &BloodType{1, "A", "A blood group"}
	Blood_B     = &BloodType{2, "B", "B blood group"}
	Blood_AB    = &BloodType{3, "AB", "AB blood group"}
	Blood_O     = &BloodType{4, "O", "O blood group"}
	Blood_Other = &BloodType{5, "其他", "Other blood group"}
	bloodMap    = map[int]*BloodType{
		1: Blood_A,
		2: Blood_B,
		3: Blood_AB,
		4: Blood_O,
		5: Blood_Other,
	}
)

// Returns BloodType of the specified id.
func BloodTypeOf(id int) *BloodType {
	if bloodType, ok := bloodMap[id]; ok {
		return bloodType
	}
	return nil
}

// Returns all BloodType
func AllBloodTypes() []*BloodType {
	return []*BloodType{Blood_A, Blood_B, Blood_AB, Blood_O, Blood_Other}
}

type Constellation struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	EName      string `json:"ename"`
	EFName     string `json:"efname"`
	StartMonth int    `json:"s_m"`
	StartDay   int    `json:"s_d"`
	EndMonth   int    `json:"e_m"`
	EndDay     int    `json:"e_d"`
}

// Returns Constellation object string
func (c Constellation) String() string {
	return fmt.Sprintf("Consteelation{%d,\"%s\",\"%s\",(%d.%d - %d.%d)}",
		c.Id, c.Name, c.EFName, c.StartMonth, c.StartDay, c.EndMonth, c.EndDay)
}

// Constellation instances
var (
	C_Ari            = &Constellation{1, "白羊座", "Ari", "Aries", 3, 21, 4, 20}
	C_Tau            = &Constellation{2, "金牛座", "Tau", "Taurus", 4, 21, 5, 20}
	C_Gem            = &Constellation{3, "双子座", "Gem", "Gemini", 5, 21, 6, 21}
	C_Cnc            = &Constellation{4, "巨蟹座", "Cnc", "Cancer", 6, 22, 7, 22}
	C_Leo            = &Constellation{5, "狮子座", "Leo", "Leo", 7, 23, 8, 22}
	C_Vir            = &Constellation{6, "处女座", "Vir", "Virgo", 8, 23, 9, 22}
	C_Lib            = &Constellation{7, "天秤座", "Lib", "Libra", 9, 23, 10, 23}
	C_Sco            = &Constellation{8, "天蝎座", "Sco", "Scorpius", 10, 24, 11, 21}
	C_Sgr            = &Constellation{9, "射手座", "Sgr", "Sagittarius", 11, 22, 12, 21}
	C_Cap            = &Constellation{10, "摩羯座", "Cap", "Capricornus", 12, 22, 1, 19}
	C_Aqr            = &Constellation{11, "水瓶座", "Aqr", "Aquarius", 1, 20, 2, 18}
	C_Psc            = &Constellation{12, "双鱼座", "Psc", "Pisces", 2, 19, 3, 20}
	constellationMap = map[int]*Constellation{
		1:  C_Ari,
		2:  C_Tau,
		3:  C_Gem,
		4:  C_Cnc,
		5:  C_Leo,
		6:  C_Vir,
		7:  C_Lib,
		8:  C_Sco,
		9:  C_Sgr,
		10: C_Cap,
		11: C_Aqr,
		12: C_Psc,
	}
)

// Returns Constellation of the specified id.
func ConstellationOf(id int) *Constellation {
	if c, ok := constellationMap[id]; ok {
		return c
	}
	return nil
}

// Returns all constellation instance.
func AllConstellations() []*Constellation {
	results := make([]*Constellation, 12)
	i := 0
	for _, inst := range constellationMap {
		results[i] = inst
		i++
	}
	return results
}
