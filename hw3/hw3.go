package main

import (
	"fmt"
	"io"
	"os"
)

var acl accessControlList

func main() {

	paramsCheck()

	// Parse the first file
	if err := parseInputFile(os.Args[1]); err != nil {
		fmt.Println(err)
		fmt.Printf("File parsing failed. Exiting program. \n")
		// Go is garbaged collected
		os.Exit(2)
	}

	// Print the inital ACL
	fmt.Printf("%v\n", &acl)

	// Parse the second input file, altering the access control list
	// that was created by the first input file.
	if err := parseCommandFile(os.Args[2]); err != nil {
		fmt.Println(err)
		fmt.Printf("Command parsing failed. Exiting program. \n")
		// Go is garbaged collected
		os.Exit(2)
	}

	// Print the resulting ACL
	fmt.Printf("%v\n", &acl)

}

// Panic helper
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func paramsCheck() {
	if len(os.Args) != 3 {
		fmt.Printf("%s error: incorrect number of parameters.\n"+
			"usage: %s <aclFile> <commandFile>\n", os.Args[0], os.Args[0])

		os.Exit(2)
	}
}

/* parseInputFile
 *
 * Description: This function reads an input file and constructs an
 *              access control list from the contents of the file.  The
 *              expected input file format is:
 *
 *     f: <filename>
 *     user1
 *     <integer describing rights for user1>
 *     user2
 *     <integer describing rights for user2>
 *     ...
 *     userN
 *     <integer describing rights for userN>
 *
 *  The set of rights available are own, read, write, and execute and
 *  are expressed in four bits.  For example, if a user owns and may read
 *  the file, his set of rights is expressed by 0b1100 or 12.
 *
 *  Only one access control list is created as a result of this function call.
 *
 *  Usernames may not begin with a colon or an asterix.
 *
 */
func parseInputFile(filename string) (err error) {

	// attempt to open file
	file, err := os.Open(filename)
	// data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Printf("%s was successfully opened.\nParsing access control entries from file.\n", filename)

	// n, err := fmt.Sscanf(string(data), "%c")
	var a rune
	var d right
	var word string

	for {
		_, err := fmt.Fscanf(file, "%c", &a)
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		}
		// fmt.Printf("%d: '%c' and %v\n", n, a, err)

		// If a colon is read, the next word is the file for the ACL
		if a == ':' {
			// Attempt to read the next word, bail if fail
			fmt.Fscanf(file, "%s", &word)

			// fmt.Printf("%d: '%s' and %v\n", m, word, err)

			// initializeACL(word, &aclist)
			acl.initialize(word)
		} else if a == '*' {
			// if there's a '*' then next word is username
			_, err1 := fmt.Fscanf(file, "%s\n", &word)
			if err1 != nil {
				return err1
			}
			_, err2 := fmt.Fscanf(file, "%d", &d)
			if err2 != nil {
				return err2
			}

			// fmt.Printf("user: %s\nrights: %d\n", word, d)

			acl.addEntry(word, d)

		}
	}

	// return true
}

/* Function: parseCommandFile()
 * Parameters: 1. filename: The name of the file to be parsed.
 *             2. listptr:  The address of the pointer to the access control
 *                          list that will be altered as a result of
 *                          parsing the file.
 *
 * Description: This function reads an input file and alters an access
 *              control list based on the contents of the file.  The
 *              possible commands are:
 *
 *    dr: Delete Right.
 *    ar: Add Right.
 *    de: Delete Entry.
 *
 *  For example, if the file reads,
 *
 *  dr
 *  vegdahl
 *  4
 *
 *  Then the access control list should be altered so that the user
 *  'vegdahl' no longer has the right to 'write' to the file.  See
 *  aclist.h for a mapping from integers to rights.
 *
 */
func parseCommandFile(filename string) (err error) {

	// attempt to open file
	file, err := os.Open(filename)

	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Printf("%s was successfully opened.\nParsing commands from file.\n", filename)

	var d right
	var username string
	var cmd string

	// File was opened, so begin parsing
	for {
		_, err := fmt.Fscanf(file, "%s\n", &cmd)
		// If EOF, be done, otherwise there was an error
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		}

		// otherwise, this *should* be a command
		switch cmd {
		// Deleting or adding rights
		case "dr", "ar":
			// Next word is username
			if _, err := fmt.Fscanf(file, "%s\n", &username); err != nil {
				return err
			}

			// and the next is the right to add/delete
			if _, err := fmt.Fscanf(file, "%d\n", &d); err != nil {
				return err
			}

			// add or delete as needed
			switch cmd {
			case "dr":
				fmt.Printf("Delete right")
				acl.deleteRight(d, username)
			case "ar":
				fmt.Printf("Add right")
				acl.addRight(d, username)
			}

			// And add some pretty info text
			fmt.Printf(" = %d on user %s \n", d, username)

		// Deleting a whole entry
		case "de":
			// Next word is username
			if _, err := fmt.Fscanf(file, "%s\n", &username); err != nil {
				return err
			}

			fmt.Printf("Delete user %s \n", username)
			acl.deleteEntry(username)
		}
	}
}
