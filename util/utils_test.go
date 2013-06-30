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
	_ "fmt"
	"reflect"
	"testing"
)

func TestDotType(t *testing.T) {
	i, i8, i16, i32, i64 := 1, int8(1), int16(1), int32(1), int64(1)
	ui, ui8, ui16, ui32, ui64 := uint(2), uint8(2), uint16(2), uint32(2), uint64(2)
	f32, f64 := float32(1.1), float64(1.1)

	i_type, i8_type, i16_type, i32_type, i64_type := reflect.TypeOf(i),
		reflect.TypeOf(i8), reflect.TypeOf(i16), reflect.TypeOf(i32), reflect.TypeOf(i64)
	ui_type, ui8_type, ui16_type, ui32_type, ui64_type := reflect.TypeOf(ui),
		reflect.TypeOf(ui8), reflect.TypeOf(ui16), reflect.TypeOf(ui32), reflect.TypeOf(ui64)
	f32_type, f64_type := reflect.TypeOf(f32), reflect.TypeOf(f64)

	if i_type.Kind() != reflect.Int || i8_type.Kind() != reflect.Int8 ||
		i16_type.Kind() != reflect.Int16 || i32_type.Kind() != reflect.Int32 ||
		i64_type.Kind() != reflect.Int64 {
		t.Errorf("Type reflect error: %v %v %v %v %v",
			i_type, i8_type, i16_type, i32_type, i64_type)
	} else {
		t.Log("PASS")
	}
	if ui_type.Kind() != reflect.Uint || ui8_type.Kind() != reflect.Uint8 ||
		ui16_type.Kind() != reflect.Uint16 || ui32_type.Kind() != reflect.Uint32 ||
		ui64_type.Kind() != reflect.Uint64 {
		t.Errorf("uint type reflect error: %v %v %v %v %v",
			ui_type, ui8_type, ui16_type, ui32_type, ui64_type)
	} else {
		t.Log("PASS")
	}

	if f32_type.Kind() != reflect.Float32 || f64_type.Kind() != reflect.Float64 {
		t.Errorf("float type reflect error: %v %v\n", f32_type, f64_type)
	} else {
		t.Log("PASS")
	}
}

func TestComparatorWithInt(t *testing.T) {
	a1_int, b1_int := 1, 3
	val, _ := comparator(a1_int, b1_int)
	if val != -1 {
		t.Errorf("The result is error, should:%v actual: %v", -1, val)
	} else {
		t.Log("PASS")
	}

	a2_int8, b2_int64 := int8(5), int64(1)
	val, _ = comparator(a2_int8, b2_int64)
	if val != 1 {
		t.Errorf("The result is error, should:%v actual:%v", 1, val)
	} else {
		t.Log("PASS")
	}

	a3_int16, b3_int32 := int16(10), int32(10)
	val, _ = comparator(a3_int16, b3_int32)
	if val != 0 {
		t.Errorf("The result is error, should:%v actual:%v", 0, val)
	} else {
		t.Log("PASS")
	}

	a4_int32, b4_int8 := int32(53), int8(10)
	val, _ = comparator(a4_int32, b4_int8)
	if val != 1 {
		t.Errorf("The result is error, should:%v actual:%v", 1, val)
	} else {
		t.Log("PASS")
	}
}
