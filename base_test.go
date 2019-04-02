package JudgeLight_Go

import (
	"fmt"
	"golang.org/x/sys/unix"
	"sync"
	"testing"
	"time"
)

func TestFork(t *testing.T) {
	var ForkLock sync.RWMutex
	ForkLock.Lock()

	pid, _, _ := unix.Syscall(unix.SYS_FORK, 0, 0, 0)
	if pid == 0 {
		unix.Exit(0)
	} else {
		time.Sleep(time.Second)
	}

	ForkLock.Unlock()
}

func TestLimit(t *testing.T) {
	var ForkLock sync.RWMutex
	ForkLock.Lock()

	pid, _, _ := unix.Syscall(unix.SYS_FORK, 0, 0, 0)
	if pid == 0 {
		var rl = unix.Rlimit{}
		rl.Cur = 1
		rl.Max = 2
		_ = unix.Setrlimit(unix.RLIMIT_CPU, &rl)
		for {

		}
		unix.Exit(0)
	} else {
		var status unix.WaitStatus
		var ru unix.Rusage
		unix.Wait4(int(pid), &status, 0, &ru)
		if status.Signal().String() != "killed" {
			t.Fail()
		}
	}

	ForkLock.Unlock()
}

func TestRealTimeLimit(t *testing.T) {
	var ForkLock sync.RWMutex
	ForkLock.Lock()
	pid, _, _ := unix.Syscall(unix.SYS_FORK, 0, 0, 0)
	if pid == 0 {
		if err := unix.Exec("/bin/sleep", []string{"sleep", "2"}, []string{}); err != nil {
			t.Fail()
		}
		unix.Exit(0)
	}
	ForkLock.Unlock()

	var status unix.WaitStatus
	var ru unix.Rusage
	ticker := time.NewTicker(time.Second)

	go func() {
		defer ticker.Stop()

		select {
		case <-ticker.C:
			ret, err := unix.Wait4(int(pid), &status, unix.WNOHANG, &ru)
			if err != nil {
				panic(err)
			}
			if ret == 0 {
				fmt.Println("out real time limit")
				unix.Kill(int(pid), unix.SIGKILL)
			}
		}
	}()

	time.Sleep(3 * time.Second)
	ticker.Stop()
}

func TestPtrace(t *testing.T) {
	var regs unix.PtraceRegs
	var ForkLock sync.RWMutex
	ForkLock.Lock()

	pid, _, _ := unix.Syscall(unix.SYS_FORK, 0, 0, 0)
	if pid == 0 {
		if _, _, err := unix.Syscall(unix.SYS_PTRACE, uintptr(unix.PTRACE_TRACEME), 0, 0); err != 0 {
			t.Fail()
		}
		if err := unix.Exec("/bin/date", []string{"date"}, []string{}); err != nil {
			t.Fail()
		}
		unix.Exit(-1)
	} else {
		if _, err := unix.Wait4(int(pid), nil, 0, nil); err != nil {
			t.Fail()
		}
		exit := true
		for {
			if exit {
				if err := unix.PtraceGetRegs(int(pid), &regs); err != nil {
					break
				}
				//fmt.Println(regs.Orig_rax)
			}

			if err := unix.PtraceSyscall(int(pid), 0); err != nil {
				t.Fail()
			}

			if _, err := unix.Wait4(int(pid), nil, 0, nil); err != nil {
				t.Fail()
			}

			exit = !exit
		}
	}

	ForkLock.Unlock()
}
