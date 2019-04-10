package JudgeLight_Go

import (
	"fmt"
	"testing"
)

func TestCpuTimeLimit(t *testing.T) {
	result, err := Run(
		"/usr/bin/python3",
		[]string{"tests/loopcpu.py"},
		[]string{},
		1000,
		5000,
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
	fmt.Println("test cpu time limit:", result)
	if result.cpuTimeUsed < 1000 || result.cpuTimeUsed >= 3000 {
		t.Fail()
	}
}

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
	if !result.reFlag {
		t.Fail()
	}
}
