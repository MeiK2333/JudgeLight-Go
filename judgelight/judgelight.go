package judgelight

import (
	"errors"
	"golang.org/x/sys/unix"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func Run(
	execFilePath string,
	execArgs []string,
	envs []string,

	cpuTimeLimit int,
	realTimeLimit int,
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
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if syscallRule == nil {
		syscallRule = make([]bool, 512)
		for i, _ := range syscallRule {
			syscallRule[i] = true
		}
	}

	result := Result{}
	var status unix.WaitStatus
	var ru unix.Rusage
	// start time
	startTime := time.Now().UnixNano() / 1e6

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

	// set real time limit
	stop := make(chan bool, 1)
	// Determine if the child process exits after a period of time
	if realTimeLimit != 0 {
		ticker := time.NewTicker(time.Millisecond * time.Duration(realTimeLimit))
		go func() {
			defer ticker.Stop()

			select {
			case <-ticker.C:
				ret, _ := unix.Wait4(int(pid), &status, unix.WNOHANG, &ru)
				if ret == 0 {
					_ = unix.Kill(int(pid), unix.SIGKILL)
				}
			case <-stop:
				return
			}
		}()
	}

	exit := true
	var regs unix.PtraceRegs
	fd, err := unix.Open("/proc/"+strconv.Itoa(int(pid))+"/status", unix.O_RDONLY, 600)
	if err != nil {
		return result, err
	}
	defer unix.Close(fd)

	for {
		if _, err := unix.Wait4(int(pid), &status, unix.WSTOPPED, &ru); err != nil {
			return result, err
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
			goto JUDGEEND
		}

		if exit { // In syscall exit
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
			if ms, err := MemoryUsage(fd); err != nil {
				return result, err
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
	stop <- true
	// cpu time used
	// user time + system time
	result.cpuTimeUsed = int(ru.Utime.Sec*1000) + int(ru.Utime.Usec/1000) +
		int(ru.Stime.Sec*1000) + int(ru.Stime.Usec/1000)

	// real time used
	endTime := time.Now().UnixNano() / 1e6
	result.realTimeUsed = int(endTime - startTime)

	// signal
	result.signal = int(status.StopSignal())

	return result, nil
}

func AllowSyscall(syscallRule []bool, syscallId uint64) bool {
	return syscallRule[int(syscallId)]
}
