// Basic implementation for alphabetic sorting with
// Go's built-in sort package.

package main

// Type alias to implement sort.Interface for
type byName CompanyList

// Implementation for sort.Interface, rather self-explanatory.
func (cl byName) Len() int           { return len(cl) }
func (cl byName) Swap(i, j int)      { cl[i], cl[j] = cl[j], cl[i] }
func (cl byName) Less(i, j int) bool { return cl[i].companyName < cl[j].companyName }
