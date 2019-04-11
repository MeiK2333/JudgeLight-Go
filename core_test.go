package JudgeLight_Go

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestRun(t *testing.T) {
	result, err := Run(
		"/bin/echo",
		[]string{"run"},
		[]string{},
		1000,
		1000,
		65535,
		65535,
		655350,
		"",
		"",
		"",
		0,
		0,
		"",
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("test run:", result)
}

func TestHello(t *testing.T) {
	result, err := Run(
		"/bin/echo",
		[]string{"Hello", "World"},
		[]string{},
		1000,
		1000,
		65535,
		65535,
		655350,
		"",
		"",
		"",
		0,
		0,
		"",
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("test hello:", result)
}

func TestOutput(t *testing.T) {
	result, err := Run(
		"/bin/echo",
		[]string{"Hello", "World"},
		[]string{},
		1000,
		1000,
		65535,
		65535,
		655350,
		"",
		"tests/test_output.txt",
		"",
		0,
		0,
		"",
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("test output:", result)

	if dat, err := ioutil.ReadFile("tests/test_output.txt"); err != nil || string(dat) != "Hello World\n" {
		t.Fail()
	}
}

func TestErrorOutput(t *testing.T) {
	result, err := Run(
		"/usr/bin/python3",
		[]string{"tests/erroroutput.py"},
		[]string{},
		1000,
		1000,
		65535,
		65535,
		655350,
		"",
		"",
		"tests/test_error.txt",
		0,
		0,
		"",
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("test error output:", result)

	if dat, err := ioutil.ReadFile("tests/test_error.txt"); err != nil || string(dat) != "Hello World" {
		t.Fail()
	}
}

func TestInput(t *testing.T) {
	result, err := Run(
		"/usr/bin/python3",
		[]string{"tests/testinput.py"},
		[]string{},
		1000,
		1000,
		65535,
		65535,
		655350,
		"tests/input.txt",
		"tests/test_input.txt",
		"",
		0,
		0,
		"",
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("test error output:", result)

	if dat, err := ioutil.ReadFile("tests/test_input.txt"); err != nil || string(dat) != "3\n" {
		t.Fail()
	}
}
