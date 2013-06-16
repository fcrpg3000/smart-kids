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
)

func TestBloodTypeOf(t *testing.T) {
	bloodType := BloodTypeOf(6)
	if bloodType != nil {
		t.Error("BloodType length is ", len(bloodMap))
	} else {
		t.Log("BloodType length is ", len(bloodMap))
	}

	bloodType = BloodTypeOf(Blood_A.Id)
	if bloodType == Blood_A {
		t.Log("BloodTypeOf func success.")
	} else {
		t.Error("BloodType id and instance not match.")
	}
	fmt.Println("The bloodType is", bloodType)
}

func TestAllBloodTypes(t *testing.T) {
	bloodTypes := AllBloodTypes()
	if len(bloodTypes) != len(bloodMap) {
		t.Error("All BloodType length is ", len(bloodMap))
	} else {
		t.Log("AllBloodTypes func success.")
	}
	for _, bloodType := range bloodTypes {
		if bloodType != BloodTypeOf(bloodType.Id) {
			t.Error("BloodTypeOf func is error. not match id.")
		}
	}
}
