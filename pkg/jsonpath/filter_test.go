package jsonpath

import (
	"testing"
)

func TestFilterExpressions(t *testing.T) {
	tests := []struct {
		name  string
		query string
		data  interface{}
		want  []interface{}
	}{
		{
			name:  "price less than 10",
			query: "$.books[?(@.price < 10)]",
			data: map[string]interface{}{
				"books": []interface{}{
					map[string]interface{}{"title": "Book A", "price": 8.99},
					map[string]interface{}{"title": "Book B", "price": 12.99},
					map[string]interface{}{"title": "Book C", "price": 5.50},
				},
			},
			want: []interface{}{
				map[string]interface{}{"title": "Book A", "price": 8.99},
				map[string]interface{}{"title": "Book C", "price": 5.50},
			},
		},
		{
			name:  "role equals admin",
			query: "$.users[?(@.role == 'admin')]",
			data: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{"name": "Alice", "role": "admin"},
					map[string]interface{}{"name": "Bob", "role": "user"},
					map[string]interface{}{"name": "Charlie", "role": "admin"},
				},
			},
			want: []interface{}{
				map[string]interface{}{"name": "Alice", "role": "admin"},
				map[string]interface{}{"name": "Charlie", "role": "admin"},
			},
		},
		{
			name:  "boolean equality true",
			query: "$.users[?(@.active == true)]",
			data: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{"name": "Alice", "active": true},
					map[string]interface{}{"name": "Bob", "active": false},
					map[string]interface{}{"name": "Charlie", "active": true},
				},
			},
			want: []interface{}{
				map[string]interface{}{"name": "Alice", "active": true},
				map[string]interface{}{"name": "Charlie", "active": true},
			},
		},
		{
			name:  "greater than",
			query: "$.products[?(@.stock > 20)]",
			data: map[string]interface{}{
				"products": []interface{}{
					map[string]interface{}{"name": "Product A", "stock": 15},
					map[string]interface{}{"name": "Product B", "stock": 25},
					map[string]interface{}{"name": "Product C", "stock": 30},
				},
			},
			want: []interface{}{
				map[string]interface{}{"name": "Product B", "stock": 25},
				map[string]interface{}{"name": "Product C", "stock": 30},
			},
		},
		{
			name:  "less than or equal",
			query: "$.items[?(@.priority <= 2)]",
			data: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"name": "Item A", "priority": 1},
					map[string]interface{}{"name": "Item B", "priority": 3},
					map[string]interface{}{"name": "Item C", "priority": 2},
				},
			},
			want: []interface{}{
				map[string]interface{}{"name": "Item A", "priority": 1},
				map[string]interface{}{"name": "Item C", "priority": 2},
			},
		},
		{
			name:  "not equal",
			query: "$.items[?(@.status != 'active')]",
			data: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"name": "Item A", "status": "active"},
					map[string]interface{}{"name": "Item B", "status": "inactive"},
					map[string]interface{}{"name": "Item C", "status": "pending"},
				},
			},
			want: []interface{}{
				map[string]interface{}{"name": "Item B", "status": "inactive"},
				map[string]interface{}{"name": "Item C", "status": "pending"},
			},
		},
		{
			name:  "field existence",
			query: "$.records[?(@.email)]",
			data: map[string]interface{}{
				"records": []interface{}{
					map[string]interface{}{"name": "Alice", "email": "alice@example.com"},
					map[string]interface{}{"name": "Bob"},
					map[string]interface{}{"name": "Charlie", "email": "charlie@example.com"},
				},
			},
			want: []interface{}{
				map[string]interface{}{"name": "Alice", "email": "alice@example.com"},
				map[string]interface{}{"name": "Charlie", "email": "charlie@example.com"},
			},
		},
		{
			name:  "regex match",
			query: "$.items[?(@.name =~ /Apple/)]",
			data: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"name": "Apple iPhone"},
					map[string]interface{}{"name": "Samsung Galaxy"},
					map[string]interface{}{"name": "Apple iPad"},
				},
			},
			want: []interface{}{
				map[string]interface{}{"name": "Apple iPhone"},
				map[string]interface{}{"name": "Apple iPad"},
			},
		},
		{
			name:  "multiple conditions with AND",
			query: "$.products[?(@.price < 15 && @.inStock == true)]",
			data: map[string]interface{}{
				"products": []interface{}{
					map[string]interface{}{"name": "Product A", "price": 10, "inStock": true},
					map[string]interface{}{"name": "Product B", "price": 20, "inStock": true},
					map[string]interface{}{"name": "Product C", "price": 12, "inStock": false},
					map[string]interface{}{"name": "Product D", "price": 14, "inStock": true},
				},
			},
			want: []interface{}{
				map[string]interface{}{"name": "Product A", "price": 10, "inStock": true},
				map[string]interface{}{"name": "Product D", "price": 14, "inStock": true},
			},
		},
		{
			name:  "filter with no matches",
			query: "$.items[?(@.value > 100)]",
			data: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"value": 10},
					map[string]interface{}{"value": 20},
					map[string]interface{}{"value": 30},
				},
			},
			want: nil,
		},
		{
			name:  "extract field from filtered results",
			query: "$.books[?(@.price < 10)].title",
			data: map[string]interface{}{
				"books": []interface{}{
					map[string]interface{}{"title": "Book A", "price": 8.99},
					map[string]interface{}{"title": "Book B", "price": 12.99},
					map[string]interface{}{"title": "Book C", "price": 5.50},
				},
			},
			want: []interface{}{"Book A", "Book C"},
		},
		{
			name:  "greater than or equal",
			query: "$.items[?(@.score >= 80)]",
			data: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"name": "Item A", "score": 75},
					map[string]interface{}{"name": "Item B", "score": 80},
					map[string]interface{}{"name": "Item C", "score": 90},
				},
			},
			want: []interface{}{
				map[string]interface{}{"name": "Item B", "score": 80},
				map[string]interface{}{"name": "Item C", "score": 90},
			},
		},
		{
			name:  "double quoted string comparison",
			query: `$.users[?(@.role == "admin")]`,
			data: map[string]interface{}{
				"users": []interface{}{
					map[string]interface{}{"name": "Alice", "role": "admin"},
					map[string]interface{}{"name": "Bob", "role": "user"},
				},
			},
			want: []interface{}{
				map[string]interface{}{"name": "Alice", "role": "admin"},
			},
		},
		{
			name:  "boolean false comparison",
			query: "$.items[?(@.disabled == false)]",
			data: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"name": "Item A", "disabled": false},
					map[string]interface{}{"name": "Item B", "disabled": true},
					map[string]interface{}{"name": "Item C", "disabled": false},
				},
			},
			want: []interface{}{
				map[string]interface{}{"name": "Item A", "disabled": false},
				map[string]interface{}{"name": "Item C", "disabled": false},
			},
		},
		{
			name:  "multiple conditions with OR",
			query: "$.items[?(@.status == 'active' || @.status == 'pending')]",
			data: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"name": "Item A", "status": "active"},
					map[string]interface{}{"name": "Item B", "status": "inactive"},
					map[string]interface{}{"name": "Item C", "status": "pending"},
				},
			},
			want: []interface{}{
				map[string]interface{}{"name": "Item A", "status": "active"},
				map[string]interface{}{"name": "Item C", "status": "pending"},
			},
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

func TestFilterEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		data      interface{}
		want      []interface{}
		wantError bool
	}{
		{
			name:  "filter on non-array returns empty",
			query: "$.user[?(@.active == true)]",
			data: map[string]interface{}{
				"user": map[string]interface{}{"name": "Alice", "active": true},
			},
			want: nil,
		},
		{
			name:  "nested field reference in filter",
			query: "$.items[?(@.details.category == 'electronics')]",
			data: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{
						"name":    "Item A",
						"details": map[string]interface{}{"category": "electronics"},
					},
					map[string]interface{}{
						"name":    "Item B",
						"details": map[string]interface{}{"category": "books"},
					},
				},
			},
			want: []interface{}{
				map[string]interface{}{
					"name":    "Item A",
					"details": map[string]interface{}{"category": "electronics"},
				},
			},
		},
		{
			name:  "filter with missing field",
			query: "$.items[?(@.missingField == 'value')]",
			data: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"name": "Item A"},
					map[string]interface{}{"name": "Item B"},
				},
			},
			want: nil,
		},
		{
			name:  "filter with float comparison",
			query: "$.items[?(@.value > 3.14)]",
			data: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"name": "Item A", "value": 2.5},
					map[string]interface{}{"name": "Item B", "value": 4.2},
					map[string]interface{}{"name": "Item C", "value": 3.14},
				},
			},
			want: []interface{}{
				map[string]interface{}{"name": "Item B", "value": 4.2},
			},
		},
		{
			name:  "regex match case sensitive",
			query: "$.items[?(@.name =~ /apple/)]",
			data: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"name": "Apple"},
					map[string]interface{}{"name": "apple"},
					map[string]interface{}{"name": "APPLE"},
				},
			},
			want: []interface{}{
				map[string]interface{}{"name": "apple"},
			},
		},
		{
			name:  "regex match case insensitive",
			query: "$.items[?(@.name =~ /(?i)apple/)]",
			data: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"name": "Apple"},
					map[string]interface{}{"name": "apple"},
					map[string]interface{}{"name": "APPLE"},
					map[string]interface{}{"name": "Banana"},
				},
			},
			want: []interface{}{
				map[string]interface{}{"name": "Apple"},
				map[string]interface{}{"name": "apple"},
				map[string]interface{}{"name": "APPLE"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := ParseString(tt.query)
			if tt.wantError {
				if err == nil {
					t.Errorf("ParseString() expected error, got nil")
				}
				return
			}
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
