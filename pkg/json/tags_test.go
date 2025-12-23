package json

import (
	"reflect"
	"testing"
)

// TestParseTag tests struct tag parsing functionality
func TestParseTag(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		expected fieldInfo
	}{
		{
			name: "empty tag",
			tag:  "",
			expected: fieldInfo{
				name:      "",
				omitEmpty: false,
				asString:  false,
				skip:      false,
			},
		},
		{
			name: "field name only",
			tag:  "fieldname",
			expected: fieldInfo{
				name:      "fieldname",
				omitEmpty: false,
				asString:  false,
				skip:      false,
			},
		},
		{
			name: "field name with omitempty",
			tag:  "fieldname,omitempty",
			expected: fieldInfo{
				name:      "fieldname",
				omitEmpty: true,
				asString:  false,
				skip:      false,
			},
		},
		{
			name: "field name with string option",
			tag:  "fieldname,string",
			expected: fieldInfo{
				name:      "fieldname",
				omitEmpty: false,
				asString:  true,
				skip:      false,
			},
		},
		{
			name: "field name with omitempty and string",
			tag:  "fieldname,omitempty,string",
			expected: fieldInfo{
				name:      "fieldname",
				omitEmpty: true,
				asString:  true,
				skip:      false,
			},
		},
		{
			name: "skip field with dash",
			tag:  "-",
			expected: fieldInfo{
				name:      "-",
				omitEmpty: false,
				asString:  false,
				skip:      true,
			},
		},
		{
			name: "omitempty only",
			tag:  ",omitempty",
			expected: fieldInfo{
				name:      "",
				omitEmpty: true,
				asString:  false,
				skip:      false,
			},
		},
		{
			name: "string only",
			tag:  ",string",
			expected: fieldInfo{
				name:      "",
				omitEmpty: false,
				asString:  true,
				skip:      false,
			},
		},
		{
			name: "field name with multiple options reordered",
			tag:  "fieldname,string,omitempty",
			expected: fieldInfo{
				name:      "fieldname",
				omitEmpty: true,
				asString:  true,
				skip:      false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseTag(tt.tag)
			if result != tt.expected {
				t.Errorf("parseTag(%q) = %+v, want %+v", tt.tag, result, tt.expected)
			}
		})
	}
}

// TestGetFieldInfo tests field information extraction from struct fields
func TestGetFieldInfo(t *testing.T) {
	type TestStruct struct {
		Name     string `json:"name"`
		Age      int    `json:"age,omitempty"`
		Count    int64  `json:"count,string"`
		Ignored  string `json:"-"`
		NoTag    string
		EmptyTag string `json:""`
		OnlyOmit string `json:",omitempty"`
		BothOpts string `json:"both,omitempty,string"`
	}

	structType := reflect.TypeOf(TestStruct{})

	tests := []struct {
		name      string
		fieldName string
		expected  fieldInfo
	}{
		{
			name:      "simple name tag",
			fieldName: "Name",
			expected: fieldInfo{
				name:      "name",
				omitEmpty: false,
				asString:  false,
				skip:      false,
			},
		},
		{
			name:      "omitempty tag",
			fieldName: "Age",
			expected: fieldInfo{
				name:      "age",
				omitEmpty: true,
				asString:  false,
				skip:      false,
			},
		},
		{
			name:      "string tag",
			fieldName: "Count",
			expected: fieldInfo{
				name:      "count",
				omitEmpty: false,
				asString:  true,
				skip:      false,
			},
		},
		{
			name:      "skip field",
			fieldName: "Ignored",
			expected: fieldInfo{
				name:      "-",
				omitEmpty: false,
				asString:  false,
				skip:      true,
			},
		},
		{
			name:      "no tag uses field name",
			fieldName: "NoTag",
			expected: fieldInfo{
				name:      "NoTag",
				omitEmpty: false,
				asString:  false,
				skip:      false,
			},
		},
		{
			name:      "empty tag uses field name",
			fieldName: "EmptyTag",
			expected: fieldInfo{
				name:      "EmptyTag",
				omitEmpty: false,
				asString:  false,
				skip:      false,
			},
		},
		{
			name:      "only omitempty uses field name",
			fieldName: "OnlyOmit",
			expected: fieldInfo{
				name:      "OnlyOmit",
				omitEmpty: true,
				asString:  false,
				skip:      false,
			},
		},
		{
			name:      "both options",
			fieldName: "BothOpts",
			expected: fieldInfo{
				name:      "both",
				omitEmpty: true,
				asString:  true,
				skip:      false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field, ok := structType.FieldByName(tt.fieldName)
			if !ok {
				t.Fatalf("field %s not found", tt.fieldName)
			}

			result := getFieldInfo(field)
			if result != tt.expected {
				t.Errorf("getFieldInfo(%s) = %+v, want %+v", tt.fieldName, result, tt.expected)
			}
		})
	}
}

// TestIsEmptyValue tests empty value detection
func TestIsEmptyValue(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"zero int", 0, true},
		{"non-zero int", 42, false},
		{"zero int64", int64(0), true},
		{"non-zero int64", int64(42), false},
		{"zero float64", 0.0, true},
		{"non-zero float64", 3.14, false},
		{"empty string", "", true},
		{"non-empty string", "hello", false},
		{"false bool", false, true},
		{"true bool", true, false},
		{"nil pointer", (*int)(nil), true},
		{"non-nil pointer", new(int), false},
		{"nil slice", []int(nil), true},
		{"empty slice", []int{}, true},
		{"non-empty slice", []int{1}, false},
		{"nil map", map[string]int(nil), true},
		{"empty map", map[string]int{}, true},
		{"non-empty map", map[string]int{"a": 1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.value)
			result := isEmptyValue(v)
			if result != tt.expected {
				t.Errorf("isEmptyValue(%v) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}
