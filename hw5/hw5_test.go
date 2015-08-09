package main

import (
	"testing"
)

func TestCirc(t *testing.T) {
	// testGraph := make(Graph, 0)

	testGraph, err := parsePrereqs("prereqs6.txt")

	if err != nil {
		t.Error(err)
	}

	if circ, ok := testGraph.findCircularity(); !ok || circ.Name != "THE105" {
		t.Errorf("Fail: %s\n", circ.Name)
	}
}
