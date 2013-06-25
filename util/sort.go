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
	"errors"
	"fmt"
	"strings"
)

const (
	ASC                 = "ASC"
	DESC                = "DESC"
	DEFAULT_DIRECTION   = ASC
	DEFAULT_IGNORE_CASE = false
)

var (
	emptyPropertyErr = errors.New("You have to provide at least one sort property to sort by!")
)

type Order struct {
	// ASC or DESC
	Direction  string
	Property   string
	IgnoreCase bool
}

func NewOrder(direction, property string, ignoreCase bool) *Order {
	return &Order{direction, property, ignoreCase}
}

func DescendingOrder(property string) *Order {
	return &Order{DESC, property, DEFAULT_IGNORE_CASE}
}

func AscendingOrder(property string) *Order {
	return &Order{DEFAULT_DIRECTION, property, DEFAULT_IGNORE_CASE}
}

// Returns true if this is a ascending order, otherwise false.
func (o Order) IsAscending() bool {
	return o.Direction == ASC
}

func (o Order) Equals(other *Order) bool {
	if other == nil {
		return false
	}
	return o.Direction == other.Direction && o.Property == other.Property
}

// This method use fmt.Sprintf(format, Order#Property, Order#Direction) to return.
func (o Order) Format(format string) string {
	return fmt.Sprintf(format, o.Property, o.Direction)
}

func (o Order) String() string {
	return o.Format("%s: %s")
}

// Sort option for queries. You have to provide at least a list of properties
// to sort for that must not include nil or empty strings.
// The direction defaults to DEFAULT_DIRECTION.
type Sort struct {
	Orders []*Order `json:",omitempty"`
}

func NewSort(direction string, properties []string) (*Sort, error) {
	propLen := len(properties)
	if propLen == 0 {
		return nil, emptyPropertyErr
	}
	var orders = make([]*Order, propLen)
	for i, property := range properties {
		orders[i] = NewOrder(direction, property, DEFAULT_IGNORE_CASE)
	}
	return &Sort{orders}, nil
}

func AscendingSort(properties []string) (*Sort, error) {
	return NewSort(DEFAULT_DIRECTION, properties)
}

func DescendingSort(properties []string) (*Sort, error) {
	return NewSort(DESC, properties)
}

// Returns the order registered for the given property.
func (s Sort) OrderFor(property string) *Order {
	if len(s.Orders) == 0 {
		return nil
	}
	for _, order := range s.Orders {
		if order.Property == property {
			return order
		}
	}
	return nil
}

func (s Sort) SqlString() string {
	var a = make([]string, len(s.Orders)+1)
	a[0] = " ORDER BY "
	for i, order := range s.Orders {
		if i > 0 {
			a[i+1] = order.Format(", %s %s")
		} else {
			a[i+1] = order.Format("%s %s")
		}
	}
	return strings.Join(a, "")
}

func (s Sort) String() string {
	var a = make([]string, len(s.Orders))
	for i, order := range s.Orders {
		a[i] = order.String()
	}
	return strings.Join(a, ", ")
}
