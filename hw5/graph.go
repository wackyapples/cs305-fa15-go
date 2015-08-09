package main

import (
	"errors"
	"fmt"
)

type Graph CourseList

// type courseMap map[string]Course

// Probably move this to graph.go eventually
const (
	WHITE int = iota
	GRAY
	BLACK
	GREY = GRAY // For the Anglophiles
)

// errors

func (g *Graph) addCourse(name string) (err error) {
	// Check if the course already exists
	if _, ok := g.lookup(name); ok {
		errStr := fmt.Sprintf("Multiple definitions of course %s", name)
		return errors.New(errStr)
	}
	// Add the course to the graph
	*g = append(*g, &Course{name, 0, make(CourseList, 0)})
	return nil
}

// Look up a course in the graph.
// If found, return the course.
// ok if false if the course is not found
func (g *Graph) lookup(name string) (c *Course, ok bool) {
	for _, val := range *g {
		if val.Name == name {
			return val, true
		}
	}
	return nil, false
}

// Adds target course as a prereq for source course.
// Returns nil error if successfully added, or returns
// and error explaining why it failed
func (g *Graph) addPrereq(source, target string) (err error) {
	srcCourse, srcOk := g.lookup(source)
	tarCourse, tarOk := g.lookup(target)
	// fmt.Printf("src(%s): %t, tarOk(%s): %t\n", source, srcOk, target, tarOk)
	switch {
	case !srcOk:
		// Source doesn't exist in the graph
		return errors.New(fmt.Sprintf("Attempt to set prerequisite "+
			"for non-existent course: %s", source))

	case !tarOk:
		// Target doesn't exist in the graph
		return errors.New(fmt.Sprintf("Attempt to set non-existent "+
			"course as prerequisite: %s", target))

	default:
		// Check if the course is already a prereq
		// fmt.Printf("srcCourse(%s) preds: %v\n", srcCourse.Name, srcCourse.Preds)
		for _, val := range srcCourse.Preds {
			if val == tarCourse {
				return errors.New(fmt.Sprintf("Course %s doubly-listed "+
					"as prerequisite for %s\n", target, source))
			}
		}
	}

	// Add the course 4reals
	srcCourse.Preds = append(srcCourse.Preds, tarCourse)
	// fmt.Printf("added %s to %s\n", target, source)

	return nil
}

// Finds any circularities among the pre-reqs in a Graph
func (g *Graph) findCircularity() (c *Course, ok bool) {
	// Reset the Graph
	g.resetMarks()

	// Iterate through the graph
	for _, ucourse := range *g {
		// Only courses with any preds/prereqs must be checked
		if len(ucourse.Preds) > 0 {
			// Reset the graph (again)
			g.resetMarks()

			// Interate through each prereq, if any are the same
			// as it's 'root' node, return that course

			for _, pcourse := range ucourse.Preds {
				if findCircHelper(ucourse, pcourse) {
					return ucourse, true
				}
			}

		}
	}

	return nil, false
}

// Recursively compares two courses for circularity,
// returning true if one exists
func findCircHelper(u, p *Course) bool {
	// Check pointers
	if u == nil || p == nil {
		return false
	}
	// Edge case
	if u == p {
		return true
	}

	// Only visit each course once
	if p.Mark == WHITE {
		p.Mark = GRAY

		// Recurse through p's prereqs (if needed)
		if len(p.Preds) > 0 {
			for _, ppred := range p.Preds {
				if findCircHelper(u, ppred) {
					return true
				}
			}
		}
	}

	// Not a circularity
	return false
}

// Actually prints off redundancies found in a Graph
// returns true if anything is found and printed
func (g *Graph) printRedundancies() (ok bool) {
	// g.resetMarks()

	// Iterate through the graph
	for _, ucourse := range *g {
		// Reset the graph for this course
		g.resetMarks()

		// Only courses with any preds/prereqs must be checked
		if len(ucourse.Preds) > 0 {

			// Check this course for a redundancy

			// Redundancy is defined as a course has another
			// as a prereq more than once, through another
			// prereq course.

			// Run the redundancy finder on the course
			ucourse.findRedundancies()

			// Check each of ucourse's pre-reqs for marked redundancies
			for _, pred := range ucourse.Preds {
				if pred.Mark == BLACK {
					ok = true
					fmt.Printf("Prerequisite %s of %s is redundant.\n",
						pred.Name, ucourse.Name)
				}
			}

		}
	}

	return ok
}

// Marks any redundant prereqs of course c as 'BLACK'
// Returns true if the course has already been marked
func (c *Course) findRedundancies() bool {

	// base case, course is marked
	if c.Mark == BLACK || c.Mark == GREY {
		return true
	}

	// Otherwise, iterate through each of c's pre-reqs (if any)
	for _, course := range c.Preds {

		// If any course is already marked...
		if course.findRedundancies() {
			// Redundancy!
			course.Mark = BLACK
		} else {
			// otherwise, mark as visited
			course.Mark = GRAY
		}
	}

	// Course is freshly marked
	return false
}

// func (c *Course) fullPrereqs() (preds CourseList, ok bool) {
// 	if len(c.Preds) == 0 {
// 		return nil, false
// 	}
// 	for _, pred := range c.Preds {
// 		var nlPreds CourseList
// 		if len(pred.Preds) != 0 {
// 			nlPreds, _ = pred.fullPrereqs()
// 		}
// 		preds = append(preds, nlPreds...)
// 	}

// 	return preds, true
// }

// Reset all marks in the graph to 'WHITE'
func (g *Graph) resetMarks() {
	for _, vertex := range *g {
		vertex.Mark = WHITE
	}
}
