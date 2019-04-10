package JudgeLight_Go

import (
	"fmt"
	"testing"
)

func TestRealTimeLimit(t *testing.T) {
	result, err := Run(
		"/bin/sleep",
		[]string{"3"},
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
	fmt.Println("test real time limit:", result)
	if result.realTimeUsed < 1000 || result.realTimeUsed >= 2000 {
		t.Fail()
	}
}

func TestOutputSizeLimit(t *testing.T) {
	result, err := Run(
		"/bin/bash",
		[]string{"tests/loopecho.sh"},
		[]string{},
		1000,
		1000,
		65535,
		65535,
		10,
		"",
		"tests/tmp.out",
		"",
		0,
		0,
		"",
		nil,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("test output size limit:", result)
}
