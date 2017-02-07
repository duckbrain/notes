
package notebook

// Finds a notebook based on the name. Returns an error if the name is not 
// found or note specific enough to limit to one.
//
// Search matches on the first few characters, similar to how git matches
// commit hashs, if the full name is not provided, the first characters
// are allowed, as long as there are no duplicates
func Search(text string) (Notebook, error) {
	return Notebook{}, nil
}

// Finds all notebooks that can be used and returns them
func All() ([]Notebook, error) {
	return nil, nil
}
