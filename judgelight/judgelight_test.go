package judgelight

import (
	"fmt"
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
