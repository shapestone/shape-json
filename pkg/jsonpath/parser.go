package jsonpath

import (
	"fmt"
	"strconv"
	"unicode"
)

// tokenType represents the type of a token in a JSONPath expression
type tokenType int

const (
	tokenRoot tokenType = iota
	tokenDot
	tokenDotDot
	tokenLeftBracket
	tokenRightBracket
	tokenWildcard
	tokenIdentifier
	tokenNumber
	tokenString
	tokenColon
	tokenQuestion
	tokenLeftParen
	tokenRightParen
	tokenAt
	tokenOperator
	tokenFilterExpr
)

// token represents a single token in a JSONPath expression
type token struct {
	value string    // string: 16 bytes
	typ   tokenType // int: 8 bytes
	pos   int       // int: 8 bytes
}

// tokenize converts a JSONPath query string into a sequence of tokens
func tokenize(query string) ([]token, error) {
	if query == "" {
		return nil, fmt.Errorf("empty query string")
	}

	var tokens []token
	pos := 0

	// First token must be root ($)
	if query[0] != '$' {
		return nil, fmt.Errorf("query must start with '$' (root selector)")
	}
	tokens = append(tokens, token{typ: tokenRoot, value: "$", pos: 0})
	pos = 1

	for pos < len(query) {
		ch := query[pos]

		switch {
		case ch == '.':
			// Check for recursive descent (..)
			if pos+1 < len(query) && query[pos+1] == '.' {
				tokens = append(tokens, token{typ: tokenDotDot, value: "..", pos: pos})
				pos += 2
			} else {
				tokens = append(tokens, token{typ: tokenDot, value: ".", pos: pos})
				pos++
			}

		case ch == '[':
			tokens = append(tokens, token{typ: tokenLeftBracket, value: "[", pos: pos})
			pos++

		case ch == ']':
			tokens = append(tokens, token{typ: tokenRightBracket, value: "]", pos: pos})
			pos++

		case ch == '*':
			tokens = append(tokens, token{typ: tokenWildcard, value: "*", pos: pos})
			pos++

		case ch == ':':
			tokens = append(tokens, token{typ: tokenColon, value: ":", pos: pos})
			pos++

		case ch == '?':
			tokens = append(tokens, token{typ: tokenQuestion, value: "?", pos: pos})
			pos++

		case ch == '(':
			tokens = append(tokens, token{typ: tokenLeftParen, value: "(", pos: pos})
			pos++

		case ch == ')':
			tokens = append(tokens, token{typ: tokenRightParen, value: ")", pos: pos})
			pos++

		case ch == '@':
			tokens = append(tokens, token{typ: tokenAt, value: "@", pos: pos})
			pos++

		case ch == '\'':
			// Parse single-quoted string
			start := pos
			pos++ // skip opening quote
			value := ""
			for pos < len(query) && query[pos] != '\'' {
				value += string(query[pos])
				pos++
			}
			if pos >= len(query) {
				return nil, fmt.Errorf("unclosed string at position %d", start)
			}
			pos++ // skip closing quote
			tokens = append(tokens, token{typ: tokenString, value: value, pos: start})

		case ch == '"':
			// Parse double-quoted string
			start := pos
			pos++ // skip opening quote
			value := ""
			for pos < len(query) && query[pos] != '"' {
				value += string(query[pos])
				pos++
			}
			if pos >= len(query) {
				return nil, fmt.Errorf("unclosed string at position %d", start)
			}
			pos++ // skip closing quote
			tokens = append(tokens, token{typ: tokenString, value: value, pos: start})

		case unicode.IsDigit(rune(ch)):
			// Parse number
			start := pos
			value := ""
			for pos < len(query) && unicode.IsDigit(rune(query[pos])) {
				value += string(query[pos])
				pos++
			}
			tokens = append(tokens, token{typ: tokenNumber, value: value, pos: start})

		case unicode.IsLetter(rune(ch)) || ch == '_':
			// Parse identifier
			start := pos
			value := ""
			for pos < len(query) && (unicode.IsLetter(rune(query[pos])) || unicode.IsDigit(rune(query[pos])) || query[pos] == '_') {
				value += string(query[pos])
				pos++
			}
			tokens = append(tokens, token{typ: tokenIdentifier, value: value, pos: start})

		case unicode.IsSpace(rune(ch)):
			// Skip whitespace
			pos++

		case ch == '/':
			// Parse regex pattern /pattern/ or /pattern/flags
			start := pos
			pos++ // skip opening /
			value := "/"
			for pos < len(query) && query[pos] != '/' {
				value += string(query[pos])
				pos++
			}
			if pos >= len(query) {
				return nil, fmt.Errorf("unclosed regex pattern at position %d", start)
			}
			value += "/" // add closing /
			pos++        // skip closing /

			// Check for flags (like 'i' for case-insensitive)
			for pos < len(query) && unicode.IsLetter(rune(query[pos])) {
				value += string(query[pos])
				pos++
			}

			// Keep the full regex pattern with slashes and flags
			tokens = append(tokens, token{typ: tokenString, value: value, pos: start})

		case ch == '<' || ch == '>' || ch == '=' || ch == '!' || ch == '&' || ch == '|':
			// Parse operators - these can be multi-character
			start := pos
			value := string(ch)
			pos++

			// Check for multi-character operators
			if pos < len(query) {
				next := query[pos]
				// Handle ==, !=, <=, >=, &&, ||, =~
				if (ch == '=' && next == '=') ||
					(ch == '!' && next == '=') ||
					(ch == '<' && next == '=') ||
					(ch == '>' && next == '=') ||
					(ch == '&' && next == '&') ||
					(ch == '|' && next == '|') ||
					(ch == '=' && next == '~') {
					value += string(next)
					pos++
				}
			}

			tokens = append(tokens, token{typ: tokenOperator, value: value, pos: start})

		default:
			return nil, fmt.Errorf("unexpected character '%c' at position %d", ch, pos)
		}
	}

	return tokens, nil
}

