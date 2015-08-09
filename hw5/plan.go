package main

import (
	"errors"
	"fmt"
	"strings"
)

// type Plan struct {
// 	Courses     []CourseList
// 	courseGraph *Graph
// }

type Plan []*CourseList

func newPlan() (plPtr *Plan) {
	pl := make(Plan, 1)
	newCl := make(CourseList, 0)
	pl[0] = &newCl
	return &pl
}

// Add a course to a specified term in a Plan
// Returns an error explaining failure if needed
// nil error otherwise
//
// TERMS START AT 0!
func (p *Plan) addCourse(term int, name string, cg *Graph) (err error) {
	// Check pointer to graph
	if cg == nil {
		return errors.New("nil Graph!")
	}

	// term number sanity check
	if term < 0 {
		return errors.New("Invalid term (negative)")
	} else if term > len(*p) {
		return errors.New("Invalid term (too big)")
	}

	// trim the name of spaces
	name = strings.Trim(name, " ")

	// Look up the course in the graph
	if c, ok := cg.lookup(name); ok {
		// Go to the specified term and add the course
		for termNum, termVal := range *p {
			if termNum == term {
				*termVal = append(*termVal, c)
			}
		}
		// p[term] = append(p[term], c)
	} else if name != "" {
		return errors.New(fmt.Sprintf("Non-existent course in "+
			"course-list: %s", name))
	}
	return nil
}

// Adds an additional term to a Plan
func (p *Plan) newTerm() {
	newCl := make(CourseList, 0)
	*p = append(*p, &newCl)
}

// Nicely prints each term and classes in each
// term in a plan
func (p *Plan) printPlan() {
	for term, line := range *p {
		fmt.Printf("Term %d", term+1)
		for _, course := range *line {
			fmt.Printf(" %s", course.Name)
		}
		fmt.Println()
	}
}

// Analyze the given Plan, using a specified graph.
//
// Checks that prereqs are not violated, and courses
// are not repeated.
// There are NO co-reqs
func (p *Plan) analyze(cg *Graph) (err error) {
	// Check pointer to graph
	if cg == nil {
		return errors.New("nil Graph!")
	}

	// Go by each term
	for termNum, _ := range *p {
		// Reset course marks
		cg.resetMarks()

		// Term error tracker
		// termErr := false

		// Check for course repeats
		repeats := p.findRepeats(termNum)

		// Reset course marks
		cg.resetMarks()
		// Check for pre-req violations
		viols := p.prereqCheck(termNum)

		if !repeats && !viols {
			fmt.Printf("Term %d: no problems found.\n", termNum+1)
		}
	}

	return nil
}

// Checks for repeats in the specified term
// returning true if any are found
func (p *Plan) findRepeats(term int) (any bool) {
	// p.countCourses()

	// Count each course's appearance
	for _, term := range *p {
		// Each course
		for _, course := range *term {
			// Add one to times appeared
			course.Mark++
		}
	}

	for termNum, termList := range *p {
		// Only check specified term
		if termNum == term {
			for _, course := range *termList {
				// If the term shows up more than once,
				// print and mark any as true
				if course.Mark > 1 {
					fmt.Printf("Term %d: %s is scheduled "+
						"more than once\n", term+1, course.Name)
					any = true
				}
			}
		}
	}

	return any
}

// Check for prereq violations for the courses
// in the specified term. That is, if every course
// with prereqs has had every prereq already
// scheduled
func (p *Plan) prereqCheck(term int) (any bool) {

	// Count each course's appearance up to
	// the previous term
	for termNum, termList := range *p {
		if termNum < term {
			// Each course
			for _, course := range *termList {
				// Add one to times appeared
				course.Mark++
			}
		}
	}

	for termNum, termList := range *p {
		// Only check specified term
		if termNum == term {
			for _, course := range *termList {
				for _, pred := range course.Preds {
					// Check each course's prereqs that they
					// have been scheduled
					if pred.Mark <= 0 {
						fmt.Printf("Term %d: %s is scheduled "+
							"without prerequisite %s\n", term+1, course.Name, pred.Name)
						any = true
					}
				}
			}
		}
	}

	return any
}

// Counts via Course marks the number of times
// each course appears in a plan
// func (p *Plan) countCourses() {
// 	// Each term
// 	for _, term := range *p {
// 		// Each course
// 		for _, course := range *term {
// 			// Add one to times appeared
// 			course.Mark++
// 		}
// 	}
// }
