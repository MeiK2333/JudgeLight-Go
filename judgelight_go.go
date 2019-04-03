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
	if syscallRule == nil {
		syscallRule = make([]bool, 512)
		for i, _ := range syscallRule {
			syscallRule[i] = true
		}
	}

	result := Result{}
	var status unix.WaitStatus
	var ru unix.Rusage

	var ForkLock sync.RWMutex
	ForkLock.Lock()

	pid, _, _ := unix.Syscall(unix.SYS_FORK, 0, 0, 0)

	if pid < 0 { // Fork error
		return result, errors.New("fork failure")
	} else if pid == 0 { // Child
		// Set limit
		if err := SetProcLimit(cpuTimeLimit, memoryLimit, stackLimit, outputSizeLimit); err != nil {
			panic(err)
		}

		// Set user
		if err := SetProcUser(uid, gid); err != nil {
			panic(err)
		}

		// Chroot
		// Before SetProcStream
		if chroot != "" {
			if err := unix.Chroot(chroot); err != nil {
				panic(err)
			}
		}

		// Set stream
		if err := SetProcStream(inputFilePath, outputFilePath, errorFilePath); err != nil {
			panic(err)
		}

		// Set ptrace
		if _, _, err := unix.Syscall(unix.SYS_PTRACE, uintptr(unix.PTRACE_TRACEME), 0, 0); err != 0 {
			panic(err)
		}

		// Exec run
		if err := unix.Exec(execFilePath, append([]string{execFilePath}, execArgs...), envs); err != nil {
			panic(err)
		}

		// You will never arrive here
		unix.Exit(-1)
	}
	ForkLock.Unlock()

	// Parent

	// TODO
	// set real time limit

	exit := true
	var regs unix.PtraceRegs
	for {
		if _, err := unix.Wait4(int(pid), &status, unix.WSTOPPED, &ru); err != nil {
			return result, errors.New("wait4 failure")
		}

		// Exited
		if status.Exited() {
			goto JUDGEEND
		}

		// The program is paused but the reason is not ptrace
		// 工地英语
		if status.StopSignal() != unix.SIGTRAP {
			_, _, _ = unix.Syscall(unix.SYS_PTRACE, uintptr(unix.PTRACE_KILL), 0, 0)
			_, _ = unix.Wait4(int(pid), nil, 0, nil)
			result.reFlag = true
			result.reSignal = int(status.StopSignal())
			goto JUDGEEND
		}

		if exit {
			exit = false

			// Copy the tracee's general-purpose or floating-point registers
			if err := unix.PtraceGetRegs(int(pid), &regs); err != nil {
				return result, err
			}

			// Check Syscall
			if !AllowSyscall(syscallRule, regs.Orig_rax) {
				_, _, _ = unix.Syscall(unix.SYS_PTRACE, uintptr(unix.PTRACE_KILL), 0, 0)
				_, _ = unix.Wait4(int(pid), nil, 0, nil)
				result.reFlag = true
				result.reSyscall = int(regs.Orig_rax)
				goto JUDGEEND
			}

			// Get memory usage
			if ms, err := MemoryUsage(int(pid)); err != nil {
				return result, nil
			} else if ms.vmData > result.memoryUsed {
				result.memoryUsed = ms.vmData
			}
		} else {
			exit = true
		}

		// Restart ptrace
		if err := unix.PtraceSyscall(int(pid), 0); err != nil {
			return result, err
		}
	}

JUDGEEND:
	return result, nil
}

func AllowSyscall(syscallRule []bool, syscallId uint64) bool {
	return syscallRule[int(syscallId)]
}