// parser holds the parsing state
type parser struct {
	tokens []token
	pos    int
}

// parse converts tokens into a compiled expression
func parse(tokens []token) (*expr, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty token list")
	}

	p := &parser{tokens: tokens, pos: 0}
	return p.parseExpr()
}

// parseExpr parses the entire expression
func (p *parser) parseExpr() (*expr, error) {
	// First token must be root
	if !p.match(tokenRoot) {
		return nil, fmt.Errorf("expression must start with root selector ($)")
	}
	p.advance()

	var selectors []selector

	// Root selector itself selects the root
	selectors = append(selectors, &rootSelector{})

	// Parse remaining selectors
	for !p.isAtEnd() {
		sel, err := p.parseSelector()
		if err != nil {
			return nil, err
		}
		selectors = append(selectors, sel)
	}

	return &expr{selectors: selectors}, nil
}

// parseSelector parses a single selector
func (p *parser) parseSelector() (selector, error) {
	current := p.current()

	switch current.typ {
	case tokenDot:
		p.advance()
		return p.parseChildOrWildcard()

	case tokenDotDot:
		p.advance()
		return p.parseRecursive()

	case tokenLeftBracket:
		return p.parseBracket()

	default:
		return nil, fmt.Errorf("unexpected token type %v at position %d", current.typ, current.pos)
	}
}

// parseChildOrWildcard parses a child selector or wildcard after a dot
func (p *parser) parseChildOrWildcard() (selector, error) {
	if p.isAtEnd() {
		return nil, fmt.Errorf("unexpected end of expression after '.'")
	}

	current := p.current()
	if current.typ == tokenWildcard {
		p.advance()
		return &wildcardSelector{}, nil
	}

	if current.typ == tokenIdentifier {
		name := current.value
		p.advance()
		return &childSelector{name: name}, nil
	}

	return nil, fmt.Errorf("expected identifier or wildcard after '.', got %v", current.typ)
}

// parseRecursive parses a recursive descent selector
func (p *parser) parseRecursive() (selector, error) {
	if p.isAtEnd() {
		return nil, fmt.Errorf("unexpected end of expression after '..'")
	}

	current := p.current()
	if current.typ == tokenIdentifier {
		name := current.value
		p.advance()
		return &recursiveSelector{name: name}, nil
	}

	if current.typ == tokenWildcard {
		p.advance()
		return &recursiveWildcardSelector{}, nil
	}

	return nil, fmt.Errorf("expected identifier or wildcard after '..', got %v", current.typ)
}

// parseBracket parses bracket notation [...]
func (p *parser) parseBracket() (selector, error) {
	if !p.match(tokenLeftBracket) {
		return nil, fmt.Errorf("expected '['")
	}
	p.advance()

	if p.isAtEnd() {
		return nil, fmt.Errorf("unexpected end of expression in bracket")
	}

	current := p.current()

	// Handle filter expression [?(...)]
	if current.typ == tokenQuestion {
		return p.parseFilter()
	}

	// Handle wildcard [*]
	if current.typ == tokenWildcard {
		p.advance()
		if !p.match(tokenRightBracket) {
			return nil, fmt.Errorf("expected ']' after '*'")
		}
		p.advance()
		return &wildcardSelector{}, nil
	}

	// Handle string ['name'] or ["name"]
	if current.typ == tokenString {
		name := current.value
		p.advance()
		if !p.match(tokenRightBracket) {
			return nil, fmt.Errorf("expected ']' after string")
		}
		p.advance()
		return &childSelector{name: name}, nil
	}

	// Handle number or slice [0], [0:5], [:5], [2:]
	if current.typ == tokenNumber {
		start, err := strconv.Atoi(current.value)
		if err != nil {
			return nil, fmt.Errorf("invalid number: %v", err)
		}
		p.advance()

		// Check for slice
		if p.match(tokenColon) {
			p.advance()
			return p.parseSlice(start, true)
		}

		// Simple index
		if !p.match(tokenRightBracket) {
			return nil, fmt.Errorf("expected ']' after number")
		}
		p.advance()
		return &indexSelector{index: start}, nil
	}

	// Handle slice starting from beginning [:5]
	if current.typ == tokenColon {
		p.advance()
		return p.parseSlice(0, false)
	}

	return nil, fmt.Errorf("unexpected token in bracket: %v", current.typ)
}

