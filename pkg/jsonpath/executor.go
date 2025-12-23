package jsonpath

// execute applies a sequence of selectors to data and returns all matching values
func execute(selectors []selector, data interface{}) []interface{} {
	if len(selectors) == 0 {
		return nil
	}

	// Start with the root data as the initial current set
	current := []interface{}{data}

	// Apply each selector in sequence
	for _, sel := range selectors {
		current = sel.apply(current)
		if len(current) == 0 {
			// No matches, early exit
			return nil
		}
	}

	return current
}
