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
	"reflect"
)

// mapped table dict_country
type Country struct {
	Id    uint   `db:"id" json:"id"`
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

func ToCountry(i interface{}, err error) *Country {
	if err != nil {
		panic(err)
	}
	if i == nil || reflect.ValueOf(i).IsNil() {
		return nil
	}
	return i.(*Country)
}

func ToCountries(results []interface{}, err error) []*Country {
	if err != nil {
		panic(err)
	}
	size := len(results)
	countries := make([]*Country, size)
	if size == 0 {
		return countries
	}
	for i, r := range results {
		countries[i] = r.(*Country)
	}
	return countries
}

// province or state
// mapped table `dict_province`
type Province struct {
	Id             uint   `db:"id" json:"id"`
	CountryId      uint   `db:"country_id" json:"countryId"`
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

func ToProvince(i interface{}, err error) *Province {
	if err != nil {
		panic(err)
	}
	if i == nil || reflect.ValueOf(i).IsNil() {
		return nil
	}
	return i.(*Province)
}

func ToProvinces(results []interface{}, err error) []*Province {
	if err != nil {
		panic(err)
	}
	size := len(results)
	provinces := make([]*Province, size)
	if size == 0 {
		return provinces
	}
	for i, r := range results {
		provinces[i] = r.(*Province)
	}
	return provinces
}

// mapped table `dict_city`
type City struct {
	Id        uint   `db:"id" json:"id"`
	CountryId uint   `db:"country_id", json:"countryId"`
	ProId     uint   `db:"pro_id" json:"proId"`
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

func ToCity(i interface{}, err error) *City {
	if err != nil {
		panic(err)
	}
	if i == nil || reflect.ValueOf(i).IsNil() {
		return nil
	}
	return i.(*City)
}

func ToCities(results []interface{}, err error) []*City {
	if err != nil {
		panic(err)
	}
	size := len(results)
	cities := make([]*City, size)
	if size == 0 {
		return cities
	}
	for i, r := range results {
		cities[i] = r.(*City)
	}
	return cities
}

// mapped table `dict_district`
type District struct {
	Id        uint   `db:"id" json:"id"`
	CountryId uint   `db:"country_id" json:"countryId"`
	CityId    uint   `db:"city_id" json:"cityId"`
	ProId     uint   `db:"pro_id" json:"proId"`
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

func ToDistrict(i interface{}, err error) *District {
	if err != nil {
		panic(err)
	}
	if i == nil || reflect.ValueOf(i).IsNil() {
		return nil
	}
	return i.(*District)
}

func ToDistricts(results []interface{}, err error) []*District {
	if err != nil {
		panic(err)
	}
	size := len(results)
	districts := make([]*District, size)
	if size == 0 {
		return districts
	}
	for i, r := range results {
		districts[i] = r.(*District)
	}
	return districts
}

// 所在地、位置信息
type Location struct {
	Country  *Country
	Province *Province
	City     *City
	District *District
}

type Education struct {
	Id    uint16 `json:"id"`
	Name  string `json:"name"`
	EName string `json:"ename"`
	Level uint16 `json:"lv"`
}

// Returns Education object string
func (e Education) String() string {
	return fmt.Sprintf("Education{%d,\"%s\",\"%s\",%d}",
		e.Id, e.Name, e.EName, e.Level)
}

// Education instances
var (
	EDU_Unknown  = &Education{uint16(0), "未选择", "Unknown", uint16(0)}
	educationMap = map[uint16]*Education{
		uint16(1): &Education{uint16(1), "博士后", "Postdoctoral", uint16(9)},
		uint16(2): &Education{uint16(2), "博士", "Doctor", uint16(8)},
		uint16(3): &Education{uint16(3), "硕士", "Master", uint16(7)},
		uint16(4): &Education{uint16(4), "大学本科", "Undergraduate", uint16(6)},
		uint16(5): &Education{uint16(5), "专科", "College", uint16(5)},
		uint16(6): &Education{uint16(6), "高中/职高", "Senior", uint16(4)},
		uint16(7): &Education{uint16(7), "中专", "Secondary", uint16(3)},
		uint16(8): &Education{uint16(8), "初中", "Junior", uint16(2)},
		uint16(9): &Education{uint16(9), "小学", "Primary", uint16(1)},
	}
)

// Returns Education of the specified id.
// if the def is nil, default return EDU_Unknown
func EducationOf(id uint16, def *Education) *Education {
	if edu, ok := educationMap[id]; ok {
		return edu
	}
	if def == nil {
		return EDU_Unknown
	}
	return nil
}

func AllEducations() []*Education {
	results := make([]*Education, len(educationMap))
	i := 0
	for _, v := range educationMap {
		results[i] = v
	}
	return results
}

type Feeling struct {
	Id    uint16 `json:"id"`
	Name  string `json:"name"`
	EName string `json:"ename"`
}

// Returns Feeling object string
func (f Feeling) String() string {
	return fmt.Sprintf("Feeling{%d,\"%s\",\"%s\"}", f.Id, f.Name, f.EName)
}

// Feeling instances
var (
	FL_Unknown = &Feeling{uint16(0), "未选择", "Unknown"}
	feelingMap = map[uint16]*Feeling{
		uint16(1): &Feeling{uint16(1), "单身", "Single"},
		uint16(2): &Feeling{uint16(2), "恋爱中", "In love"},
		uint16(3): &Feeling{uint16(3), "已定婚", "Engadged"},
		uint16(4): &Feeling{uint16(4), "已婚", "Marriaged"},
		uint16(5): &Feeling{uint16(5), "分居", "Separated"},
		uint16(6): &Feeling{uint16(6), "离异", "Divorced"},
		uint16(7): &Feeling{uint16(7), "保密", "Secret"},
	}
)

// Returns Feeling of the specified id.
// if the def is nil, default returns FL_Unknown
func FeelingOf(id uint16, def *Feeling) *Feeling {
	if feeling, ok := feelingMap[id]; ok {
		return feeling
	}
	if def == nil {
		return FL_Unknown
	}
	return def
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
	Id    uint16 `json:"id"`
	Name  string `json:"name"`
	EName string `json:"ename"`
}

// Returns BloodType object string
func (b BloodType) String() string {
	return fmt.Sprintf("BloodType{%d,\"%s\",\"%s\"}", b.Id, b.Name, b.EName)
}

// BloodType instances
var (
	BL_Unknown = &BloodType{uint16(0), "未选择", "Unknown"}
	bloodMap   = map[uint16]*BloodType{
		uint16(1): &BloodType{uint16(1), "A", "A blood group"},
		uint16(2): &BloodType{uint16(2), "B", "B blood group"},
		uint16(3): &BloodType{uint16(3), "AB", "AB blood group"},
		uint16(4): &BloodType{uint16(4), "O", "O blood group"},
		uint16(5): &BloodType{uint16(5), "其他", "Other blood group"},
	}
)

// Returns BloodType of the specified id.
// if the def is nil, default returns BL_Unknown
func BloodTypeOf(id uint16, def *BloodType) *BloodType {
	if bloodType, ok := bloodMap[id]; ok {
		return bloodType
	}
	if def == nil {
		return BL_Unknown
	}
	return def
}

// Returns all BloodType
func AllBloodTypes() []*BloodType {
	results := make([]*BloodType, len(bloodMap))
	i := 0
	for _, v := range bloodMap {
		results[i] = v
	}
	return results
}

type Constellation struct {
	Id         uint16 `json:"id"`
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
	Cons_Unknown     = &Constellation{uint16(0), "未选择", "", "", 0, 0, 0, 0}
	constellationMap = map[uint16]*Constellation{
		uint16(1):  &Constellation{uint16(1), "白羊座", "Ari", "Aries", 3, 21, 4, 20},
		uint16(2):  &Constellation{uint16(2), "金牛座", "Tau", "Taurus", 4, 21, 5, 20},
		uint16(3):  &Constellation{uint16(3), "双子座", "Gem", "Gemini", 5, 21, 6, 21},
		uint16(4):  &Constellation{uint16(4), "巨蟹座", "Cnc", "Cancer", 6, 22, 7, 22},
		uint16(5):  &Constellation{uint16(5), "狮子座", "Leo", "Leo", 7, 23, 8, 22},
		uint16(6):  &Constellation{uint16(6), "处女座", "Vir", "Virgo", 8, 23, 9, 22},
		uint16(7):  &Constellation{uint16(7), "天秤座", "Lib", "Libra", 9, 23, 10, 23},
		uint16(8):  &Constellation{uint16(8), "天蝎座", "Sco", "Scorpius", 10, 24, 11, 21},
		uint16(9):  &Constellation{uint16(9), "射手座", "Sgr", "Sagittarius", 11, 22, 12, 21},
		uint16(10): &Constellation{uint16(10), "摩羯座", "Cap", "Capricornus", 12, 22, 1, 19},
		uint16(11): &Constellation{uint16(11), "水瓶座", "Aqr", "Aquarius", 1, 20, 2, 18},
		uint16(12): &Constellation{uint16(12), "双鱼座", "Psc", "Pisces", 2, 19, 3, 20},
	}
)

// Returns Constellation of the specified id.
func ConstellationOf(id uint16, def *Constellation) *Constellation {
	if c, ok := constellationMap[id]; ok {
		return c
	}
	if def == nil {
		return Cons_Unknown
	}
	return def
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
