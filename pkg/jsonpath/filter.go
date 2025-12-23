package jsonpath

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// filterSelector represents a filter expression selector [?(...)]
type filterSelector struct {
	expr *filterExpression
}

// apply applies the filter to array elements
func (s *filterSelector) apply(current []interface{}) []interface{} {
	var results []interface{}

	for _, item := range current {
		// Filter only applies to arrays
		if arr, ok := item.([]interface{}); ok {
			for _, elem := range arr {
				if s.expr.evaluate(elem) {
					results = append(results, elem)
				}
			}
		}
	}

	return results
}

// filterExpression represents a filter condition
type filterExpression struct {
	left     *filterOperand    // pointer: 8 bytes
	right    *filterOperand    // pointer: 8 bytes
	next     *filterExpression // pointer: 8 bytes
	operator string            // string: 16 bytes
	logicOp  string            // "&&" or "||" (string: 16 bytes)
}

// filterOperand represents a value in a filter expression
type filterOperand struct {
	value   interface{} // literal value (string, number, bool) (interface: 16 bytes)
	field   string      // field path like "price" or "details.category" (string: 16 bytes)
	isField bool        // bool: 1 byte (+ 7 bytes padding at end)
}

// evaluate checks if an item matches the filter expression
func (f *filterExpression) evaluate(item interface{}) bool {
	// Evaluate the current expression
	result := f.evaluateSingle(item)

	// If there's a chained expression, combine with logical operator
	if f.next != nil {
		nextResult := f.next.evaluate(item)
		switch f.logicOp {
		case "&&":
			return result && nextResult
		case "||":
			return result || nextResult
		}
	}

	return result
}

// evaluateSingle evaluates a single comparison expression
func (f *filterExpression) evaluateSingle(item interface{}) bool {
	// Get left operand value
	leftVal := f.getOperandValue(f.left, item)

	// Field existence check (no operator)
	if f.operator == "" {
		return leftVal != nil
	}

	// Get right operand value
	rightVal := f.getOperandValue(f.right, item)

	// Perform comparison
	return f.compare(leftVal, rightVal)
}

// getOperandValue extracts the value of an operand
func (f *filterExpression) getOperandValue(op *filterOperand, item interface{}) interface{} {
	if !op.isField {
		return op.value
	}

	// Navigate through nested fields
	current := item
	fields := strings.Split(op.field, ".")

	for _, field := range fields {
		if obj, ok := current.(map[string]interface{}); ok {
			val, exists := obj[field]
			if !exists {
				return nil
			}
			current = val
		} else {
			return nil
		}
	}

	return current
}

// compare performs the comparison based on the operator
func (f *filterExpression) compare(left, right interface{}) bool {
	switch f.operator {
	case "==":
		return compareEqual(left, right)
	case "!=":
		return !compareEqual(left, right)
	case "<":
		return compareLessThan(left, right)
	case ">":
		return compareGreaterThan(left, right)
	case "<=":
		return compareLessThanOrEqual(left, right)
	case ">=":
		return compareGreaterThanOrEqual(left, right)
	case "=~":
		return compareRegex(left, right)
	}
	return false
}

// compareEqual checks equality
func compareEqual(left, right interface{}) bool {
	if left == nil || right == nil {
		return left == right
	}

	// Handle numeric comparisons with type conversion
	leftNum, leftIsNum := toNumber(left)
	rightNum, rightIsNum := toNumber(right)
	if leftIsNum && rightIsNum {
		return leftNum == rightNum
	}

	// Direct equality check for other types
	return left == right
}

// compareLessThan checks if left < right
func compareLessThan(left, right interface{}) bool {
	leftNum, leftOk := toNumber(left)
	rightNum, rightOk := toNumber(right)
	if !leftOk || !rightOk {
		return false
	}
	return leftNum < rightNum
}

