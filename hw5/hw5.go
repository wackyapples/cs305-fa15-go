package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// types
type CourseList []*Course

type Course struct {
	Name  string
	Mark  int
	Preds CourseList
}

// Package vars. Maybe keep.
var (
	myGraph  *Graph
	myPlan   *Plan
	maxPreds int
)

func main() {

	// Parameter check
	if len(os.Args) != 3 {
		fmt.Println("Need exactly two arguments")
		os.Exit(1)
	}

	// var myGraph Graph

	// Parse the pre-reqs file
	if cl, err := parsePrereqs(os.Args[1]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		myGraph = cl
	}

	// Print sizes of things
	fmt.Printf("Total number of courses: %d; "+
		"largest number of prerequisites for a course: %d.\n", len(*myGraph), maxPreds)

	// Check for circularities, bail if one is found
	if circ, ok := myGraph.findCircularity(); ok {
		fmt.Printf("Prerequisite circularity detected for %s\n", circ.Name)
		os.Exit(1)
	} else {
		fmt.Println("No prerequisite circularities found.")
	}

	// Check for redundancies, not a fatal error
	if !myGraph.printRedundancies() {
		fmt.Println("No prerequisite redundancies found.")
	}

	// Time to tackle schedule, which means another parser!
	// The graph is included for course lookup
	if pl, err := parsePlan(os.Args[2], myGraph); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		myPlan = pl
	}

	// Print proposed schedule
	fmt.Println("================proposed schedule================")
	myPlan.printPlan()

	// analyze the student's schedule
	fmt.Println("================schedule analysis================")
	myPlan.analyze(myGraph)

	// And we're done :)

}

// Parse the prereqs file and return a Graph of the courses.
// If an error occurs, it is returned as well (with a nil Graph)
func parsePrereqs(filename string) (cgPtr *Graph, err error) {

	// Attempt to open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	cg := make(Graph, 0)

	prereqTemp := map[string][]string{}

Loop:
	for {
		str, err := reader.ReadString('\n')
		switch {
		case err == io.EOF:
			// Go breaks out of switches, so a label is needed
			break Loop
		case err != nil:
			return nil, err
		}

		// Remove newline
		line := strings.Trim(str, "\n")

		if strings.HasSuffix(line, ":") {
			// The line has only one entry, and it's a
			// course without pre-reqs
			line = strings.Trim(line, ":")
			cg.addCourse(line)

		} else if line != "" {
			// There's more than one course, and it's not an empty line

			// Split the first entry from its prereqs
			lineList := strings.Split(line, ":")

			// lineList[0] is the first entry that ends with ":"
			// and is to be added to the graph.
			//
			// lineList[1] is the remaining courses separated by
			// spaces to be added as prereqs.

			// set up the new course
			cg.addCourse(lineList[0])

			// Get a list of the courses to add as prereqs,
			// discarding the first empty entry (caused by
			// preceding space)
			prereqs := strings.Split(lineList[1], " ")

			// Need to store pre-reqs for when done parsing
			// top-level classes (GC'd)
			prereqTemp[lineList[0]] = prereqs
		}
	}

	// Attempt to add the prereqs
	for key, pres := range prereqTemp {
		numPreds := 0

		for _, pre := range pres {
			// ignore the occasional blank string from extra spaces
			if pre != "" {
				if preErr := cg.addPrereq(key, pre); preErr != nil {
					fmt.Println(preErr)
				} else {
					numPreds++
				}
			}
		}

		// Count maximum number of prereqs in whole graph
		if numPreds > maxPreds {
			maxPreds = numPreds
		}
	}

	return &cg, nil
}

// Reads and parses a plan file
// In addtion to the filename, the Graph of
// all the courses is needed.
func parsePlan(filename string, cg *Graph) (pl *Plan, err error) {
	// Lots of boilerplate stuff,
	// maybe merge parsing functions someday

	// Attempt to open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	// pl = make(Plan, 0)
	// Build new plan
	pl = newPlan()

Loop:
	for term := 0; ; term++ {
		str, err := reader.ReadString('\n')
		switch {
		case err == io.EOF:
			// Go breaks out of switches, so a label is needed
			break Loop
		case err != nil:
			return nil, err
		}

		// Add a new term if needed
		if term >= len(*pl) {
			pl.newTerm()
		}

		// Remove newline
		line := strings.Trim(str, "\n")

		// Don't mess with empty lines
		if line != "" {
			// Split the line up, with a slice containing
			// the name of each course
			lineList := strings.Split(line, " ")

			// Attempt to add each course to the current term
			for _, course := range lineList {
				if err = pl.addCourse(term, course, cg); err != nil {
					fmt.Println(err)
				}
			}
		}
	}

	return pl, nil
}
