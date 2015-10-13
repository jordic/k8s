package main

import "testing"

func Test_changed(t *testing.T) {

	a := []string{"a", "b"}
	b := []string{"b", "b"}

	if changed(a, b) != true {
		t.Error("Lists a, and b are different")
	}

	b = []string{"a", "b"}
	if changed(a, b) != false {
		t.Error("Lists a, and b are same")
	}

	c := []string{"b", "a"}
	if changed(a, c) != false {
		t.Error("Lists a, and b are same")
	}

}