// compareGreaterThan checks if left > right
func compareGreaterThan(left, right interface{}) bool {
	leftNum, leftOk := toNumber(left)
	rightNum, rightOk := toNumber(right)
	if !leftOk || !rightOk {
		return false
	}
	return leftNum > rightNum
}

// compareLessThanOrEqual checks if left <= right
func compareLessThanOrEqual(left, right interface{}) bool {
	leftNum, leftOk := toNumber(left)
	rightNum, rightOk := toNumber(right)
	if !leftOk || !rightOk {
		return false
	}
	return leftNum <= rightNum
}

// compareGreaterThanOrEqual checks if left >= right
func compareGreaterThanOrEqual(left, right interface{}) bool {
	leftNum, leftOk := toNumber(left)
	rightNum, rightOk := toNumber(right)
	if !leftOk || !rightOk {
		return false
	}
	return leftNum >= rightNum
}

// compareRegex checks if left matches the regex pattern in right
func compareRegex(left, right interface{}) bool {
	leftStr, ok := left.(string)
	if !ok {
		return false
	}

	pattern, ok := right.(string)
	if !ok {
		return false
	}

	// Remove the leading and trailing slashes from regex pattern
	pattern = strings.TrimPrefix(pattern, "/")
	pattern = strings.TrimSuffix(pattern, "/")

	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	return re.MatchString(leftStr)
}

// toNumber converts a value to float64 if possible
func toNumber(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	case float32:
		return float64(v), true
	case string:
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			return num, true
		}
	}
	return 0, false
}

// parseFilterExpression parses a filter expression string
// Example: "@.price < 10" or "@.price < 15 && @.inStock == true"
func parseFilterExpression(filterStr string) (*filterExpression, error) {
	// First, check for logical operators (&&, ||) to split the expression
	if idx := findLogicalOperator(filterStr); idx != -1 {
		op := ""
		splitIdx := idx
		if strings.HasPrefix(filterStr[idx:], "&&") {
			op = "&&"
		} else if strings.HasPrefix(filterStr[idx:], "||") {
			op = "||"
		}

		left := strings.TrimSpace(filterStr[:splitIdx])
		right := strings.TrimSpace(filterStr[splitIdx+2:])

		leftExpr, err := parseSimpleFilterExpression(left)
		if err != nil {
			return nil, err
		}

		rightExpr, err := parseFilterExpression(right)
		if err != nil {
			return nil, err
		}

		leftExpr.logicOp = op
		leftExpr.next = rightExpr
		return leftExpr, nil
	}

	return parseSimpleFilterExpression(filterStr)
}

// findLogicalOperator finds the index of the first logical operator (&&, ||)
// that is not inside quotes
func findLogicalOperator(s string) int {
	inSingleQuote := false
	inDoubleQuote := false

	for i := 0; i < len(s)-1; i++ {
		ch := s[i]

		if ch == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
		} else if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
		}

		if !inSingleQuote && !inDoubleQuote {
			if s[i] == '&' && s[i+1] == '&' {
				return i
			}
			if s[i] == '|' && s[i+1] == '|' {
				return i
			}
		}
	}

	return -1
}

