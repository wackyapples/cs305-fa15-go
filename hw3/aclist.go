package main

import (
	"fmt"
)

// type alias for rights, actually only use 4 bits
type right uint8

// Constant for rights values
const (
	R_EXEC right = 1 << iota
	R_WRITE
	R_READ
	R_OWN
)

// Each ACE has the username and associated rights
type accessControlEntry struct {
	user   string
	rights right
}

// Each ACL has a filename and a slice of ACEs
type accessControlList struct {
	filename string
	ace      []accessControlEntry
}

// var maxEntries int = 100

func (acl *accessControlList) initialize(filename string) (ok bool) {
	// set filename
	acl.filename = filename
	// make ace slice, empty for now.
	acl.ace = make([]accessControlEntry, 0)

	return true
}

func (acl *accessControlList) addEntry(newUser string, rights right) (ok bool) {
	// If the list is empty, add the user and return
	if len(acl.ace) == 0 {
		acl.ace = append(acl.ace, accessControlEntry{newUser, rights})
		return true
	}

	// Check for duplicate usernames
	for _, checkEntry := range acl.ace {
		if checkEntry.user == newUser {
			return false
		}
	}

	// Check if the new user is an owner, if they are, add to the front
	// and return.
	if (rights & R_OWN) == R_OWN {
		tempAce := []accessControlEntry{accessControlEntry{newUser, rights}}
		acl.ace = append(tempAce, acl.ace...)

		return true
	}

	// Otherwise, just add the user and rights to the end
	acl.ace = append(acl.ace, accessControlEntry{newUser, rights})

	return true
}

// Creates a string for printing of the contents of an ACL
func (acl *accessControlList) String() (str string) {
	str = fmt.Sprintf("printList: (")

	// empty check
	if acl.ace == nil {
		str += fmt.Sprintf(" empty access control list )\n")
		return
	}

	// Add the filename
	str += fmt.Sprintf("File: %s. ", acl.filename)

	// Add the stringified version of the ACE slice
	// if there are any entries, stringify them
	if len(acl.ace) > 0 {
		for _, entry := range acl.ace {
			str += fmt.Sprintf(", %s (%s)", entry.user, entry.rights)
		}
	} else {
		str += fmt.Sprintf(" No entries.")
	}

	str += fmt.Sprintf(") \n")

	return
}

// Generate a string of rights (one character a piece) from a rights
// variable
func (r right) String() (str string) {
	// if there are no rights return an empty string
	if r == 0 {
		return ""
	}

	// If the right bit for each permission is 'on', add the
	// appropriate letter to the string to return.
	if (r & R_OWN) == R_OWN {
		str += "o"
	}
	if (r & R_READ) == R_READ {
		str += "r"
	}
	if (r & R_WRITE) == R_WRITE {
		str += "w"
	}
	if (r & R_EXEC) == R_EXEC {
		str += "x"
	}

	return
}

// Returns if the right is valid
func (r right) valid() bool {
	return r <= 16 && r >= 0
}

// Returns if the right is valid single right (read OR write, not both, etc)
func (r right) validSingle() bool {
	return r == R_EXEC ||
		r == R_WRITE ||
		r == R_READ ||
		r == R_OWN
}

// Delete a right from a user in an ACL
func (acl *accessControlList) deleteRight(r right, username string) (ok bool) {
	// checks
	// If the thing is empty, fail
	if acl.ace == nil || len(acl.ace) < 1 {
		return false
	}

	// Check if the right is valid
	if !r.validSingle() {
		return false
	}

	// go through the contents and remove the rights
	for idx, entry := range acl.ace {
		if entry.user == username {
			// Clear the bit and return success
			acl.ace[idx].rights &= ^(r)
			return true
		}
	}

	// Fallthrough fail
	return false
}

// Add a right to a user in an ACL
func (acl *accessControlList) addRight(r right, username string) (ok bool) {
	// checks
	// If the thing is empty, fail
	if acl.ace == nil || len(acl.ace) < 1 {
		return false
	}

	// Check if the right is valid
	if !r.validSingle() {
		return false
	}

	// go through the contents and remove the rights
	for idx, entry := range acl.ace {
		if entry.user == username {
			// Clear the bit and return success
			acl.ace[idx].rights |= (r)
			return true
		}
	}

	// Fallthrough fail
	return false

}

// Delete an entry by username from an ACL
func (acl *accessControlList) deleteEntry(username string) (ok bool) {
	// checks
	// If the thing is empty, succeed
	if acl.ace == nil || len(acl.ace) < 1 {
		return true
	}

	// Go through the ACL until the username is found or end is reached
	for idx, entry := range acl.ace {
		if entry.user == username {
			// If the username is there, cut it out, so to speak
			acl.ace = append(acl.ace[:idx], acl.ace[idx+1:]...)
			// Since len(s) is a zero length slice, this is safe
			// Go is garbaged collected, so the now dead element
			// will be disposed of
			return true
		}
	}

	// If the user is not found, they aren't in the list, therefore
	// success!
	return true
}
