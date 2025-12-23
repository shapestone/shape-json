package jsonpath

import (
	"testing"
)

// FuzzJSONPath is a comprehensive fuzzing test for JSONPath parsing and execution
// Run with: go test -fuzz=FuzzJSONPath -fuzztime=30s
func FuzzJSONPath(f *testing.F) {
	// Seed corpus with valid JSONPath examples
	seeds := []string{
		// Basic selectors
		"$",
		"$.store",
		"$.store.book",
		"$.store.book[0]",
		"$[0]",

		// Wildcards
		"$.*",
		"$.store.*",
		"$[*]",
		"$.store.book[*]",

		// Recursive descent
		"$..price",
		"$..book",
		"$..*",

		// Array slicing
		"$[0:5]",
		"$[:10]",
		"$[5:]",
		"$[-3:]",
		"$[:-2]",
		"$[-5:-2]",

		// Bracket notation
		"$['store']",
		`$["store"]`,
		"$['property name with spaces']",
		"$['special-chars!@#']",

		// Filters
		"$[?(@.price < 10)]",
		"$[?(@.price > 5)]",
		"$[?(@.price <= 10)]",
		"$[?(@.price >= 5)]",
		`$[?(@.name == "Alice")]`,
		"$[?(@.active == true)]",
		"$[?(@.value == null)]",
		`$[?(@.email =~ ".*@example.com")]`,

		// Logical operators
		"$[?(@.price < 10 && @.active == true)]",
		"$[?(@.price < 5 || @.price > 20)]",

		// Complex paths
		"$.store.book[0].title",
		"$.store.book[*].author",
		"$.store.book[?(@.price < 10)].title",
		"$.store.book[0:2].author",
		"$.store..price",

		// Edge cases
		"$.''",
		"$['']",
		"$.",
		"$..",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, path string) {
		// The parser should never panic, regardless of input
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("JSONPath parsing panicked on input %q: %v", path, r)
			}
		}()

		// Try to parse the path
		expr, err := ParseString(path)

		// We don't care if parsing fails (most random inputs will be invalid)
		// We only care that it doesn't panic
		if err != nil {
			return
		}

		// If parsing succeeded, try to execute it on some sample data
		// This should also never panic
		testData := map[string]interface{}{
			"store": map[string]interface{}{
				"book": []interface{}{
					map[string]interface{}{
						"price":  10,
						"title":  "Book 1",
						"active": true,
					},
					map[string]interface{}{
						"price": 20,
						"title": "Book 2",
					},
				},
			},
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("JSONPath execution panicked on path %q: %v", path, r)
			}
		}()

		_ = expr.Get(testData)
	})
}

// FuzzJSONPathFilters fuzzes specifically filter expressions
func FuzzJSONPathFilters(f *testing.F) {
	seeds := []string{
		"$[?(@.price < 10)]",
		"$[?(@.price > 5)]",
		"$[?(@.price <= 10)]",
		"$[?(@.price >= 5)]",
		`$[?(@.name == "test")]`,
		"$[?(@.active == true)]",
		"$[?(@.active == false)]",
		"$[?(@.value == null)]",
		`$[?(@.email =~ ".*@test.com")]`,
		"$[?(@.price < 10 && @.active == true)]",
		"$[?(@.price < 5 || @.price > 20)]",
		"$[?((@.price < 10 && @.active == true) || @.featured == true)]",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, path string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Filter parsing panicked on input %q: %v", path, r)
			}
		}()

		expr, err := ParseString(path)
		if err != nil {
			return
		}

		testData := []interface{}{
			map[string]interface{}{"price": 10, "active": true, "name": "test", "email": "user@test.com"},
			map[string]interface{}{"price": 20, "active": false, "value": nil},
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Filter execution panicked on path %q: %v", path, r)
			}
		}()

		_ = expr.Get(testData)
	})
}

// FuzzJSONPathSlices fuzzes array slicing
func FuzzJSONPathSlices(f *testing.F) {
	seeds := []string{
		"$[0:5]",
		"$[:10]",
		"$[5:]",
		"$[-3:]",
		"$[:-2]",
		"$[-5:-2]",
		"$[0:0]",
		"$[100:200]",
		"$[-100:-50]",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, path string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Slice parsing panicked on input %q: %v", path, r)
			}
		}()

		expr, err := ParseString(path)
		if err != nil {
			return
		}

		testData := []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Slice execution panicked on path %q: %v", path, r)
			}
		}()

		_ = expr.Get(testData)
	})
}

// FuzzJSONPathRecursive fuzzes recursive descent
func FuzzJSONPathRecursive(f *testing.F) {
	seeds := []string{
		"$..price",
		"$..book",
		"$..*",
		"$..book[0]",
		"$..book[*]",
		"$..book[?(@.price < 10)]",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, path string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Recursive descent panicked on input %q: %v", path, r)
			}
		}()

		expr, err := ParseString(path)
		if err != nil {
			return
		}

		testData := map[string]interface{}{
			"level1": map[string]interface{}{
				"price": 10,
				"level2": map[string]interface{}{
					"price": 20,
					"book": []interface{}{
						map[string]interface{}{"price": 5},
					},
				},
			},
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Recursive descent execution panicked on path %q: %v", path, r)
			}
		}()

		_ = expr.Get(testData)
	})
}

// FuzzJSONPathBracket fuzzes bracket notation
func FuzzJSONPathBracket(f *testing.F) {
	seeds := []string{
		"$['store']",
		`$["store"]`,
		"$['property with spaces']",
		"$['special-chars']",
		"$['']",
		"$['123']",
		"$['unicode™']",
		"$[0]",
		"$[999]",
		"$[-1]",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, path string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Bracket notation panicked on input %q: %v", path, r)
			}
		}()

		expr, err := ParseString(path)
		if err != nil {
			return
		}

		testData := map[string]interface{}{
			"store":                "value1",
			"property with spaces": "value2",
			"special-chars":        "value3",
			"":                     "empty",
			"123":                  "numeric",
			"unicode™":             "unicode",
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Bracket notation execution panicked on path %q: %v", path, r)
			}
		}()

		_ = expr.Get(testData)
	})
}

// FuzzJSONPathComplex fuzzes complex nested paths
func FuzzJSONPathComplex(f *testing.F) {
	seeds := []string{
		"$.store.book[0].title",
		"$.store.book[*].author",
		"$.store.book[?(@.price < 10)].title",
		"$.store.book[0:2].author",
		"$.store..price",
		"$..book[?(@.price < 10)]",
		"$.store.book[*].author",
		"$.*.book[0]",
		"$..book[*].price",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, path string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Complex path panicked on input %q: %v", path, r)
			}
		}()

		expr, err := ParseString(path)
		if err != nil {
			return
		}

		testData := map[string]interface{}{
			"store": map[string]interface{}{
				"book": []interface{}{
					map[string]interface{}{
						"title":  "Book 1",
						"author": "Author 1",
						"price":  5,
					},
					map[string]interface{}{
						"title":  "Book 2",
						"author": "Author 2",
						"price":  15,
					},
				},
			},
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Complex path execution panicked on path %q: %v", path, r)
			}
		}()

		_ = expr.Get(testData)
	})
}
