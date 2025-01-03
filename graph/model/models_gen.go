// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Author struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Posts   []*Post   `json:"posts,omitempty"`
	Friends []*Author `json:"friends,omitempty"`
}

type IntFilter struct {
	Equals *int `json:"equals,omitempty"`
	Gt     *int `json:"gt,omitempty"`
	Gte    *int `json:"gte,omitempty"`
	Lt     *int `json:"lt,omitempty"`
	Lte    *int `json:"lte,omitempty"`
}

type Post struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Characters    int      `json:"characters"`
	Text          *string  `json:"text,omitempty"`
	Score         *float64 `json:"score,omitempty"`
	Completed     bool     `json:"completed"`
	DatePublished *string  `json:"datePublished,omitempty"`
	Author        *Author  `json:"author"`
}

type PostAggregateResult struct {
	Posts    []*Post  `json:"posts,omitempty"`
	Count    *string  `json:"count,omitempty"`
	AvgScore *float64 `json:"avgScore,omitempty"`
}

type PostFilter struct {
	Title      *StringFilter `json:"title,omitempty"`
	Characters *IntFilter    `json:"characters,omitempty"`
	IsComplete *bool         `json:"isComplete,omitempty"`
	And        []*PostFilter `json:"and,omitempty"`
	Or         []*PostFilter `json:"or,omitempty"`
	Not        *PostFilter   `json:"not,omitempty"`
}

type PostOrder struct {
	Field *SortableField `json:"field,omitempty"`
	Order *SortOrder     `json:"order,omitempty"`
}

type Query struct {
}

type StringFilter struct {
	Equals   *string `json:"equals,omitempty"`
	Contains *string `json:"contains,omitempty"`
}

type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

var AllSortOrder = []SortOrder{
	SortOrderAsc,
	SortOrderDesc,
}

func (e SortOrder) IsValid() bool {
	switch e {
	case SortOrderAsc, SortOrderDesc:
		return true
	}
	return false
}

func (e SortOrder) String() string {
	return string(e)
}

func (e *SortOrder) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortOrder(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortOrder", str)
	}
	return nil
}

func (e SortOrder) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SortableField string

const (
	SortableFieldTitle         SortableField = "title"
	SortableFieldCharacters    SortableField = "characters"
	SortableFieldDatePublished SortableField = "datePublished"
)

var AllSortableField = []SortableField{
	SortableFieldTitle,
	SortableFieldCharacters,
	SortableFieldDatePublished,
}

func (e SortableField) IsValid() bool {
	switch e {
	case SortableFieldTitle, SortableFieldCharacters, SortableFieldDatePublished:
		return true
	}
	return false
}

func (e SortableField) String() string {
	return string(e)
}

func (e *SortableField) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortableField(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortableField", str)
	}
	return nil
}

func (e SortableField) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
