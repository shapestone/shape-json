package jsonpath

import (
	"reflect"
	"testing"
)

func TestExecute(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		selectors []selector
		want      []interface{}
	}{
		{
			name:      "root selector only",
			data:      map[string]interface{}{"name": "John"},
			selectors: []selector{&rootSelector{}},
			want:      []interface{}{map[string]interface{}{"name": "John"}},
		},
		{
			name: "child selector",
			data: map[string]interface{}{"name": "John", "age": 30},
			selectors: []selector{
				&rootSelector{},
				&childSelector{name: "name"},
			},
			want: []interface{}{"John"},
		},
		{
			name: "nested child selectors",
			data: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "John",
					"age":  30,
				},
			},
			selectors: []selector{
				&rootSelector{},
				&childSelector{name: "user"},
				&childSelector{name: "name"},
			},
			want: []interface{}{"John"},
		},
		{
			name: "array index selector",
			data: map[string]interface{}{
				"users": []interface{}{"John", "Jane", "Bob"},
			},
			selectors: []selector{
				&rootSelector{},
				&childSelector{name: "users"},
				&indexSelector{index: 1},
			},
			want: []interface{}{"Jane"},
		},
		{
			name: "wildcard selector on object",
			data: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "John",
					"age":  30,
				},
			},
			selectors: []selector{
				&rootSelector{},
				&childSelector{name: "user"},
				&wildcardSelector{},
			},
			want: []interface{}{"John", 30},
		},
		{
			name: "wildcard selector on array",
			data: map[string]interface{}{
				"users": []interface{}{"John", "Jane", "Bob"},
			},
			selectors: []selector{
				&rootSelector{},
				&childSelector{name: "users"},
				&wildcardSelector{},
			},
			want: []interface{}{"John", "Jane", "Bob"},
		},
		{
			name: "slice selector",
			data: []interface{}{"a", "b", "c", "d", "e"},
			selectors: []selector{
				&rootSelector{},
				&sliceSelector{start: 1, end: 4, hasStart: true, hasEnd: true},
			},
			want: []interface{}{"b", "c", "d"},
		},
		{
			name: "slice selector from start",
			data: []interface{}{"a", "b", "c", "d", "e"},
			selectors: []selector{
				&rootSelector{},
				&sliceSelector{end: 3, hasStart: false, hasEnd: true},
			},
			want: []interface{}{"a", "b", "c"},
		},
		{
			name: "slice selector to end",
			data: []interface{}{"a", "b", "c", "d", "e"},
			selectors: []selector{
				&rootSelector{},
				&sliceSelector{start: 2, hasStart: true, hasEnd: false},
			},
			want: []interface{}{"c", "d", "e"},
		},
		{
			name: "recursive selector",
			data: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "John",
					"profile": map[string]interface{}{
						"name": "Johnny",
						"bio":  "Developer",
					},
				},
			},
			selectors: []selector{
				&rootSelector{},
				&recursiveSelector{name: "name"},
			},
			want: []interface{}{"John", "Johnny"},
		},
		{
			name: "no match returns empty",
			data: map[string]interface{}{"name": "John"},
			selectors: []selector{
				&rootSelector{},
				&childSelector{name: "age"},
			},
			want: nil,
		},
		{
			name:      "empty selectors",
			data:      map[string]interface{}{"name": "John"},
			selectors: []selector{},
			want:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := execute(tt.selectors, tt.data)
			if !slicesEqualUnordered(got, tt.want) {
				t.Errorf("execute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChildSelector(t *testing.T) {
	tests := []struct {
		name    string
		current []interface{}
		field   string
		want    []interface{}
	}{
		{
			name: "select existing field",
			current: []interface{}{
				map[string]interface{}{"name": "John", "age": 30},
			},
			field: "name",
			want:  []interface{}{"John"},
		},
		{
			name: "select non-existing field",
			current: []interface{}{
				map[string]interface{}{"name": "John"},
			},
			field: "age",
			want:  nil,
		},
		{
			name: "select from multiple objects",
			current: []interface{}{
				map[string]interface{}{"name": "John"},
				map[string]interface{}{"name": "Jane"},
			},
			field: "name",
			want:  []interface{}{"John", "Jane"},
		},
		{
			name:    "select from non-object",
			current: []interface{}{"string", 123},
			field:   "name",
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sel := &childSelector{name: tt.field}
			got := sel.apply(tt.current)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("childSelector.apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexSelector(t *testing.T) {
	tests := []struct {
		name    string
		current []interface{}
		want    []interface{}
		index   int
	}{
		{
			name:    "select valid index",
			current: []interface{}{[]interface{}{"a", "b", "c"}},
			index:   1,
			want:    []interface{}{"b"},
		},
		{
			name:    "select negative index",
			current: []interface{}{[]interface{}{"a", "b", "c"}},
			index:   -1,
			want:    []interface{}{"c"},
		},
		{
			name:    "select out of bounds index",
			current: []interface{}{[]interface{}{"a", "b", "c"}},
			index:   10,
			want:    nil,
		},
		{
			name:    "select from non-array",
			current: []interface{}{map[string]interface{}{"name": "John"}},
			index:   0,
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sel := &indexSelector{index: tt.index}
			got := sel.apply(tt.current)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("indexSelector.apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWildcardSelector(t *testing.T) {
	tests := []struct {
		name    string
		current []interface{}
		want    []interface{}
	}{
		{
			name: "wildcard on object",
			current: []interface{}{
				map[string]interface{}{"name": "John", "age": 30},
			},
			want: []interface{}{"John", 30},
		},
		{
			name: "wildcard on array",
			current: []interface{}{
				[]interface{}{"a", "b", "c"},
			},
			want: []interface{}{"a", "b", "c"},
		},
		{
			name:    "wildcard on multiple items",
			current: []interface{}{[]interface{}{"a", "b"}, []interface{}{"c", "d"}},
			want:    []interface{}{"a", "b", "c", "d"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sel := &wildcardSelector{}
			got := sel.apply(tt.current)
			if !slicesEqualUnordered(got, tt.want) {
				t.Errorf("wildcardSelector.apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceSelector(t *testing.T) {
	tests := []struct {
		name     string
		current  []interface{}
		want     []interface{}
		start    int
		end      int
		hasStart bool
		hasEnd   bool
	}{
		{
			name:     "slice with start and end",
			current:  []interface{}{[]interface{}{"a", "b", "c", "d", "e"}},
			start:    1,
			end:      4,
			hasStart: true,
			hasEnd:   true,
			want:     []interface{}{"b", "c", "d"},
		},
		{
			name:     "slice from start",
			current:  []interface{}{[]interface{}{"a", "b", "c", "d", "e"}},
			end:      3,
			hasStart: false,
			hasEnd:   true,
			want:     []interface{}{"a", "b", "c"},
		},
		{
			name:     "slice to end",
			current:  []interface{}{[]interface{}{"a", "b", "c", "d", "e"}},
			start:    2,
			hasStart: true,
			hasEnd:   false,
			want:     []interface{}{"c", "d", "e"},
		},
		{
			name:     "slice empty range",
			current:  []interface{}{[]interface{}{"a", "b", "c"}},
			start:    2,
			end:      2,
			hasStart: true,
			hasEnd:   true,
			want:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sel := &sliceSelector{
				start:    tt.start,
				end:      tt.end,
				hasStart: tt.hasStart,
				hasEnd:   tt.hasEnd,
			}
			got := sel.apply(tt.current)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sliceSelector.apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecursiveSelector(t *testing.T) {
	tests := []struct {
		name    string
		current []interface{}
		field   string
		want    []interface{}
	}{
		{
			name: "recursive select at multiple levels",
			current: []interface{}{
				map[string]interface{}{
					"name": "root",
					"child": map[string]interface{}{
						"name": "child1",
						"grandchild": map[string]interface{}{
							"name": "grandchild1",
						},
					},
				},
			},
			field: "name",
			want:  []interface{}{"root", "child1", "grandchild1"},
		},
		{
			name: "recursive select in arrays",
			current: []interface{}{
				map[string]interface{}{
					"items": []interface{}{
						map[string]interface{}{"name": "item1"},
						map[string]interface{}{"name": "item2"},
					},
				},
			},
			field: "name",
			want:  []interface{}{"item1", "item2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sel := &recursiveSelector{name: tt.field}
			got := sel.apply(tt.current)
			if !slicesEqualUnordered(got, tt.want) {
				t.Errorf("recursiveSelector.apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to compare slices without order
func slicesEqualUnordered(a, b []interface{}) bool {
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
