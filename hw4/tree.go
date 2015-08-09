package main

import (
	"fmt"
)

type CompanyTree struct {
	node        *CompanyEntry
	left, right *CompanyTree
}

// func NewTree() *CompanyTree {
// 	tree := &CompanyTree{}
// 	tree.node = nil
// 	return tree
// }

// Inserts a CompanyList into the tree
func (ct *CompanyTree) insertList(cl CompanyList) {
	for _, entry := range cl {
		// Using entry is OK because each element
		// is a POINTER to a Company Entry
		ct.insertEntry(entry)
		// fmt.Printf("format", ...)
	}
}

func (t *CompanyTree) insertEntry(ce *CompanyEntry) {
	// Empty node, place ce here
	// fmt.Printf("Tree pointers: %v\n*CompanyEntry: %v\n", t, ce)
	if t.node == nil {
		// t = new(CompanyTree)
		// fmt.Println("NODE IS NULL")
		t.node = ce
		t.left = new(CompanyTree)
		t.right = new(CompanyTree)
		// t.left = NewTree()
		// t.right = NewTree()
		return
	} else {
		if ce.companyName < t.node.companyName {
			t.left.insertEntry(ce)
		} else {
			t.right.insertEntry(ce)
		}
	}

	// Empty node, place ce here
	// if t.node == nil {
	// 	t.node = ce
	// }
}

func (t *CompanyTree) printTree() {
	// good, old fashioned, tree recursion
	if t.node == nil {
		return
	}

	if t.left.node != nil {
		t.left.printTree()
	}

	fmt.Printf("* %s: ", t.node.companyName)
	fmt.Printf("(%f, %f)\n", t.node.latitude, t.node.longitude)

	if t.right.node != nil {
		t.right.printTree()
	}

}
