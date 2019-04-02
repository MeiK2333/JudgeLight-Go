package JudgeLight_Go

import (
	"errors"
	"golang.org/x/sys/unix"
)

func SetProcLimit(
	cpuTimeLimit int,
	memoryLimit int,
	stackLimit int,
	outputSizeLimit int,
) error {
	var rl = unix.Rlimit{}

	// cpu time limit
	if cpuTimeLimit != 0 {
		rl.Cur = uint64(cpuTimeLimit/1000 + 1)
		rl.Max = rl.Cur
		if err := unix.Setrlimit(unix.RLIMIT_CPU, &rl); err != nil {
			return errors.New("set cpu limit failure")
		}
	}

	return nil
}
