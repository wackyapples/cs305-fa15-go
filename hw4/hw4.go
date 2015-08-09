package main

import (
	"flag"
	"fmt"
	"io"
	// "log"
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"
)

// The glorious CompanyEntry type.
type CompanyEntry struct {
	companyName        string
	companyDescription string
	website            string
	streetAddr         string
	suiteNumber        string
	city               string
	state              string
	zip                int64
	latitude           float64
	longitude          float64
}

// Package variables
var (
	compList CompanyList
	compTree *CompanyTree
)

// flags
var (
	verboseFlag  = flag.Bool("v", false, "Verbose mode")
	treeFlag     = flag.Bool("t", false, "Tree mode")
	mergeFlag    = flag.Bool("m", false, "Mergesort mode")
	distanceFlag = flag.Bool("d", false, "Distance mode")
	gosortFlag   = flag.Bool("g", false, "GoSort Mode")
)

// Initialize stuff
func init() {
	// init package variables
	// Empty, yet initialized slice of CompanyEntries
	compList = make(CompanyList, 0)
	// This makes a new tree defaulting values to their
	// 'zeros', that is, strings are empty, numbers are 0
	// and pointers are nil. Isn't Go cool?
	compTree = new(CompanyTree)

	// flag errors print custom usage info
	flag.Usage = func() {
		printUsage()
	}

	// Parse the flags
	flag.Parse()
}

// Prints usage info and exits with value of 1
func printUsage() {
	fmt.Printf("usage: %s -[v|t|m|d|g] <input file> . \n", os.Args[0])
	os.Exit(1)
}

func main() {

	// Check for legal parameter usage, 1 or 0 flags and 1 non-flag
	if flag.NFlag() > 1 || flag.NArg() != 1 {
		printUsage()
	}

	// Set the filename as the first (and only) non-flag variable
	filename := flag.Args()[0]

	// Parse the file into a slice of CompanyEntries
	// Need to decide if that slice is a package variable
	// or a pointer to be passed
	// I hath decreed there to be package variables
	if err := parseFile(filename); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Successfully parsed file,", filename)

	// The list is printed 'backwards' first, as this
	// is how the original assignment printed its list.
	compList.printListInverse()
	fmt.Printf("There are %d nodes in the list.\n\n", len(compList))

	// Do what the set flag says, if any flag other than -v is
	// set at all, as -v is already done!
	switch {
	case *treeFlag:
		fmt.Printf("\n****Sorted alphabetically****\n")
		compTree.insertList(compList)
		compTree.printTree()

	case *mergeFlag:
		fmt.Printf("\n****Sorted alphabetically****\n")
		compList.sortName()
		compList.printList()

	case *distanceFlag:
		fmt.Printf("\n****Sorted by distance****\n")
		compList.sortDist()
		compList.printList()

	case *gosortFlag:
		fmt.Printf("\n****Sorted alphabetically****\n")
		sort.Sort(byName(compList))
		compList.printList()
	}
}

func parseFile(filename string) error {
	// var file *os.File

	// Try to open the file, return any errs
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	// Close file upon return
	defer file.Close()

	// new Reader object for parsing the file
	reader := bufio.NewReader(file)

	// FROM ORIGINAL SOURCE
	// If the input file is well-formed, the first character
	// of each company entry is preceded by a '*'.  At the
	// top of each loop iteration, affirm that the * is there.
	// if not, the file is not well-formed, it has not followed
	// the expected format, and we should return with an error code.

	for {
		// First character, first line
		a, _, err := reader.ReadRune()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		}

		// Temp entry to build from file
		var newEntry *CompanyEntry

		// If the first character was a '*', it is the start of
		// an entry
		if a == '*' {
			// make a new empty entry
			newEntry = new(CompanyEntry)
			// Parse the entry!!!
			chomp(reader, &newEntry.companyName)
			chomp(reader, &newEntry.companyDescription)
			chomp(reader, &newEntry.website)
			chomp(reader, &newEntry.streetAddr)
			chomp(reader, &newEntry.suiteNumber)
			chomp(reader, &newEntry.city)
			chomp(reader, &newEntry.state)
			chomp(reader, &newEntry.zip)
			chomp(reader, &newEntry.latitude)
			chomp(reader, &newEntry.longitude)

			// Add the new entry
			// NOTE: In the original program, items added
			// to the front of the list, here they are
			// added to the end.
			compList = append(compList, newEntry)

			// Adds new entries to the front, VERY SLOW
			// compList = append(CompanyList{newEntry}, compList...)
		}
	}
}

// Yeah, remember all those giant C functions that were reduced
// to like, 3 lines? Yeah, not this one. On the bright side,
// great chance to play around with a type switch!
//
// Anyway, this function takes a *bufio.Reader, and an
// input pointer to a string, int64, or float64.
// The next line from the reader is read and the newline
// character at the end is removed and the line is parsed
// as appropriate (determined by passed pointer).
//
// **The 'input' itself is modified!**
func chomp(reader *bufio.Reader, input interface{}) {

	// Used during number parsing
	var readErr error

	// Attempt to read the next line
	if readStr, _, err := reader.ReadLine(); err != nil {
		fmt.Println(err)
		return
	} else {
		// Trim the newline from the line and surrounding
		// spaces.
		// readStr = strings.TrimSuffix(readStr, "\n")
		// readStr = strings.TrimPrefix(readStr, " ")
		// readStr = strings.TrimSuffix(readStr, " ")
		readStr := string(readStr)
		readStr = strings.Trim(readStr, " ")

		// Woo! A type switch. Fancy.
		switch t := input.(type) {
		// If no type matches, print an error and return.
		default:
			fmt.Printf("Bad type: %T\n", t)
			return

		// If the input is a *string, set it to readStr (and be done)
		case *string:
			*t = readStr

		// If the input is a *int64 (zip code), parse it
		case *int64:
			if readInt, err := strconv.ParseInt(readStr, 10, 64); err == nil {
				*t = readInt
			} else {
				readErr = err
			}

		// If the input is a *float64 (lon/lat), parse it
		case *float64:
			if readFloat, err := strconv.ParseFloat(readStr, 64); err == nil {
				*t = readFloat
			} else {
				readErr = err
			}
		}
	}

	// If any parse errors occurred, print them.
	if readErr != nil {
		fmt.Println(readErr)
	}
}
