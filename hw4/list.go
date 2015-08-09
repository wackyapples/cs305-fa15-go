package main

import (
	"fmt"
)

// CompanyList is a slice to pointers to
// company entries
type CompanyList []*CompanyEntry

// Creates a string from items from a CompanyList
func (cl CompanyList) String() (str string) {
	// Need to mess with this number, but too many
	// entries all being added to the same string
	// causes issues. As in hanging issues.
	if len(cl) > 5000 {
		return "Too many entries to stringify, use printList()"
	}

	for _, entry := range cl {
		str += entry.String()
	}

	return
}

// Prints the items of a CompanyList
// PRINTS, not stringifies.
func (cl CompanyList) printList() {

	for _, entry := range cl {
		fmt.Printf("%v", entry)
	}

	return
}

// Prints the items of a CompanyList inversely.
// Why? This is how the original C assignment
// appeared.
func (cl CompanyList) printListInverse() {

	for i := len(cl) - 1; i >= 0; i-- {
		fmt.Printf("%v", cl[i])
	}

	return
}

// Generates a string from a CompanyEntry,
// if the verbose flag is set, EVERY entry is added
// to the string, otherwise only the company name
// and coordinates are printed.
func (ce *CompanyEntry) String() (str string) {
	// if verbose mode is on, include every field
	// Formatting taken directly from original assignment.
	if *verboseFlag {
		str += fmt.Sprintf("* %s\n  %s\n  %s\n", ce.companyName, ce.companyDescription, ce.website)
		str += fmt.Sprintf("  %s, %s\n", ce.streetAddr, ce.suiteNumber)
		str += fmt.Sprintf("  %s, %s %d\n", ce.city, ce.state, ce.zip)
		str += fmt.Sprintf("  (%f, %f)\n", ce.latitude, ce.longitude)
	} else {
		str += fmt.Sprintf("* %s: ", ce.companyName)
		str += fmt.Sprintf("(%f, %f)\n", ce.latitude, ce.longitude)

	}

	return
}
