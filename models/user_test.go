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
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestGender(t *testing.T) {
	if Male != GenderOf(1) {
		t.Error("GenderOf get gender ptr is error, actual: ", GenderOf(1))
	} else {
		t.Log("GenderOf match Male success.")
	}
	if Female != GenderOf(2) {
		t.Error("GenderOf get gender ptr is error, actual: ", GenderOf(2))
	} else {
		t.Log("GenderOf match Female success.")
	}
	if SecretGender != GenderOf(3) {
		t.Error("GenderOf get gender ptr is error, actual: ", GenderOf(3))
	} else {
		t.Log("GenderOf match SecretGender success.")
	}
}

func TestHashPassword(t *testing.T) {
	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(fmt.Sprintf("admin{%s}", "admin")))
	fmt.Printf("SourcePwd: %s, Salt: %s, HashPwd: %x\n", "admin", "admin", sha1Hash.Sum(nil))
	fmt.Printf("HashPwd(16): %s\n", hex.EncodeToString(sha1Hash.Sum(nil)))
}
