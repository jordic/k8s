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

func Test_toIpString_when_slice_has_multiple_values(t *testing.T) {

	slice := []string{"a  ", "b"}

	ipString := toIpString(slice)

	expected := "a/32,b/32,"
	if ipString != expected {
		t.Error("toIpString doesn't concatenate IP's correctly. Expected: " + expected + ", Actual: " + ipString)
	}
}

func Test_toIpString_when_slice_has_no_values(t *testing.T) {

	slice := []string{}

	ipString := toIpString(slice)

	expected := ""
	if ipString != expected {
		t.Error("toIpString doesn't concatenate IP's correctly. Expected: " + expected + ", Actual: " + ipString)
	}
}

func Test_toIpString_when_slice_has_empty_value(t *testing.T) {

	slice := []string{""}

	ipString := toIpString(slice)

	expected := ""
	if ipString != expected {
		t.Error("toIpString doesn't concatenate IP's correctly. Expected: " + expected + ", Actual: " + ipString)
	}
}

func Test_toIpString_when_slice_has_single_value(t *testing.T) {

	slice := []string{"10.10.10.10"}

	ipString := toIpString(slice)

	expected := "10.10.10.10/32,"
	if ipString != expected {
		t.Error("toIpString doesn't concatenate IP's correctly. Expected: " + expected + ", Actual: " + ipString)
	}
}
