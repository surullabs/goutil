package template

import (
	. "github.com/surullabs/goutil/testing"
	. "launchpad.net/gocheck"
)

var noError interface{} = nil

type FieldCheckTestSuite struct{}

var _ = Suite(&FieldCheckTestSuite{})

type fieldCheckTest struct {
	check FieldCheck
	data  MapData
	err   interface{}
}

type noData struct{}

func (noData) Get(string) interface{} { return nil }

type mapData map[string]interface{}

func (m mapData) Get(key string) interface{} {
	value := m[key]
	switch typed := value.(type) {
	case map[string]interface{}:
		return mapData(typed)
	default:
		return value
	}
}

func (*FieldCheckTestSuite) TestFieldCheck(c *C) {
	for _, t := range []fieldCheckTest{
		{FieldCheck{Name: "empty"}, noData{}, noError},
		{FieldCheck{Name: "comment", Text: "{{/* This is a comment */}}"}, noData{}, noError},
		{FieldCheck{Name: "invalid template", Text: "{{}}"}, noData{}, "template: invalid template.*"},
		{FieldCheck{Name: "missing field", Text: "{{.Field1}}"}, noData{}, "missing field Field1"},
		{FieldCheck{Name: "missing nested field", Text: "{{.Field1.Field2}}"}, noData{}, "missing field Field1.Field2"},
		{FieldCheck{Name: "missing nested field with top level present", Text: "{{.Field1.Field2}}"}, mapData{"Field1": mapData{}}, "missing field Field1.Field2"},
		{FieldCheck{Name: "existing field", Text: "{{.Field1}}"}, mapData{"Field1": "value"}, noError},
		{FieldCheck{Name: "existing nested field", Text: "{{.Field1.Field2}}"}, mapData{"Field1": mapData{"Field2": "value"}}, noError},
		{FieldCheck{Name: "field with variable", Text: "{{with $x := .Field1}} $x.Field2 {{end}}"}, mapData{"Field1": mapData{"Field2": "value"}}, noError},
	} {
		err := CheckMissingFields(t.check, t.data)
		c.Log(t.check.Name)
		c.Check(err, NilOrErrorMatches, t.err)
	}
}