// parseSlice parses the rest of an array slice
func (p *parser) parseSlice(start int, hasStart bool) (selector, error) {
	hasEnd := false
	end := 0

	if !p.match(tokenRightBracket) {
		// There's an end value
		if !p.match(tokenNumber) {
			return nil, fmt.Errorf("expected number or ']' in slice")
		}
		var err error
		end, err = strconv.Atoi(p.current().value)
		if err != nil {
			return nil, fmt.Errorf("invalid number in slice: %v", err)
		}
		hasEnd = true
		p.advance()
	}

	if !p.match(tokenRightBracket) {
		return nil, fmt.Errorf("expected ']' after slice")
	}
	p.advance()

	return &sliceSelector{start: start, end: end, hasStart: hasStart, hasEnd: hasEnd}, nil
}

// Helper methods for parser
func (p *parser) current() token {
	if p.pos >= len(p.tokens) {
		return token{typ: -1}
	}
	return p.tokens[p.pos]
}

func (p *parser) advance() {
	p.pos++
}

func (p *parser) match(typ tokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.current().typ == typ
}

func (p *parser) isAtEnd() bool {
	return p.pos >= len(p.tokens)
}

// parseFilter parses a filter expression [?(...)]
func (p *parser) parseFilter() (selector, error) {
	// Consume the '?'
	if !p.match(tokenQuestion) {
		return nil, fmt.Errorf("expected '?' for filter")
	}
	p.advance()

	// Expect '('
	if !p.match(tokenLeftParen) {
		return nil, fmt.Errorf("expected '(' after '?'")
	}
	p.advance()

	// Collect all tokens until we find the matching ')'
	// This includes handling nested parentheses and quotes
	filterTokens, err := p.collectFilterTokens()
	if err != nil {
		return nil, err
	}

	// Reconstruct the filter expression string from tokens
	filterStr := p.reconstructFilterString(filterTokens)

	// Parse the filter expression
	expr, err := parseFilterExpression(filterStr)
	if err != nil {
		return nil, fmt.Errorf("invalid filter expression: %w", err)
	}

	// Expect ']'
	if !p.match(tokenRightBracket) {
		return nil, fmt.Errorf("expected ']' after filter expression")
	}
	p.advance()

	return &filterSelector{expr: expr}, nil
}

// collectFilterTokens collects tokens until the matching ')' is found
func (p *parser) collectFilterTokens() ([]token, error) {
	var tokens []token
	parenDepth := 1 // We already consumed the opening '('

	for !p.isAtEnd() {
		current := p.current()

		if current.typ == tokenLeftParen {
			parenDepth++
		} else if current.typ == tokenRightParen {
			parenDepth--
			if parenDepth == 0 {
				p.advance() // consume the ')'
				return tokens, nil
			}
		}

		tokens = append(tokens, current)
		p.advance()
	}

	return nil, fmt.Errorf("unclosed filter expression: missing ')'")
}

// reconstructFilterString rebuilds the filter expression string from tokens
func (p *parser) reconstructFilterString(tokens []token) string {
	var parts []string

	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]

		switch tok.typ {
		case tokenAt:
			parts = append(parts, "@")
		case tokenDot:
			parts = append(parts, ".")
		case tokenIdentifier:
			parts = append(parts, tok.value)
		case tokenNumber:
			parts = append(parts, tok.value)
		case tokenString:
			// Check if it's a regex pattern (starts with /)
			if tok.value != "" && tok.value[0] == '/' {
				parts = append(parts, tok.value)
			} else {
				// Re-add quotes for regular strings
				parts = append(parts, "'"+tok.value+"'")
			}
		case tokenOperator:
			parts = append(parts, tok.value)
		default:
			// For other tokens, just use their value
			parts = append(parts, tok.value)
		}
	}

	// Join parts and handle operators specially
	result := ""
	for i, part := range parts {
		// Add space around operators for better parsing
		if isOperator(part) {
			if i > 0 && !isOperator(parts[i-1]) {
				result += " "
			}
			result += part
			if i < len(parts)-1 && !isOperator(parts[i+1]) {
				result += " "
			}
		} else {
			result += part
		}
	}

	return result
}

