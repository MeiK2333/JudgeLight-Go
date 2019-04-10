package JudgeLight_Go

import "fmt"

type Result struct {
	cpuTimeUsed  int
	realTimeUsed int
	memoryUsed   int
	signal       int
	reFlag       bool
	reSyscall    int
}

func (r Result) String() string {
	return fmt.Sprintf("\nCPU time used:\t%d\n"+
		"real time used:\t%d\n"+
		"memory used:\t%d\n"+
		"signal:\t\t%d\n"+
		"reFlag:\t\t%t\n"+
		"reSyscall:\t%d\n", r.cpuTimeUsed, r.realTimeUsed, r.memoryUsed, r.signal, r.reFlag, r.reSyscall)
}