// parseSimpleFilterExpression parses a simple filter expression (no logical operators)
// Example: "@.price < 10" or "@.role == 'admin'" or "@.email"
func parseSimpleFilterExpression(filterStr string) (*filterExpression, error) {
	filterStr = strings.TrimSpace(filterStr)

	// Check for existence expression (just a field reference)
	if !strings.Contains(filterStr, "==") &&
		!strings.Contains(filterStr, "!=") &&
		!strings.Contains(filterStr, "<=") &&
		!strings.Contains(filterStr, ">=") &&
		!strings.Contains(filterStr, "=~") &&
		!strings.Contains(filterStr, "<") &&
		!strings.Contains(filterStr, ">") {
		// Field existence check
		field := parseFieldReference(filterStr)
		if field == "" {
			return nil, fmt.Errorf("invalid field reference: %s", filterStr)
		}
		return &filterExpression{
			left: &filterOperand{
				isField: true,
				field:   field,
			},
		}, nil
	}

	// Parse comparison operators
	operators := []string{"==", "!=", "<=", ">=", "=~", "<", ">"}
	for _, op := range operators {
		idx := findOperator(filterStr, op)
		if idx == -1 {
			continue
		}

		left := strings.TrimSpace(filterStr[:idx])
		right := strings.TrimSpace(filterStr[idx+len(op):])

		leftOp, err := parseOperand(left)
		if err != nil {
			return nil, fmt.Errorf("invalid left operand: %w", err)
		}

		rightOp, err := parseOperand(right)
		if err != nil {
			return nil, fmt.Errorf("invalid right operand: %w", err)
		}

		return &filterExpression{
			left:     leftOp,
			operator: op,
			right:    rightOp,
		}, nil
	}

	return nil, fmt.Errorf("invalid filter expression: %s", filterStr)
}

// findOperator finds the index of an operator, skipping operators inside quotes
func findOperator(s, op string) int {
	inSingleQuote := false
	inDoubleQuote := false

	for i := 0; i <= len(s)-len(op); i++ {
		ch := s[i]

		if ch == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
		} else if ch == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
		}

		if !inSingleQuote && !inDoubleQuote {
			if strings.HasPrefix(s[i:], op) {
				// Make sure it's not part of a longer operator
				// e.g., don't match "<" when it's part of "<="
				if op == "<" || op == ">" {
					if i+len(op) < len(s) && s[i+len(op)] == '=' {
						continue
					}
				}
				return i
			}
		}
	}

	return -1
}

// parseOperand parses a filter operand (field reference or literal value)
func parseOperand(s string) (*filterOperand, error) {
	s = strings.TrimSpace(s)

	// Check if it's a field reference (starts with @)
	if strings.HasPrefix(s, "@") {
		field := parseFieldReference(s)
		if field == "" {
			return nil, fmt.Errorf("invalid field reference: %s", s)
		}
		return &filterOperand{
			isField: true,
			field:   field,
		}, nil
	}

	// Check if it's a string literal
	if (strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'")) ||
		(strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`)) {
		return &filterOperand{
			isField: false,
			value:   s[1 : len(s)-1],
		}, nil
	}

	// Check if it's a regex pattern
	if strings.HasPrefix(s, "/") && strings.HasSuffix(s, "/") {
		return &filterOperand{
			isField: false,
			value:   s, // Keep the slashes for regex matching
		}, nil
	}

	// Check for regex pattern with flags (e.g., /pattern/i)
	if strings.HasPrefix(s, "/") {
		lastSlash := strings.LastIndex(s[1:], "/")
		if lastSlash != -1 {
			pattern := s[:lastSlash+2]
			flags := s[lastSlash+2:]
			// Convert flags to regex syntax
			if flags != "" {
				// For now, support case insensitive flag
				if strings.Contains(flags, "i") {
					pattern = "/(?i)" + pattern[1:]
				}
			}
			return &filterOperand{
				isField: false,
				value:   pattern,
			}, nil
		}
	}

	// Check if it's a boolean
	if s == "true" {
		return &filterOperand{
			isField: false,
			value:   true,
		}, nil
	}
	if s == "false" {
		return &filterOperand{
			isField: false,
			value:   false,
		}, nil
	}

	// Try to parse as number
	if num, err := strconv.ParseFloat(s, 64); err == nil {
		return &filterOperand{
			isField: false,
			value:   num,
		}, nil
	}

	return nil, fmt.Errorf("invalid operand: %s", s)
}

// parseFieldReference extracts the field path from a field reference
// Example: "@.price" -> "price", "@.details.category" -> "details.category"
func parseFieldReference(s string) string {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "@") {
		return ""
	}
	s = strings.TrimPrefix(s, "@")
	s = strings.TrimPrefix(s, ".")
	return s
}
