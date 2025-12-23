package jsonpath_test

import (
	"fmt"
	"log"

	"github.com/shapestone/shape-json/pkg/jsonpath"
)

// Example demonstrates basic JSONPath query usage
func Example() {
	// Sample data
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
			},
		},
	}

	// Parse and execute a query
	expr, err := jsonpath.ParseString("$.store.book[*].author")
	if err != nil {
		log.Fatal(err)
	}

	results := expr.Get(data)
	for _, author := range results {
		fmt.Println(author)
	}
	// Output:
	// Nigel Rees
	// Evelyn Waugh
}

// Example_childSelector demonstrates child property selection
func Example_childSelector() {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  30,
		},
	}

	expr, _ := jsonpath.ParseString("$.user.name")
	results := expr.Get(data)
	fmt.Println(results[0])
	// Output: John
}

// Example_arrayIndex demonstrates array indexing
func Example_arrayIndex() {
	data := map[string]interface{}{
		"items": []interface{}{"apple", "banana", "cherry"},
	}

	expr, _ := jsonpath.ParseString("$.items[1]")
	results := expr.Get(data)
	fmt.Println(results[0])
	// Output: banana
}

// Example_wildcard demonstrates wildcard selection
func Example_wildcard() {
	data := map[string]interface{}{
		"items": []interface{}{"apple", "banana", "cherry"},
	}

	expr, _ := jsonpath.ParseString("$.items[*]")
	results := expr.Get(data)
	for _, item := range results {
		fmt.Println(item)
	}
	// Output:
	// apple
	// banana
	// cherry
}

// Example_recursiveDescent demonstrates recursive descent
func Example_recursiveDescent() {
	data := map[string]interface{}{
		"products": []interface{}{
			map[string]interface{}{
				"name":  "Book",
				"price": 8.95,
			},
			map[string]interface{}{
				"name":  "Pen",
				"price": 12.99,
			},
			map[string]interface{}{
				"name":  "Bicycle",
				"price": 19.95,
			},
		},
	}

	expr, _ := jsonpath.ParseString("$..price")
	results := expr.Get(data)
	for _, price := range results {
		fmt.Println(price)
	}
	// Output:
	// 8.95
	// 12.99
	// 19.95
}

// Example_arraySlice demonstrates array slicing
func Example_arraySlice() {
	data := []interface{}{"a", "b", "c", "d", "e"}

	expr, _ := jsonpath.ParseString("$[1:4]")
	results := expr.Get(data)
	for _, item := range results {
		fmt.Println(item)
	}
	// Output:
	// b
	// c
	// d
}
