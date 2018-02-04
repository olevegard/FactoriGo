package main

import (
	"runtime/debug"
	"strconv"
	"strings"
	"testing"
)

func printStack(maxLines int) {
	for line_no, line := range strings.Split(string(debug.Stack()), "\n") {
		if line_no > maxLines {
			break
		}
		println(line)
	}
}

func assertEqual_int(expected int, actual int, t *testing.T, variable_name string) {
	assertEqual(strconv.Itoa(expected), strconv.Itoa(actual), t, variable_name)
}

func assertEqual(expected string, actual string, t *testing.T, variable_name string) {
	if expected != actual {
		printStack(10)
		t.Errorf("%v Expected to be %v but actually is %v", variable_name, expected, actual)
	}
}

func assertNotEqual(expected string, actual string, t *testing.T, variable_name string) {
	/*
		Go doesn't have assert* built in, this provides a basic assert with call_lvl
	*/
	if expected == actual {
		printStack(10)
		t.Errorf("%v Should not be the same : %v", variable_name, expected)
	}
}

func assertTrue(is_true bool, t *testing.T, variable_name string) {
	if !is_true {
		printStack(10)
		t.Errorf("%v Should be true", variable_name)
	}
}

func assertFalse(is_false bool, t *testing.T, variable_name string) {
	if is_false {
		printStack(10)
		t.Errorf("%v Should not be true", variable_name)
	}
}

func assertIntInStr(value int, str string, t *testing.T) {
	if !strings.Contains(str, strconv.Itoa(value)) {
		printStack(10)
		t.Errorf("Should have found '%d' in '%s'", value, str)
	}
}
