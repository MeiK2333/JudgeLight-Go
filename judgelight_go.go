package JudgeLight_Go

import (
	"errors"
	"golang.org/x/sys/unix"
	"sync"
)

func Run(
	execFilePath string,
	execArgs []string,
	envs []string,

	cpuTimeLimit int,
	readTimeLimit int,
	memoryLimit int,
	stackLimit int,
	outputSizeLimit int,

	inputFilePath string,
	outputFilePath string,
	errorFilePath string,

	uid int,
	gid int,
	chroot string,
	syscallRule []bool,
) (Result, error) {
	result := Result{}

	var ForkLock sync.RWMutex
	ForkLock.Lock()
	pid, _, _ := unix.Syscall(unix.SYS_FORK, 0, 0, 0)

	if pid < 0 { // fork error
		return result, errors.New("fork failure")
	} else if pid == 0 { // child

		// run
		if err := unix.Exec(execFilePath, append([]string{execFilePath}, execArgs...), envs); err != nil {
			panic(err)
		}

		// You will never arrive here
		unix.Exit(-1)
	}
	ForkLock.Unlock()

	// parent
	return result, nil
}