// isOperator checks if a string is an operator
func isOperator(s string) bool {
	operators := []string{"==", "!=", "<=", ">=", "=~", "<", ">", "&&", "||", "!"}
	for _, op := range operators {
		if s == op {
			return true
		}
	}
	return false
}

// String returns a string representation of a tokenType
func (t tokenType) String() string {
	switch t {
	case tokenRoot:
		return "ROOT"
	case tokenDot:
		return "DOT"
	case tokenDotDot:
		return "DOTDOT"
	case tokenLeftBracket:
		return "LBRACKET"
	case tokenRightBracket:
		return "RBRACKET"
	case tokenWildcard:
		return "WILDCARD"
	case tokenIdentifier:
		return "IDENTIFIER"
	case tokenNumber:
		return "NUMBER"
	case tokenString:
		return "STRING"
	case tokenColon:
		return "COLON"
	default:
		return "UNKNOWN"
	}
}

// Selector implementations

// rootSelector selects the root element
type rootSelector struct{}

func (s *rootSelector) apply(current []interface{}) []interface{} {
	return current
}

// childSelector selects a named child property
type childSelector struct {
	name string
}

func (s *childSelector) apply(current []interface{}) []interface{} {
	var results []interface{}
	for _, item := range current {
		if obj, ok := item.(map[string]interface{}); ok {
			if val, exists := obj[s.name]; exists {
				results = append(results, val)
			}
		}
	}
	return results
}

// indexSelector selects an array element by index
type indexSelector struct {
	index int
}

func (s *indexSelector) apply(current []interface{}) []interface{} {
	var results []interface{}
	for _, item := range current {
		if arr, ok := item.([]interface{}); ok {
			idx := s.index
			if idx < 0 {
				idx = len(arr) + idx
			}
			if idx >= 0 && idx < len(arr) {
				results = append(results, arr[idx])
			}
		}
	}
	return results
}

// wildcardSelector selects all children
type wildcardSelector struct{}

func (s *wildcardSelector) apply(current []interface{}) []interface{} {
	var results []interface{}
	for _, item := range current {
		switch v := item.(type) {
		case map[string]interface{}:
			for _, val := range v {
				results = append(results, val)
			}
		case []interface{}:
			results = append(results, v...)
		}
	}
	return results
}

// sliceSelector selects a slice of array elements
type sliceSelector struct {
	start    int
	end      int
	hasStart bool
	hasEnd   bool
}

func (s *sliceSelector) apply(current []interface{}) []interface{} {
	var results []interface{}
	for _, item := range current {
		arr, ok := item.([]interface{})
		if !ok {
			continue
		}

		start := 0
		if s.hasStart {
			start = s.start
			if start < 0 {
				start = len(arr) + start
			}
			if start < 0 {
				start = 0
			}
		}

		end := len(arr)
		if s.hasEnd {
			end = s.end
			if end < 0 {
				end = len(arr) + end
			}
			if end > len(arr) {
				end = len(arr)
			}
		}

		if start < end && start < len(arr) {
			for i := start; i < end && i < len(arr); i++ {
				results = append(results, arr[i])
			}
		}
	}
	return results
}

// recursiveSelector recursively searches for a named property
type recursiveSelector struct {
	name string
}

func (s *recursiveSelector) apply(current []interface{}) []interface{} {
	var results []interface{}

	var recurse func(interface{}, int)
	recurse = func(item interface{}, depth int) {
		// Limit recursion depth to prevent infinite loops
		if depth > 100 {
			return
		}

		switch v := item.(type) {
		case map[string]interface{}:
			// Check if this object has the property
			if val, exists := v[s.name]; exists {
				results = append(results, val)
			}
			// Recurse into all values
			for _, child := range v {
				recurse(child, depth+1)
			}
		case []interface{}:
			// Recurse into all array elements
			for _, child := range v {
				recurse(child, depth+1)
			}
		}
	}

	for _, item := range current {
		recurse(item, 0)
	}

	return results
}

// recursiveWildcardSelector recursively selects all values
type recursiveWildcardSelector struct{}

func (s *recursiveWildcardSelector) apply(current []interface{}) []interface{} {
	var results []interface{}

	var recurse func(interface{}, int)
	recurse = func(item interface{}, depth int) {
		// Limit recursion depth to prevent infinite loops
		if depth > 100 {
			return
		}

		results = append(results, item)

		switch v := item.(type) {
		case map[string]interface{}:
			for _, child := range v {
				recurse(child, depth+1)
			}
		case []interface{}:
			for _, child := range v {
				recurse(child, depth+1)
			}
		}
	}

	for _, item := range current {
		recurse(item, 0)
	}

	return results
}
