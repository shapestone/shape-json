package jsonpath

import (
	"reflect"
	"testing"
)

func TestParseString(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		wantError bool
	}{
		{
			name:      "valid root only",
			query:     "$",
			wantError: false,
		},
		{
			name:      "valid simple path",
			query:     "$.name",
			wantError: false,
		},
		{
			name:      "valid nested path",
			query:     "$.user.name",
			wantError: false,
		},
		{
			name:      "valid array index",
			query:     "$[0]",
			wantError: false,
		},
		{
			name:      "valid wildcard",
			query:     "$.*",
			wantError: false,
		},
		{
			name:      "valid recursive",
			query:     "$..name",
			wantError: false,
		},
		{
			name:      "empty query",
			query:     "",
			wantError: true,
		},
		{
			name:      "invalid - no root",
			query:     "name",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseString(tt.query)
			if tt.wantError {
				if err == nil {
					t.Errorf("ParseString() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("ParseString() unexpected error: %v", err)
			}
		})
	}
}

func TestExprGet(t *testing.T) {
	tests := []struct {
		name  string
		query string
		data  interface{}
		want  []interface{}
	}{
		{
			name:  "root selector",
			query: "$",
			data:  map[string]interface{}{"name": "John"},
			want:  []interface{}{map[string]interface{}{"name": "John"}},
		},
		{
			name:  "child selector",
			query: "$.name",
			data:  map[string]interface{}{"name": "John", "age": 30},
			want:  []interface{}{"John"},
		},
		{
			name:  "nested child selectors",
			query: "$.user.name",
			data: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "John",
					"age":  30,
				},
			},
			want: []interface{}{"John"},
		},
		{
			name:  "array index",
			query: "$.users[1]",
			data: map[string]interface{}{
				"users": []interface{}{"John", "Jane", "Bob"},
			},
			want: []interface{}{"Jane"},
		},
		{
			name:  "array wildcard",
			query: "$.users[*]",
			data: map[string]interface{}{
				"users": []interface{}{"John", "Jane", "Bob"},
			},
			want: []interface{}{"John", "Jane", "Bob"},
		},
		{
			name:  "bracket notation",
			query: "$['name']",
			data:  map[string]interface{}{"name": "John"},
			want:  []interface{}{"John"},
		},
		{
			name:  "array slice",
			query: "$[1:3]",
			data:  []interface{}{"a", "b", "c", "d", "e"},
			want:  []interface{}{"b", "c"},
		},
		{
			name:  "array slice from start",
			query: "$[:3]",
			data:  []interface{}{"a", "b", "c", "d", "e"},
			want:  []interface{}{"a", "b", "c"},
		},
		{
			name:  "array slice to end",
			query: "$[2:]",
			data:  []interface{}{"a", "b", "c", "d", "e"},
			want:  []interface{}{"c", "d", "e"},
		},
		{
			name:  "recursive descent",
			query: "$..name",
			data: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "John",
					"profile": map[string]interface{}{
						"name": "Johnny",
					},
				},
			},
			want: []interface{}{"John", "Johnny"},
		},
		{
			name:  "wildcard",
			query: "$.*",
			data:  map[string]interface{}{"a": 1, "b": 2, "c": 3},
			want:  []interface{}{1, 2, 3},
		},
		{
			name:  "complex path with array and object",
			query: "$.users[0].name",
			data: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{"name": "John", "age": 30},
					map[string]interface{}{"name": "Jane", "age": 25},
				},
			},
			want: []interface{}{"John"},
		},
		{
			name:  "complex path with wildcard",
			query: "$.users[*].name",
			data: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{"name": "John", "age": 30},
					map[string]interface{}{"name": "Jane", "age": 25},
				},
			},
			want: []interface{}{"John", "Jane"},
		},
		{
			name:  "no match",
			query: "$.nonexistent",
			data:  map[string]interface{}{"name": "John"},
			want:  nil,
		},
		{
			name:  "nested arrays",
			query: "$[0][1]",
			data: []interface{}{
				[]interface{}{"a", "b", "c"},
				[]interface{}{"d", "e", "f"},
			},
			want: []interface{}{"b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.query)
			if err != nil {
				t.Fatalf("ParseString() error: %v", err)
			}

			got := expr.Get(tt.data)
			if !slicesEqualUnordered(got, tt.want) {
				t.Errorf("Expr.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExprGetWithComplexData(t *testing.T) {
	// Complex nested structure
	data := map[string]interface{}{
		"store": map[string]interface{}{
			"book": []interface{}{
				map[string]interface{}{
					"category": "reference",
					"author":   "Nigel Rees",
					"title":    "Sayings of the Century",
					"price":    8.95,
				},
				map[string]interface{}{
					"category": "fiction",
					"author":   "Evelyn Waugh",
					"title":    "Sword of Honour",
					"price":    12.99,
				},
				map[string]interface{}{
					"category": "fiction",
					"author":   "Herman Melville",
					"title":    "Moby Dick",
					"isbn":     "0-553-21311-3",
					"price":    8.99,
				},
			},
			"bicycle": map[string]interface{}{
				"color": "red",
				"price": 19.95,
			},
		},
	}

	tests := []struct {
		name  string
		query string
		want  []interface{}
	}{
		{
			name:  "all book authors",
			query: "$.store.book[*].author",
			want:  []interface{}{"Nigel Rees", "Evelyn Waugh", "Herman Melville"},
		},
		{
			name:  "first book",
			query: "$.store.book[0]",
			want: []interface{}{
				map[string]interface{}{
					"category": "reference",
					"author":   "Nigel Rees",
					"title":    "Sayings of the Century",
					"price":    8.95,
				},
			},
		},
		{
			name:  "all prices recursively",
			query: "$..price",
			want:  []interface{}{8.95, 12.99, 8.99, 19.95},
		},
		{
			name:  "bicycle color",
			query: "$.store.bicycle.color",
			want:  []interface{}{"red"},
		},
		{
			name:  "book slice",
			query: "$.store.book[0:2]",
			want: []interface{}{
				map[string]interface{}{
					"category": "reference",
					"author":   "Nigel Rees",
					"title":    "Sayings of the Century",
					"price":    8.95,
				},
				map[string]interface{}{
					"category": "fiction",
					"author":   "Evelyn Waugh",
					"title":    "Sword of Honour",
					"price":    12.99,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.query)
			if err != nil {
				t.Fatalf("ParseString() error: %v", err)
			}

			got := expr.Get(data)
			if !deepSlicesEqualUnordered(got, tt.want) {
				t.Errorf("Expr.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		query string
		data  interface{}
		want  []interface{}
	}{
		{
			name:  "empty object",
			query: "$.*",
			data:  map[string]interface{}{},
			want:  nil,
		},
		{
			name:  "empty array",
			query: "$[*]",
			data:  []interface{}{},
			want:  nil,
		},
		{
			name:  "null value",
			query: "$.value",
			data:  map[string]interface{}{"value": nil},
			want:  []interface{}{nil},
		},
		{
			name:  "boolean values",
			query: "$.*",
			data:  map[string]interface{}{"a": true, "b": false},
			want:  []interface{}{true, false},
		},
		{
			name:  "numeric values",
			query: "$[*]",
			data:  []interface{}{0, 1, 2, 3.14, -5},
			want:  []interface{}{0, 1, 2, 3.14, -5},
		},
		{
			name:  "mixed types in array",
			query: "$[*]",
			data:  []interface{}{"string", 123, true, nil},
			want:  []interface{}{"string", 123, true, nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.query)
			if err != nil {
				t.Fatalf("ParseString() error: %v", err)
			}

			got := expr.Get(tt.data)
			if !deepSlicesEqualUnordered(got, tt.want) {
				t.Errorf("Expr.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function for deep comparison without order
func deepSlicesEqualUnordered(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Create a copy of b to mark matches
	bCopy := make([]bool, len(b))

	for _, aVal := range a {
		found := false
		for j, bVal := range b {
			if !bCopy[j] && reflect.DeepEqual(aVal, bVal) {
				bCopy[j] = true
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
