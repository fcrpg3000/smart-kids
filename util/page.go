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
)

var (
	currPageError = errors.New("Current page must not be less than zero!")
	pageSizeError = errors.New("Page size must not be less than zero!")
)

type Pageable struct {
	Current     int   `json:"page"`
	PageSize    int   `json:"pageSize"`
	Offset      int   `json:"offset"`
	Sort        *Sort `json:"sort,omitempty"`
	HasPrevious bool  `json:"hasPrev"`
}

func NewPageable0(current int, pageSize int, sort *Sort) (*Pageable, error) {
	if current <= 0 {
		current = 1
	}
	if pageSize < 0 {
		return nil, pageSizeError
	}
	pageable := &Pageable{Current: current, PageSize: pageSize, Sort: sort}
	pageable.Offset = (pageable.Current - 1) * pageable.PageSize
	pageable.HasPrevious = pageable.Current > 1
	return pageable, nil
}

func NewPageable(current int, pageSize int, direction string, properties []string) (*Pageable, error) {
	sort, err := NewSort(direction, properties)
	if err != nil {
		return nil, err
	}
	return NewPageable0(current, pageSize, sort)
}

// Returns the Pageable requesting the next Page.
func (p *Pageable) Next() *Pageable {
	pageable, _ := NewPageable0(p.Current+1, p.PageSize, p.Sort)
	return pageable
}

// Returns the Pageable requesting the first page.
func (p *Pageable) First() *Pageable {
	if p.Current == 1 {
		return p
	}
	first, _ := NewPageable0(1, p.PageSize, p.Sort)
	return first
}

// Returns the previous Pageable or the first Pageable
// if the current one already is the first one.
func (p *Pageable) PrevOrFirst() *Pageable {
	if p.HasPrevious {
		pageable, _ := NewPageable0(p.Current-1, p.PageSize, p.Sort)
		return pageable
	} else {
		return p
	}
}

type Page struct {
	Current          int           `json:"page"`
	PageSize         int           `json:"size"`
	TotalPages       int           `json:"totalPages"`
	NumberOfElements int           `json:"-"`
	TotalElements    int64         `json:"total"`
	HasPrevPage      bool          `json:"hasPrev"`
	IsFirstPage      bool          `json:"firstPage"`
	HasNextPage      bool          `json:"hasNext"`
	IsLastPage       bool          `json:"lastPage"`
	Content          []interface{} `json:"content"`
	HasContent       bool          `hasContent`
	Sort             *Sort         `json:"sort"`
	pageable         *Pageable     `json:"-"`
}

func (p *Page) initialize() *Page {
	if p.pageable == nil {
		p.Current = 1
		p.PageSize = 0
		p.TotalPages = 1
	} else {
		p.Current = p.pageable.Current
		p.PageSize = p.pageable.PageSize
		p.Sort = p.pageable.Sort
	}
	if p.PageSize > 0 {
		intTotal := int(p.TotalElements)
		if intTotal%p.PageSize == 0 {
			p.TotalPages = intTotal / p.PageSize
		} else {
			p.TotalPages = intTotal/p.PageSize + 1
		}
	}

	p.NumberOfElements = len(p.Content)
	p.HasPrevPage = p.Current > 1
	p.HasNextPage = p.Current+1 <= p.TotalPages
	p.IsFirstPage = !p.HasPrevPage
	p.IsLastPage = !p.HasNextPage
	p.HasContent = p.NumberOfElements > 0

	return p
}

func (p *Page) PageRange() []int {
	if p.TotalPages == 1 {
		return []int{1}
	}
	var pageRange = make([]int, p.TotalPages)
	for i := 1; i <= p.TotalPages; i++ {
		pageRange[i-1] = i
	}
	return pageRange
}

func (p *Page) NextPageable() *Pageable {
	if p.HasNextPage && p.pageable != nil {
		return p.pageable.Next()
	}
	return nil
}

func (p *Page) PrevPageable() *Pageable {
	if p.HasPrevPage && p.pageable != nil {
		return p.pageable.PrevOrFirst()
	}
	return nil
}

func (p *Page) PrevPage() int {
	if p.HasPrevPage {
		return p.Current - 1
	}
	return 0
}

func (p *Page) NextPage() int {
	if p.HasNextPage {
		return p.Current + 1
	}
	return p.Current
}

func NewPage(content []interface{}, pageable *Pageable, total int64) *Page {
	page := &Page{Content: content, pageable: pageable, TotalElements: total}
	return page.initialize()
}

func NewPageWithContent(content []interface{}) *Page {
	return NewPage(content, nil, int64(len(content)))
}
