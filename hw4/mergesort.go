package main

import (
	"math"
)

// Function that compares two CEs
type compareFunc func(ce1, ce2 *CompanyEntry) bool

// Constants used in distance calculation
const (
	upLat    = 45.571
	upLon    = -122.726
	earthRad = 367.4447
)

// sort the CL by name
func (cl *CompanyList) sortName() {
	*cl = mergeSort(*cl, nameCompare)
}

// sort the CL by distance
func (cl *CompanyList) sortDist() {
	*cl = mergeSort(*cl, distCompare)
}

// Merge sorts a CompanyList by a specified compare function
func mergeSort(cl CompanyList, comp compareFunc) CompanyList {
	// Edge case, if the list is 1 or 0 entries
	if len(cl) < 2 {
		return cl
	}

	// Calculate the middle of the CL slice
	mid := len(cl) / 2

	// Set the left and right lists,
	// that is up to mid and mid to the end
	l := mergeSort(cl[:mid], comp)
	r := mergeSort(cl[mid:], comp)

	// And finally merge the two based on
	// the compare function provided
	return merge(l, r, comp)
}

// Merge sort merge between two lists based on a compare
// function
func merge(l, r CompanyList, comp compareFunc) CompanyList {
	// Temporary (and garbaged collected) CL slice to store
	// the newly merged entries.
	merged := make(CompanyList, 0, len(l)+len(r))

	// Loop until one of the two lists are empty
	for len(l) > 0 || len(r) > 0 {
		switch {
		// Left is empty, add right to the end and return
		case len(l) == 0:
			return append(merged, r...)

		// Right is empty, add left to the end and return
		case len(r) == 0:
			return append(merged, l...)

		// Compare the two, this is basically
		// l[0] < r[0] based on whatever comp is comparing.
		// Append the first element of l and remove it from l.
		case comp(l[0], r[0]):
			merged = append(merged, l[0])
			l = l[1:]

		// Basically l[0] >= r[0]
		// Append the first element of r and remove it from r.
		default:
			merged = append(merged, r[0])
			r = r[1:]
		}
	}

	// return the newly merged list!
	return merged
}

/* nameCompare()
 *
 * Given two pointers to CompanyEntries, return true if the first one's
 * companyName is lexicographically less than the second one.
 * Otherwise, return false.
 *
 * For example, if ce1.companyName is 'aaaa' and ce2.companyName is
 * 'aaab' then ce1.companyName is less than ce2.companyName.
 *
 */
func nameCompare(ce1, ce2 *CompanyEntry) bool {
	return ce1.companyName < ce2.companyName
}

/* distCompare()
 *
 * Given two pointers to CompanyEntry, return true if the first one's
 * latitude + longitude place it closer to the University of Portland
 * Bell Tower than the second one.  Otherwise, return false.
 *
 */
func distCompare(ce1, ce2 *CompanyEntry) bool {
	// Calculate the distances from the bell tower of both CEs
	d1 := distCalc(upLat, upLon, ce1.latitude, ce1.longitude)
	d2 := distCalc(upLat, upLon, ce2.latitude, ce2.longitude)

	// Return true if ce1 is closer
	return d1 < d2
}

/* distCalc()
 *
 * Calculates the 'great circle' distance between two points using
 * the haversine formula.
 * Details here: https://en.wikipedia.org/wiki/Haversine_formula
 *
 * Takes two sets of longitude and latitude in degrees
 * and returns distance in KM
 */
func distCalc(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert degrees to radians
	lat1 *= (math.Pi / 180)
	lon1 *= (math.Pi / 180)
	lat2 *= (math.Pi / 180)
	lon2 *= (math.Pi / 180)

	// delta lat and lon
	dlat := lat2 - lat1
	dlon := lon2 - lon1

	// First part of haversine formula
	a := (math.Sin(dlat/2) * math.Sin(dlat/2)) +
		math.Cos(lat1)*math.Cos(lat2)*(math.Sin(dlon/2)*math.Sin(dlon/2))

	// Second part of haversine formula
	c := 2 * math.Asin(math.Sqrt(a))

	// Convert to KM and return
	return earthRad * c
}
