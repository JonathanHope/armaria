package order

// Initial gets the initial order for an item in a list.
// Basically it is the midpoint of the address space.
func Initial() (string, error) {
	base62, err := MakeSymbolTable(Base62)
	if err != nil {
		return "", err
	}

	order, err := base62.Between("", "", 1, 10000, 4)
	if err != nil {
		return "", err
	}

	return order[0], nil
}

// Start gets the order for an item at the start of a list.
func Start(next string) (string, error) {
	base62, err := MakeSymbolTable(Base62)
	if err != nil {
		return "", err
	}

	order, err := base62.Between("", next, 1, 10000, 4)
	if err != nil {
		return "", err
	}

	return order[0], nil
}

// End gets the order for an item at the end of the list.
func End(previous string) (string, error) {
	base62, err := MakeSymbolTable(Base62)
	if err != nil {
		return "", err
	}

	order, err := base62.Between(previous, "", 1, 10000, 4)
	if err != nil {
		return "", err
	}

	return order[0], nil
}

// Between gets the order for an item between two items on the list.
func Between(previous string, next string) (string, error) {
	base62, err := MakeSymbolTable(Base62)
	if err != nil {
		return "", err
	}

	order, err := base62.Between(previous, next, 1, 0, 4)
	if err != nil {
		return "", err
	}

	return order[0], nil
}
