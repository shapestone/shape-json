package jsonpath

import (
	"fmt"
)

func Example_filterExpression() {
	data := map[string]interface{}{
		"books": []interface{}{
			map[string]interface{}{"title": "The Great Gatsby", "price": 10.99, "inStock": true},
			map[string]interface{}{"title": "1984", "price": 8.99, "inStock": true},
			map[string]interface{}{"title": "To Kill a Mockingbird", "price": 12.99, "inStock": false},
			map[string]interface{}{"title": "Pride and Prejudice", "price": 7.99, "inStock": true},
		},
	}

	// Get books under $10 that are in stock
	expr, _ := ParseString("$.books[?(@.price < 10 && @.inStock == true)].title")
	results := expr.Get(data)

	for _, title := range results {
		fmt.Println(title)
	}
	// Output:
	// 1984
	// Pride and Prejudice
}

func Example_filterWithRegex() {
	data := map[string]interface{}{
		"products": []interface{}{
			map[string]interface{}{"name": "Apple iPhone 15"},
			map[string]interface{}{"name": "Samsung Galaxy S24"},
			map[string]interface{}{"name": "Apple MacBook Pro"},
			map[string]interface{}{"name": "Dell XPS 15"},
		},
	}

	// Get Apple products using regex
	expr, _ := ParseString("$.products[?(@.name =~ /Apple/)].name")
	results := expr.Get(data)

	for _, name := range results {
		fmt.Println(name)
	}
	// Output:
	// Apple iPhone 15
	// Apple MacBook Pro
}

func Example_filterFieldExistence() {
	data := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "Alice", "email": "alice@example.com"},
			map[string]interface{}{"name": "Bob"},
			map[string]interface{}{"name": "Charlie", "email": "charlie@example.com"},
		},
	}

	// Get users that have an email address
	expr, _ := ParseString("$.users[?(@.email)].name")
	results := expr.Get(data)

	for _, name := range results {
		fmt.Println(name)
	}
	// Output:
	// Alice
	// Charlie
}
