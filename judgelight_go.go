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
		// set limit
		if err := SetProcLimit(cpuTimeLimit, memoryLimit, stackLimit, outputSizeLimit); err != nil {
			panic(err)
		}

		// set user
		if err := SetProcUser(uid, gid); err != nil {
			panic(err)
		}

		// chroot
		// before SetProcStream
		if chroot != "" {
			if err := unix.Chroot(chroot); err != nil {
				panic(err)
			}
		}

		// set stream
		if err := SetProcStream(inputFilePath, outputFilePath, errorFilePath); err != nil {
			panic(err)
		}

		// set ptrace
		if _, _, err := unix.Syscall(unix.SYS_PTRACE, uintptr(unix.PTRACE_TRACEME), 0, 0); err != 0 {
			panic(err)
		}

		// exec run
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
