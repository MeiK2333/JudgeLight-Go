package JudgeLight_Go

import (
	"fmt"
	"testing"
)

func TestRun(t *testing.T) {
	result, err := Run(
		"/bin/echo",
		[]string{"Hello", "World"},
		[]string{},
		1000,
		1000,
		65535,
		655350,
		65535,
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
	fmt.Println(result)
}
